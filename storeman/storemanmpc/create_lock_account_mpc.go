package storemanmpc

import (
	"github.com/wanchain/schnorr-mpc/log"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"github.com/wanchain/schnorr-mpc/storeman/storemanmpc/step"
)

//send create LockAccount from leader
func requestCreateLockAccountMpc(mpcID uint64, peers []mpcprotocol.PeerInfo, preSetValue ...MpcValue) (*MpcContext, error) {
	result := createMpcBaseMpcResult()
	result.InitializeValue(preSetValue...)
	mpc := createMpcContext(mpcID, peers, result)
	requestMpc := step.CreateRequestMpcStep(&mpc.peers, mpcprotocol.MpcCreateLockAccountLeader)
	mpcReady := step.CreateMpcReadyStep(&mpc.peers)
	return generateCreateLockAccountMpc(mpc, requestMpc, mpcReady)

}

//get message from leader and create Context
func acknowledgeCreateLockAccountMpc(mpcID uint64, peers []mpcprotocol.PeerInfo, preSetValue ...MpcValue) (*MpcContext, error) {
	log.SyslogInfo("acknowledgeCreateLockAccountMpc begin.")
	for _, preSetValuebyteData := range preSetValue {
		log.SyslogInfo("acknowledgeCreateLockAccountMpc", "byteValue", string(preSetValuebyteData.ByteValue[:]))
	}

	findMap := make(map[uint64]bool)
	for _, item := range peers {
		if item.Seed > 0xffffff {
			log.SyslogErr("acknowledgeCreateLockAccountMpc fail", "err", mpcprotocol.ErrMpcSeedOutRange.Error())
			return nil, mpcprotocol.ErrMpcSeedOutRange
		}

		_, exist := findMap[item.Seed]
		if exist {
			log.SyslogErr("acknowledgeCreateLockAccountMpc fail", "err", mpcprotocol.ErrMpcSeedDuplicate.Error())
			return nil, mpcprotocol.ErrMpcSeedDuplicate
		}

		findMap[item.Seed] = true
	}

	result := createMpcBaseMpcResult()
	result.InitializeValue(preSetValue...)
	mpc := createMpcContext(mpcID, peers, result)
	AcknowledgeMpc := step.CreateAcknowledgeMpcStep(&mpc.peers, mpcprotocol.MpcCreateLockAccountPeer)
	mpcReady := step.CreateGetMpcReadyStep(&mpc.peers)
	return generateCreateLockAccountMpc(mpc, AcknowledgeMpc, mpcReady)
}

func generateCreateLockAccountMpc(mpc *MpcContext, firstStep MpcStepFunc, readyStep MpcStepFunc) (*MpcContext, error) {
	var accTypeStr string
	accType, err := mpc.mpcResult.GetByteValue(mpcprotocol.MpcStmAccType)
	if err != nil {
		return nil, err
	} else if accType == nil {
		accTypeStr = ""
	} else {
		accTypeStr = string(accType[:])
	}

	JRSS := step.CreateMpcJRSS_Step(mpcprotocol.MPCDegree, &mpc.peers)
	PublicKey := step.CreateMpcAddressStep(&mpc.peers, accTypeStr)
	ackAddress := step.CreateAckMpcAccountStep(&mpc.peers)
	mpc.setMpcStep(firstStep, readyStep, JRSS, PublicKey, ackAddress)
	return mpc, nil
}
