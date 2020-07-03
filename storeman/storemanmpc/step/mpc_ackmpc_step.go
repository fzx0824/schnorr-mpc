package step

import (
	"github.com/wanchain/schnorr-mpc/log"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"math/big"
)

type AckMpcStep struct {
	BaseStep
	messageType int64
}

func CreateAckMpcStep(peers *[]mpcprotocol.PeerInfo, messageType int64) *AckMpcStep {
	log.SyslogInfo("CreateAcknowledgeMpcStep begin")

	return &AckMpcStep{
		*CreateBaseStep(peers, 0), messageType}
}

func (ack *AckMpcStep) InitStep(result mpcprotocol.MpcResultInterface) error {
	ack.BaseStep.InitStep(result)
	return nil
}

func (ack *AckMpcStep) CreateMessage() []mpcprotocol.StepMessage {
	log.SyslogInfo("AcknowledgeMpcStep.CreateMessage begin")

	data := make([]big.Int, 1)
	data[0].SetInt64(ack.messageType)
	return []mpcprotocol.StepMessage{mpcprotocol.StepMessage{
		MsgCode:   mpcprotocol.MPCMessage,
		PeerID:    nil,
		Peers:     nil,
		Data:      data,
		BytesData: nil}}
}

func (ack *AckMpcStep) FinishStep(result mpcprotocol.MpcResultInterface, mpc mpcprotocol.StoremanManager) error {
	log.SyslogInfo("AcknowledgeMpcStep.FinishStep begin")

	err := ack.BaseStep.FinishStep()
	if err != nil {
		return err
	}

	data := make([]big.Int, 1)
	data[0].SetInt64(ack.messageType)
	result.SetValue(mpcprotocol.MPCAction, data)

	log.SyslogInfo("AcknowledgeMpcStep.FinishStep succeed")
	return nil
}

func (ack *AckMpcStep) HandleMessage(msg *mpcprotocol.StepMessage) bool {
	return true
}
