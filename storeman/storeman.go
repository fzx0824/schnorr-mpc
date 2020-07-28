package storeman

import (
	"context"
	"github.com/wanchain/schnorr-mpc/storeman/osmconf"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"os"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/wanchain/schnorr-mpc/accounts"
	"github.com/wanchain/schnorr-mpc/log"
	"github.com/wanchain/schnorr-mpc/p2p"
	"github.com/wanchain/schnorr-mpc/p2p/discover"
	"github.com/wanchain/schnorr-mpc/rpc"
	"github.com/wanchain/schnorr-mpc/storeman/storemanmpc"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"github.com/wanchain/schnorr-mpc/storeman/validator"
)

type Config struct {
	StoremanNodes     []*discover.Node
	Password          string
	WorkingPwd        string
	DataPath          string
	SchnorrThreshold  int
	SchnorrTotalNodes int
}

var DefaultConfig = Config{
	StoremanNodes:     make([]*discover.Node, 0),
	SchnorrThreshold:  26,
	SchnorrTotalNodes: 50,
}

type StrmanKeepAlive struct {
	version   int
	magic     int
	recipient discover.NodeID
}

type StrmanKeepAliveOk struct {
	version int
	magic   int
	status  int
}

type StrmanAllPeers struct {
	Ip     []string
	Port   []string
	Nodeid []string
}

type StrmanGetPeers struct {
	LocalPort string
}

const keepaliveMagic = 0x33

// New creates a Whisper client ready to communicate through the Ethereum P2P network.
func New(cfg *Config, accountManager *accounts.Manager, aKID, secretKey, region string) *Storeman {
	storeman := &Storeman{
		peers:      make(map[discover.NodeID]*Peer),
		quit:       make(chan struct{}),
		cfg:        cfg,
		isSentPeer: false,
		peersPort:  make(map[discover.NodeID]string),
	}

	storeman.mpcDistributor = storemanmpc.CreateMpcDistributor(accountManager,
		storeman,
		aKID,
		secretKey,
		region,
		cfg.Password)

	dataPath := filepath.Join(cfg.DataPath, "storeman", "data")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dataPath, 0700); err != nil {
			log.SyslogErr("make Storeman path fail", "err", err.Error())
		}
	}
	log.Info("==================================")
	log.Info("=========New storeman", "DB file path", dataPath)
	log.Info("==================================")
	validator.NewDatabase(dataPath)
	// p2p storeman sub protocol handler
	storeman.protocol = p2p.Protocol{
		Name:    mpcprotocol.PName,
		Version: uint(mpcprotocol.PVer),
		Length:  mpcprotocol.NumberOfMessageCodes,
		Run:     storeman.HandlePeer,
		NodeInfo: func() interface{} {
			return map[string]interface{}{
				"version": mpcprotocol.PVerStr,
			}
		},
	}

	return storeman
}

////////////////////////////////////
// Storeman
////////////////////////////////////
type Storeman struct {
	protocol       p2p.Protocol
	peers          map[discover.NodeID]*Peer
	storemanPeers  map[discover.NodeID]bool
	peerMu         sync.RWMutex  // Mutex to sync the active peer set
	quit           chan struct{} // Channel used for graceful exit
	mpcDistributor *storemanmpc.MpcDistributor
	cfg            *Config
	server         *p2p.Server
	isSentPeer     bool
	peersPort      map[discover.NodeID]string

	//allPeersConnected chan bool
}

// MaxMessageSize returns the maximum accepted message size.
func (sm *Storeman) MaxMessageSize() uint32 {
	// TODO what is the max size of storeman???
	return uint32(1024 * 1024)
}

// runMessageLoop reads and processes inbound messages directly to merge into client-global state.
func (sm *Storeman) runMessageLoop(p *Peer, rw p2p.MsgReadWriter) error {

	log.SyslogInfo("runMessageLoop begin")
	defer log.SyslogInfo("runMessageLoop exit")

	for {
		// fetch the next packet
		packet, err := rw.ReadMsg()
		if err != nil {
			log.SyslogErr("runMessageLoop", "peer", p.Peer.ID().String(), "err", err.Error())
			return err
		}

		switch packet.Code {

		default:

			log.SyslogInfo("runMessageLoop, received a msg", "peer", p.Peer.ID().String(), "packet size", packet.Size)
			if packet.Size > sm.MaxMessageSize() {
				log.SyslogWarning("runMessageLoop, oversized message received", "peer", p.Peer.ID().String(), "packet size", packet.Size)
			} else {
				err = sm.mpcDistributor.GetMessage(p.Peer.ID(), rw, &packet)
				if err != nil {
					log.SyslogErr("runMessageLoop, distributor handle msg fail", "err", err.Error())
				}
			}
		}

		packet.Discard()
	}
}

// APIs returns the RPC descriptors the Whisper implementation offers
func (sm *Storeman) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: mpcprotocol.PName,
			Version:   mpcprotocol.PVerStr,
			Service:   &StoremanAPI{sm: sm},
			Public:    true,
		},
	}
}

// Protocols returns the whisper sub-protocols ran by this particular client.
func (sm *Storeman) Protocols() []p2p.Protocol {
	return []p2p.Protocol{sm.protocol}
}

// Start implements node.Service, starting the background data propagation thread
// of the Whisper protocol.
func (sm *Storeman) Start(server *p2p.Server) error {

	sm.mpcDistributor.Self = server.Self()

	// set self node id into the osm config
	osmconf.GetOsmConf().SetSelfNodeId(&sm.mpcDistributor.Self.ID)
	storemanNodes, err := osmconf.GetOsmConf().GetAllPeersNodeIds()
	if err != nil {
		log.SyslogErr("Storeman", "start err ", err.Error())
		return err
	}

	sm.mpcDistributor.StoreManGroup = make([]discover.NodeID, len(storemanNodes))
	sm.storemanPeers = make(map[discover.NodeID]bool)
	sm.server = server

	for i, item := range storemanNodes {
		sm.mpcDistributor.StoreManGroup[i] = item
		sm.storemanPeers[item] = true
	}

	return nil

}

func (sm *Storeman) checkPeerInfo() {

	// Start the tickers for the updates
	keepQuest := time.NewTicker(mpcprotocol.KeepaliveCycle * time.Second)

	leaderid, err := discover.BytesID(sm.cfg.StoremanNodes[0].ID.Bytes())
	if err != nil {
		log.Info("err decode leader node id from config")
	}

	if sm.cfg.StoremanNodes[0].ID.String() == sm.server.Self().ID.String() {
		return
	}

	//leader will not checkPeerInfo
	log.Info("Entering checkPeerInfo")
	// Loop and transmit until termination is requested
	for {

		select {
		case <-keepQuest.C:
			//log.Info("Entering checkPeerInfo for loop")
			if sm.IsActivePeer(&leaderid) {
				splits := strings.Split(sm.server.ListenAddr, ":")
				sm.SendToPeer(&leaderid, mpcprotocol.GetPeersInfo, StrmanGetPeers{splits[len(splits)-1]})
			} else {
				log.Info("leader is connecting...")
				sm.server.AddPeer(sm.server.StoremanNodes[0])
			}

		}
	}
}

// Stop implements node.Service, stopping the background data propagation thread
// of the Whisper protocol.
func (sm *Storeman) Stop() error {
	return nil
}

func (sm *Storeman) SendToPeer(peerID *discover.NodeID, msgcode uint64, data interface{}) error {
	sm.peerMu.RLock()
	defer sm.peerMu.RUnlock()
	peer, exist := sm.peers[*peerID]
	if exist {
		return p2p.Send(peer.ws, msgcode, data)
	} else {
		log.SyslogWarning("peer not find", "peer", peerID.String())
	}
	return nil
}

func (sm *Storeman) IsActivePeer(peerID *discover.NodeID) bool {
	sm.peerMu.RLock()
	defer sm.peerMu.RUnlock()
	_, exist := sm.peers[*peerID]
	return exist
}

// HandlePeer is called by the underlying P2P layer when the whisper sub-protocol
// connection is negotiated.
func (sm *Storeman) HandlePeer(peer *p2p.Peer, rw p2p.MsgReadWriter) error {

	if _, exist := sm.storemanPeers[peer.ID()]; !exist {
		return errors.New("Peer is not in storemangroup")
	}

	log.Info("handle new peer", "remoteAddr", peer.RemoteAddr().String(), "peerID", peer.ID().String())

	// Create the new peer and start tracking it
	storemanPeer := newPeer(sm, peer, rw)

	sm.peerMu.Lock()
	sm.peers[storemanPeer.ID()] = storemanPeer
	sm.peerMu.Unlock()

	// Run the peer handshake and state updates
	if err := storemanPeer.handshake(); err != nil {
		log.SyslogErr("storemanPeer.handshake failed", "peerID", peer.ID().String(), "err", err.Error())
		return err
	}

	defer func() {
		sm.peerMu.Lock()

		delete(sm.peers, storemanPeer.ID())

		for _, smnode := range sm.server.StoremanNodes {
			if smnode.ID == storemanPeer.ID() {
				log.Info("remove peer", "pid", smnode.ID)
				sm.server.RemovePeer(smnode)
				break
			}
		}

		sm.peerMu.Unlock()
	}()

	storemanPeer.start()
	defer storemanPeer.stop()

	return sm.runMessageLoop(storemanPeer, rw)
}

////////////////////////////////////
// StoremanAPI
////////////////////////////////////
type StoremanAPI struct {
	sm *Storeman
}

func (sa *StoremanAPI) Version(ctx context.Context) (v string) {
	return mpcprotocol.PVerStr
}

func (sa *StoremanAPI) Peers(ctx context.Context) []*p2p.PeerInfo {
	var ps []*p2p.PeerInfo
	for _, p := range sa.sm.peers {
		ps = append(ps, p.Peer.Info())
	}

	return ps
}

func (sa *StoremanAPI) FreshGrpInfo(ctx context.Context) error {
	log.SyslogInfo("FreshGrpInfo begin")
	err := osmconf.GetOsmConf().FreshCnf(osmconf.GetOsmConf().GrpInfoPath())
	if err != nil {
		log.SyslogErr("FreshGrpInfo error", "error", err.Error())
	}
	log.SyslogInfo("FreshGrpInfo end")
	return err
}

func (sa *StoremanAPI) SignDataByApprove(ctx context.Context, data mpcprotocol.SendData) (result interface{}, err error) {

	sa.sm.mpcDistributor.SetCurPeerCount(uint16(len(sa.sm.peers)))
	PKBytes := data.PKBytes
	CurveBytes := data.Curve
	//signed, err := sa.sm.mpcDistributor.CreateReqMpcSign([]byte(data.Data), PKBytes)
	signed, err := sa.sm.mpcDistributor.CreateReqMpcSign([]byte(data.Data), []byte(data.Extern), PKBytes, 1, CurveBytes)

	// signed   R // s
	if err == nil {
		log.SyslogInfo("SignMpcTransaction end")
	} else {
		log.SyslogErr("SignMpcTransaction end", "err", err.Error())
		return mpcprotocol.SignedResult{R: []byte{}, S: []byte{}}, err
	}

	return signed, nil
}

func (sa *StoremanAPI) SignData(ctx context.Context, data mpcprotocol.SendData) (result interface{}, err error) {

	PKBytes := data.PKBytes
	CurveBytes := data.Curve
	sa.sm.mpcDistributor.SetCurPeerCount(uint16(len(sa.sm.peers)))
	//signed, err := sa.sm.mpcDistributor.CreateReqMpcSign([]byte(data.Data), PKBytes)
	signed, err := sa.sm.mpcDistributor.CreateReqMpcSign([]byte(data.Data), []byte(data.Extern), PKBytes, 0, CurveBytes)

	// signed   R // s
	if err == nil {
		log.SyslogInfo("SignMpcTransaction end")
	} else {
		log.SyslogErr("SignMpcTransaction end", "err", err.Error())
		return mpcprotocol.SignedResult{R: []byte{}, S: []byte{}}, err
	}

	return signed, nil
}

func (sa *StoremanAPI) AddValidData(ctx context.Context, data mpcprotocol.SendData) error {
	return validator.AddValidData(&data)
}

// non leader node polling the data received from leader node
func (sa *StoremanAPI) GetDataForApprove(ctx context.Context) ([]mpcprotocol.SendData, error) {
	return validator.GetDataForApprove()
}

//// non leader node ApproveData, and make sure that the data is really required to be signed by them.
func (sa *StoremanAPI) ApproveData(ctx context.Context, data []mpcprotocol.SendData) []error {
	return validator.ApproveData(data)
}
