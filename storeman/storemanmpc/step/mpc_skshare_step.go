package step

import (
	"crypto/ecdsa"
	"github.com/wanchain/schnorr-mpc/crypto"
	"github.com/wanchain/schnorr-mpc/log"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"math/big"
)

type MpcSKShareStep struct {
	BaseMpcStep
}

func CreateMpcSKShareStep(degree int, peers *[]mpcprotocol.PeerInfo) *MpcSKShareStep {
	mpc := &MpcSKShareStep{*CreateBaseMpcStep(peers, 1)}
	mpc.messages[0] = createSkPolyValue(degree, len(*peers))
	return mpc
}

func (jrss *MpcSKShareStep) CreateMessage() []mpcprotocol.StepMessage {
	message := make([]mpcprotocol.StepMessage, len(*jrss.peers))
	JRSSvalue := jrss.messages[0].(*RandomPolynomialValue)
	for i := 0; i < len(*jrss.peers); i++ {
		message[i].MsgCode = mpcprotocol.MPCMessage
		message[i].PeerID = &(*jrss.peers)[i].PeerID
		message[i].Data = make([]big.Int, 1)
		message[i].Data[0] = JRSSvalue.polyValue[i]
	}

	return message
}

func (jrss *MpcSKShareStep) FinishStep(result mpcprotocol.MpcResultInterface, mpc mpcprotocol.StoremanManager) error {
	err := jrss.BaseMpcStep.FinishStep()
	if err != nil {
		return err
	}

	// gskshare
	JRSSvalue := jrss.messages[0].(*RandomPolynomialValue)
	err = result.SetValue(mpcprotocol.MpcPrivateShare, []big.Int{*JRSSvalue.result})
	if err != nil {
		return err
	}
	// gpkshare
	var gpkShare ecdsa.PublicKey
	gpkShare.Curve = crypto.S256()
	gpkShare.X, gpkShare.Y = crypto.S256().ScalarBaseMult((*JRSSvalue.result).Bytes())
	err = result.SetValue(mpcprotocol.MpcPublicShare, []big.Int{*gpkShare.X, *gpkShare.Y})
	if err != nil {
		return err
	}

	return nil
}

func (jrss *MpcSKShareStep) HandleMessage(msg *mpcprotocol.StepMessage) bool {
	seed := jrss.getPeerSeed(msg.PeerID)
	log.Info("MpcSKShareStep::HandleMessage received message ",
		"peerId", msg.PeerID.String(),
		"seed", seed)
	if seed == 0 {
		log.SyslogErr("MpcJRSS_Step, can't find peer seed. peerID:%s", msg.PeerID.String())
	}

	JRSSvalue := jrss.messages[0].(*RandomPolynomialValue)
	_, exist := JRSSvalue.message[seed]
	if exist {
		log.SyslogErr("MpcJRSS_Step, can't find msg . peerID:%s, seed:%d", msg.PeerID.String(), seed)
		return false
	}

	JRSSvalue.message[seed] = msg.Data[0] //message.Value
	return true
}
