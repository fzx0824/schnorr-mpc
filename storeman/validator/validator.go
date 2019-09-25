package validator

import (
	"encoding/json"
	"errors"
	"github.com/wanchain/schnorr-mpc/common"
	"github.com/wanchain/schnorr-mpc/crypto"
	"github.com/wanchain/schnorr-mpc/log"
	mpcprotocol "github.com/wanchain/schnorr-mpc/storeman/storemanmpc/protocol"
	"time"
)

var noticeFuncIds [][4]byte

func init() {

}

// TODO add ValidateData
func ValidateData(data []byte) bool {
	// 1. check in local db or not
	// 2. check has approved or not
	return true
}

//func ValidateTx(signer mpccrypto.MPCTxSigner, from common.Address, chainType string, chainId *big.Int, leaderTxRawData []byte, leaderTxLeaderHashBytes []byte) bool {
//	log.SyslogInfo("ValidateTx",
//		"from", from.String(),
//		"chainType", chainType,
//		"chainId", chainId.String(),
//		"leaderTxLeaderHashBytes", common.ToHex(leaderTxLeaderHashBytes),
//		"leaderTxRawData", common.ToHex(leaderTxRawData))
//
//	var leaderTx types.Transaction
//	err := rlp.DecodeBytes(leaderTxRawData, &leaderTx)
//	if err != nil {
//		log.SyslogErr("ValidateTx leader tx data decode fail", "err", err.Error())
//		return false
//	}
//
//	log.SyslogInfo("ValidateTx", "leaderTxData", common.ToHex(leaderTx.Data()))
//	isNotice, err := IsNoticeTransaction(leaderTx.Data())
//	if err != nil {
//		log.SyslogErr("ValidateTx, check notice transaction fail", "err", err.Error())
//	} else if isNotice {
//		log.SyslogInfo("ValidateTx, is notice transaction, skip validating")
//		return true
//	}
//
//	key := GetKeyFromTx(&from, leaderTx.To(), leaderTx.Value(), leaderTx.Data(), &chainType, chainId)
//	log.SyslogInfo("mpc ValidateTx", "key", common.ToHex(key))
//
//	followerDB, err := GetDB()
//	if err != nil {
//		log.SyslogErr("ValidateTx leader get database fail", "err", err.Error())
//		return false
//	}
//
//	_, err = waitKeyFromDB([][]byte{key})
//	if err != nil {
//		log.SyslogErr("ValidateTx, check has fail", "err", err.Error())
//		return false
//	}
//
//	followerTxRawData, err := followerDB.Get(key)
//	if err != nil {
//		log.SyslogErr("ValidateTx, getting followerTxRawData fail", "err", err.Error())
//		return false
//	}
//
//	log.SyslogInfo("ValidateTx, followerTxRawData is got")
//
//	var followerRawTx mpcprotocol.SendTxArgs
//	err = json.Unmarshal(followerTxRawData, &followerRawTx)
//	if err != nil {
//		log.SyslogErr("ValidateTx, follower tx data decode fail", "err", err.Error())
//		return false
//	}
//
//	followerCreatedTx := types.NewTransaction(leaderTx.Nonce(), *followerRawTx.To, followerRawTx.Value.ToInt(),
//		leaderTx.Gas(), leaderTx.GasPrice(), followerRawTx.Data)
//	followerCreatedHash := signer.Hash(followerCreatedTx)
//	leaderTxLeaderHash := common.BytesToHash(leaderTxLeaderHashBytes)
//
//	if followerCreatedHash == leaderTxLeaderHash {
//		log.SyslogInfo("ValidateTx, validate success")
//		return true
//	} else {
//		log.SyslogErr("ValidateTx, leader tx hash is not same with follower tx hash",
//			"leaderTxLeaderHash", leaderTxLeaderHash.String(),
//			"followerCreatedHash", followerCreatedHash.String())
//		return false
//	}
//}

func waitKeyFromDB(keys [][]byte) ([]byte, error) {
	log.SyslogInfo("waitKeyFromDB, begin")

	for i, key := range keys {
		log.SyslogInfo("waitKeyFromDB", "i", i, "key", common.ToHex(key))
	}

	db, err := GetDB()
	if err != nil {
		log.SyslogErr("waitKeyFromDB get database fail", "err", err.Error())
		return nil, err
	}

	start := time.Now()
	for {
		for _, key := range keys {
			isExist, err := db.Has(key)
			if err != nil {
				log.SyslogErr("waitKeyFromDB fail", "err", err.Error())
				return nil, err
			} else if isExist {
				log.SyslogInfo("waitKeyFromDB, got it", "key", common.ToHex(key))
				return key, nil
			}
		}

		if time.Now().Sub(start) >= mpcprotocol.MPCTimeOut {
			log.SyslogInfo("waitKeyFromDB, time out")
			return nil, errors.New("waitKeyFromDB, time out")
		}

		time.Sleep(200 * time.Microsecond)
	}

	return nil, errors.New("waitKeyFromDB, unknown error")
}

func AddValidData(data *mpcprotocol.SendData) error {
	log.SyslogInfo("AddValidData", "data", data.String())
	val, err := json.Marshal(&data)
	if err != nil {
		log.SyslogErr("AddValidData, marshal fail", "err", err.Error())
		return err
	}

	key := crypto.Keccak256(data.Data[:])
	return addKeyValueToDB(key, val)
}

func addKeyValueToDB(key, value []byte) error {
	log.SyslogInfo("addKeyValueToDB, begin", "key:", common.ToHex(key))
	sdb, err := GetDB()
	if err != nil {
		log.SyslogErr("addKeyValueToDB, getting storeman database fail", "err", err.Error())
		return err
	}

	err = sdb.Put(key, value)
	if err != nil {
		log.SyslogErr("addKeyValueToDB, getting storeman database fail", "err", err.Error())
		return err
	}

	log.SyslogInfo("addKeyValueToDB", "key", common.ToHex(key))
	ret, err := sdb.Get(key)
	if err != nil {
		log.SyslogErr("addKeyValueToDB, getting storeman database fail", "err", err.Error())
		return err
	}

	log.SyslogInfo("addKeyValueToDB succeed to get data from level db after putting key-val pair", "ret", string(ret))
	return nil
}
