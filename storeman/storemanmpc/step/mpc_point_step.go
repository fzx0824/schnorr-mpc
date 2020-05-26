package step

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/wanchain/schnorr-mpc/common/hexutil"
	"github.com/wanchain/schnorr-mpc/crypto"
	"github.com/wanchain/schnorr-mpc/log"
	"github.com/wanchain/schnorr-mpc/storeman/osmconf"
	"github.com/wanchain/schnorr-mpc/storeman/schnorrmpc"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"math/big"
	"strconv"
)

type MpcPointStep struct {
	BaseMpcStep
	resultKeys []string
	signNum    int

	RShareErrNum    int
	rpkshareOKIndex []uint16
	rpkshareKOIndex []uint16
	rpkshareNOIndex []uint16
}

func CreateMpcPointStep(peers *[]mpcprotocol.PeerInfo, preValueKeys []string, resultKeys []string) *MpcPointStep {
	log.SyslogInfo("CreateMpcPointStep begin")

	signNum := len(preValueKeys)
	mpc := &MpcPointStep{*CreateBaseMpcStep(peers, signNum), resultKeys, signNum, 0,
		make([]uint16, 0), make([]uint16, 0), make([]uint16, 0)}

	for i := 0; i < signNum; i++ {
		mpc.messages[i] = createPointGenerator(preValueKeys[i])
	}

	return mpc
}

func (ptStep *MpcPointStep) CreateMessage() []mpcprotocol.StepMessage {
	log.SyslogInfo("MpcPointStep.CreateMessage begin")
	message := make([]mpcprotocol.StepMessage, 1)
	message[0].MsgCode = mpcprotocol.MPCMessage
	message[0].PeerID = nil

	pointer := ptStep.messages[0].(*mpcPointGenerator)
	var buf bytes.Buffer
	buf.Write(crypto.FromECDSAPub(&pointer.seed))
	h := sha256.Sum256(buf.Bytes())

	prv, _ := osmconf.GetOsmConf().GetSelfPrvKey()
	r, s, _ := schnorrmpc.SignInternalData(prv, h[:])
	// send rpkshare, sig of rpkshare
	message[0].Data = make([]big.Int, 2)
	message[0].Data[0] = *r
	message[0].Data[1] = *s

	// only one point, rpkShare
	message[0].BytesData = append(message[0].BytesData, crypto.FromECDSAPub(&pointer.seed))

	return message
}

func (ptStep *MpcPointStep) HandleMessage(msg *mpcprotocol.StepMessage) bool {

	log.SyslogInfo("MpcPointStep.HandleMessage begin ",
		"peerID", msg.PeerID.String(),
		"gpk x", hex.EncodeToString(msg.BytesData[0]))

	pointer := ptStep.messages[0].(*mpcPointGenerator)
	_, exist := pointer.message[*msg.PeerID]
	if exist {
		log.SyslogErr("HandleMessage", "MpcPointStep.HandleMessage, get msg from seed fail. peer", msg.PeerID.String())
		return false
	}

	pointPk := crypto.ToECDSAPub(msg.BytesData[0])
	r := msg.Data[0]
	s := msg.Data[1]

	_, grpIdString, _ := osmconf.GetGrpId(ptStep.mpcResult)

	senderPk, _ := osmconf.GetOsmConf().GetPKByNodeId(grpIdString, msg.PeerID)
	err := schnorrmpc.CheckPK(senderPk)
	if err != nil {
		log.SyslogErr("MpcPointStep", "HandleMessage", err.Error())
	}

	senderIndex, _ := osmconf.GetOsmConf().GetInxByNodeId(grpIdString, msg.PeerID)

	var buf bytes.Buffer
	buf.Write(msg.BytesData[0])
	h := sha256.Sum256(buf.Bytes())

	bVerifySig := schnorrmpc.VerifyInternalData(senderPk, h[:], &r, &s)

	if bVerifySig {
		pointer.message[*msg.PeerID] = *pointPk

		// save rpkshare for check data of s
		key := mpcprotocol.RPkShare + strconv.Itoa(int(senderIndex))
		ptStep.mpcResult.SetByteValue(key, msg.BytesData[0])
		log.SyslogInfo("@@@@@@@@@@@@save rpkshare", "key", key, "rpkshare", hexutil.Encode(msg.BytesData[0]))

		ptStep.rpkshareOKIndex = append(ptStep.rpkshareOKIndex, uint16(senderIndex))
	} else {
		log.SyslogErr("MpcPointStep::HandleMessage", " check sig fail")
		ptStep.rpkshareKOIndex = append(ptStep.rpkshareKOIndex, uint16(senderIndex))
	}

	// update no work indexes.

	return true
}

func (ptStep *MpcPointStep) FinishStep(result mpcprotocol.MpcResultInterface, mpc mpcprotocol.StoremanManager) error {
	log.SyslogInfo("MpcPointStep.FinishStep begin")

	//grpId,_ := ptStep.mpcResult.GetByteValue(mpcprotocol.MpcGrpId)
	//grpIdString := string(grpId)
	//allIndex,_ := osmconf.GetOsmConf().GetGrpElemsInxes(grpIdString)
	//tempIndex := osmconf.Difference(*allIndex,ptStep.rpkshareOKIndex)
	//ptStep.rpkshareNOIndex = osmconf.Difference(tempIndex,ptStep.rpkshareKOIndex)
	//
	//log.Info(">>>>>>MpcPointStep", "allIndex", allIndex)
	//log.Info(">>>>>>MpcPointStep", "rpkshareOKIndex", ptStep.rpkshareOKIndex)
	//log.Info(">>>>>>MpcPointStep", "rpkshareKOIndex", ptStep.rpkshareKOIndex)
	//log.Info(">>>>>>MpcPointStep", "rpkshareNOIndex", ptStep.rpkshareNOIndex)
	//
	//okIndex := make([]big.Int,len(ptStep.rpkshareOKIndex))
	//koIndex := make([]big.Int,len(ptStep.rpkshareKOIndex))
	//noIndex := make([]big.Int,len(ptStep.rpkshareNOIndex))
	//
	//for i,value := range ptStep.rpkshareOKIndex{
	//	okIndex[i].SetInt64(int64(value))
	//}
	//
	//for i,value := range ptStep.rpkshareKOIndex{
	//	koIndex[i].SetInt64(int64(value))
	//}
	//
	//for i,value := range ptStep.rpkshareNOIndex{
	//	noIndex[i].SetInt64(int64(value))
	//}
	//
	//ptStep.mpcResult.SetValue(mpcprotocol.ROKIndex,okIndex)
	//ptStep.mpcResult.SetValue(mpcprotocol.RKOIndex,koIndex)
	//ptStep.mpcResult.SetValue(mpcprotocol.RNOIndex,noIndex)

	err := ptStep.BaseMpcStep.FinishStep()
	if err != nil {
		return err
	}

	pointer := ptStep.messages[0].(*mpcPointGenerator)
	log.Info("generated gpk MpcPointStep::FinishStep",
		"result key", ptStep.resultKeys[0],
		"result value ", hexutil.Encode(crypto.FromECDSAPub(&pointer.result)))

	//err = result.SetByteValue(ptStep.resultKeys[0], crypto.FromECDSAPub(&pointer.result))
	err = result.SetByteValue(mpcprotocol.RPk, crypto.FromECDSAPub(&pointer.result))
	if err != nil {
		log.SyslogErr("HandleMessage", "MpcPointStep.FinishStep, SetValue fail. err", err.Error())
		return err
	}

	log.SyslogInfo("@@@@@@@@@@@@save RPK", "key", mpcprotocol.RPk, "RPK", hexutil.Encode(crypto.FromECDSAPub(&pointer.result)))
	log.SyslogInfo("MpcPointStep.FinishStep succeed")
	return nil
}
