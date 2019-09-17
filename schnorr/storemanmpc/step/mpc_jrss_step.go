package step

import (
	mpcprotocol "github.com/wanchain/schnorr-mpc/schnorr/storemanmpc/protocol"
	"github.com/wanchain/schnorr-mpc/log"
	"math/big"
)

type MpcJRSS_Step struct {
	BaseMpcStep
}

func CreateMpcJRSS_Step(degree int, peers *[]mpcprotocol.PeerInfo) *MpcJRSS_Step {
	mpc := &MpcJRSS_Step{*CreateBaseMpcStep(peers, 1)}
	mpc.messages[0] = createJRSSValue(degree, len(*peers))
	return mpc
}

func (jrss *MpcJRSS_Step) CreateMessage() []mpcprotocol.StepMessage {
	message := make([]mpcprotocol.StepMessage, len(*jrss.peers))
	JRSSvalue := jrss.messages[0].(*RandomPolynomialValue)
	for i := 0; i < len(*jrss.peers); i++ {
		message[i].Msgcode = mpcprotocol.MPCMessage
		message[i].PeerID = &(*jrss.peers)[i].PeerID
		message[i].Data = make([]big.Int, 1)
		message[i].Data[0] = JRSSvalue.polyValue[i]
	}

	return message
}

func (jrss *MpcJRSS_Step) FinishStep(result mpcprotocol.MpcResultInterface, mpc mpcprotocol.StoremanManager) error {
	err := jrss.BaseMpcStep.FinishStep()
	if err != nil {
		return err
	}

	JRSSvalue := jrss.messages[0].(*RandomPolynomialValue)
	err = result.SetValue(mpcprotocol.MpcPrivateShare, []big.Int{*JRSSvalue.result})
	if err != nil {
		return err
	}

	err = result.SetValue(mpcprotocol.MpcPublicShare, []big.Int{JRSSvalue.randCoefficient[0]})
	if err != nil {
		return err
	}

	return nil
}

func (jrss *MpcJRSS_Step) HandleMessage(msg *mpcprotocol.StepMessage) bool {
	seed := jrss.getPeerSeed(msg.PeerID)
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
