package step

import (
	"bytes"
	"github.com/wanchain/schnorr-mpc/common"
	"github.com/wanchain/schnorr-mpc/log"
	"github.com/wanchain/schnorr-mpc/p2p/discover"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
)

type AckMpcGPKStep struct {
	BaseStep
	message       map[discover.NodeID]bool
	mpcGPK        []byte
	remoteMpcGPKs map[discover.NodeID][]byte
}

func CreateAckMpcGPKStep(peers *[]mpcprotocol.PeerInfo) *AckMpcGPKStep {
	return &AckMpcGPKStep{
		*CreateBaseStep(peers, -1),
		make(map[discover.NodeID]bool),
		nil,
		make(map[discover.NodeID][]byte)}
}

func (ack *AckMpcGPKStep) InitStep(result mpcprotocol.MpcResultInterface) error {
	log.SyslogInfo("AckMpcAccountStep.InitStep begin")
	mpcGpk, err := result.GetByteValue(mpcprotocol.MpcContextResult)
	if err != nil {
		log.SyslogErr("AckMpcGPKStep::InitStep","ack mpc account step, init fail. err", err.Error())
		return err
	}

	// Check valid of PK ?

	ack.mpcGPK = mpcGpk
	return nil
}

func (ack *AckMpcGPKStep) CreateMessage() []mpcprotocol.StepMessage {
	return []mpcprotocol.StepMessage{mpcprotocol.StepMessage{
		MsgCode:   mpcprotocol.MPCMessage,
		PeerID:    nil,
		Peers:     nil,
		Data:      nil,
		BytesData: [][]byte{ack.mpcGPK}}}
}

func (ack *AckMpcGPKStep) FinishStep(result mpcprotocol.MpcResultInterface, mpc mpcprotocol.StoremanManager) error {
	log.SyslogInfo("AckMpcAccountStep.FinishStep begin")
	err := ack.BaseStep.FinishStep()
	if err != nil {
		return err
	}

	if len(ack.remoteMpcGPKs) != len(*ack.BaseStep.peers) {
		log.SyslogErr("AckMpcGPKStep:FinishStep","ack mpc account step, finish, invalid remote mpc address. peer num, mpc addr num",
			len(*ack.BaseStep.peers), len(ack.remoteMpcGPKs))
		return mpcprotocol.ErrInvalidMPCAddr
	}

	for peerID, mpcGpk := range ack.remoteMpcGPKs {
		if mpcGpk == nil {
			log.SyslogErr("AckMpcGPKStep:FinishStep","ack mpc account step, finish, invalid remote mpc address: nil. peerID",
				peerID.String())
			return mpcprotocol.ErrInvalidMPCAddr
		}

		if !bytes.Equal(ack.mpcGPK, mpcGpk) {
			log.SyslogErr("AckMpcGPKStep:FinishStep","ack mpc account step, finish, invalid remote mpc address. local, received, peerID",
				common.ToHex(ack.mpcGPK),
				common.ToHex(mpcGpk),
				peerID.String())
			return mpcprotocol.ErrInvalidMPCAddr
		}
	}

	return nil
}

func (ack *AckMpcGPKStep) HandleMessage(msg *mpcprotocol.StepMessage) bool {
	log.SyslogInfo("AckMpcAccountStep.HandleMessage begin")
	_, exist := ack.message[*msg.PeerID]
	if exist {
		log.SyslogErr("AckMpcGPKStep:HandleMessage","AckMpcAccountStep.HandleMessage fail. peer doesn't exist in task peer group. peerID",
			msg.PeerID.String())
		return false
	}

	if len(msg.BytesData) >= 1 {
		ack.remoteMpcGPKs[*msg.PeerID] = msg.BytesData[0]
	}

	ack.message[*msg.PeerID] = true
	return true
}
