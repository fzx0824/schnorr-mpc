package step

import (
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"math/big"
)

type MpcReadyStep struct {
	BaseStep
}

func (ready *MpcReadyStep) InitStep(result mpcprotocol.MpcResultInterface) error {
	ready.BaseStep.InitStep(result)
	return nil
}

func CreateMpcReadyStep(peers *[]mpcprotocol.PeerInfo) *MpcReadyStep {
	return &MpcReadyStep{
		*CreateBaseStep(peers, 0)}
}

func (ready *MpcReadyStep) CreateMessage() []mpcprotocol.StepMessage {
	data := make([]big.Int, 1)
	data[0].SetInt64(1)
	return []mpcprotocol.StepMessage{mpcprotocol.StepMessage{MsgCode: mpcprotocol.MPCMessage,
		PeerID:    nil,
		Peers:     nil,
		Data:      data,
		BytesData: nil}}
}

func (ready *MpcReadyStep) FinishStep(result mpcprotocol.MpcResultInterface, mpc mpcprotocol.StoremanManager) error {
	err := ready.BaseStep.FinishStep()
	if err != nil {
		return err
	}

	return nil
}

func (ready *MpcReadyStep) HandleMessage(msg *mpcprotocol.StepMessage) bool {
	return true
}
