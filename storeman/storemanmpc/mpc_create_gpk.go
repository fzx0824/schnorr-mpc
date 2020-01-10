package storemanmpc

import (
	"github.com/wanchain/schnorr-mpc/log"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"github.com/wanchain/schnorr-mpc/storeman/storemanmpc/step"
)

//send create LockAccount from leader
func reqGPKMpc(mpcID uint64, peers []mpcprotocol.PeerInfo, preSetValue ...MpcValue) (*MpcContext, error) {

	result := createMpcBaseMpcResult()
	result.InitializeValue(preSetValue...)
	mpc := createMpcContext(mpcID, peers, result)
	reqMpc := step.CreateRequestMpcStep(&mpc.peers, mpcprotocol.MpcGPKLeader)
	mpcReady := step.CreateMpcReadyStep(&mpc.peers)
	return genCreateGPKMpc(mpc, reqMpc, mpcReady)

}

//get message from leader and create Context
func ackGPKMpc(mpcID uint64, peers []mpcprotocol.PeerInfo, preSetValue ...MpcValue) (*MpcContext, error) {

	log.SyslogInfo("ackGPKMpc begin.")
	for _, preSetValuebyteData := range preSetValue {
		log.SyslogInfo("ackGPKMpc", "byteValue", string(preSetValuebyteData.ByteValue[:]))
	}

	findMap := make(map[uint64]bool)
	for _, item := range peers {
		if item.Seed > 0xffffff {
			log.SyslogErr("ackGPKMpc fail", "err", mpcprotocol.ErrMpcSeedOutRange.Error())
			return nil, mpcprotocol.ErrMpcSeedOutRange
		}

		_, exist := findMap[item.Seed]
		if exist {
			log.SyslogErr("ackGPKMpc fail", "err", mpcprotocol.ErrMpcSeedDuplicate.Error())
			return nil, mpcprotocol.ErrMpcSeedDuplicate
		}

		findMap[item.Seed] = true
	}

	result := createMpcBaseMpcResult()
	result.InitializeValue(preSetValue...)
	mpc := createMpcContext(mpcID, peers, result)
	AckMpc := step.CreateAckMpcStep(&mpc.peers, mpcprotocol.MpcGPKPeer)
	mpcReady := step.CreateGetMpcReadyStep(&mpc.peers)
	return genCreateGPKMpc(mpc, AckMpc, mpcReady)
}

func genCreateGPKMpc(mpc *MpcContext, firstStep MpcStepFunc, readyStep MpcStepFunc) (*MpcContext, error) {

	accTypeStr := ""
	skShare := step.CreateMpcSKShareStep(mpcprotocol.MPCDegree, &mpc.peers)
	gpk := step.CreateMpcGPKStep(&mpc.peers, accTypeStr)
	ackGpk := step.CreateAckMpcGPKStep(&mpc.peers)
	mpc.setMpcStep(firstStep, readyStep, skShare, gpk, ackGpk)

	for stepId, stepItem := range mpc.MpcSteps {
		stepItem.SetStepId(stepId)
	}

	return mpc, nil

}
