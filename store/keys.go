package store

import (
	"encoding/binary"
)

var (
	BabylonBlockHeaderKeyPrefix        = []byte{0x01}
	EthBlockHeaderKeyPrefix            = []byte{0x02}
	TxMessageKeyPrefix                 = []byte{0x03}
	EthScannedHeightKeyPrefix          = []byte{0x04}
	BabylonScannedHeightKeyPrefix      = []byte{0x05}
	NewFinalityProviderKeyPrefix       = []byte{0x06}
	CreateBTCDelegationKeyPrefix       = []byte{0x07}
	CommitPubRandListKeyPrefix         = []byte{0x08}
	SignatureKeyPrefix                 = []byte{0x09}
	ContractEventKeyPrefix             = []byte{0x10}
	ActiveMemberKeyPrefix              = []byte{0x11}
	BtcUndelegateKeyPrefix             = []byte{0x12}
	BTCDelegateAmountKeyPrefix         = []byte{0x13}
	SelectiveSlashingEvidenceKeyPrefix = []byte{0x14}
	BabylonDelegationKeyPrefix         = []byte{0x15}
	SubmitFinalitySignatureKeyPrefix   = []byte{0x16}
	StakeDetailsKeyPrefix              = []byte{0x17}
	BatchStakeDetailsKeyPrefix         = []byte{0x18}
	SymbioticFpIdsKeyPrefix            = []byte{0x19}
	L2OutputOracleKeyPrefix            = []byte{0x20}
)

func getBabylonBlockHeaderKey(number int64) []byte {
	numberBz := make([]byte, 8)
	binary.BigEndian.PutUint64(numberBz, uint64(number))
	return append(BabylonBlockHeaderKeyPrefix, numberBz...)
}

func getEthBlockHeaderKey(number int64) []byte {
	numberBz := make([]byte, 8)
	binary.BigEndian.PutUint64(numberBz, uint64(number))
	return append(EthBlockHeaderKeyPrefix, numberBz...)
}

func getSignatureKey(number int64) []byte {
	numberBz := make([]byte, 8)
	binary.BigEndian.PutUint64(numberBz, uint64(number))
	return append(SignatureKeyPrefix, numberBz...)
}

func getTxMessageKey(txHash []byte) []byte {
	return append(TxMessageKeyPrefix, txHash[:]...)
}

func getContractEventKey(txHash []byte) []byte {
	return append(ContractEventKeyPrefix, txHash[:]...)
}

func getNewFinalityProviderKey(txHash []byte) []byte {
	return append(NewFinalityProviderKeyPrefix, txHash[:]...)
}

func getCreateBTCDelegationKey(txHash []byte) []byte {
	return append(CreateBTCDelegationKeyPrefix, txHash[:]...)
}

func getCommitPubRandListKey(txHash []byte) []byte {
	return append(CommitPubRandListKeyPrefix, txHash[:]...)
}

func getBtcUndelegateKey(txHash []byte) []byte {
	return append(BtcUndelegateKeyPrefix, txHash[:]...)
}

func getSelectiveSlashingEvidenceKey(txHash []byte) []byte {
	return append(SelectiveSlashingEvidenceKeyPrefix, txHash[:]...)
}

func getBTCDelegateAmountKey() []byte {
	return append(BTCDelegateAmountKeyPrefix)
}

func getBabylonDelegationKey(btcTx []byte) []byte {
	return append(BabylonDelegationKeyPrefix, btcTx[:]...)
}

func getEthScannedHeightKey() []byte {
	return EthScannedHeightKeyPrefix
}

func getBabylonScannedHeightKey() []byte {
	return BabylonScannedHeightKeyPrefix
}

func getActiveMemberKey() []byte {
	return ActiveMemberKeyPrefix
}

func getSubmitFinalitySignatureKey(txHash []byte) []byte {
	return append(SubmitFinalitySignatureKeyPrefix, txHash[:]...)
}

func getStakeDetailsKey() []byte {
	return StakeDetailsKeyPrefix
}

func getBatchStakeDetailsKey(batchId uint64) []byte {
	numberBz := make([]byte, 8)
	binary.BigEndian.PutUint64(numberBz, batchId)
	return append(BatchStakeDetailsKeyPrefix, numberBz...)
}

func getSymbioticFpIdsKey(batchId uint64) []byte {
	numberBz := make([]byte, 8)
	binary.BigEndian.PutUint64(numberBz, batchId)
	return append(SymbioticFpIdsKeyPrefix, numberBz...)
}

func getL2OutputOracleKey(l1Number uint64) []byte {
	numberBz := make([]byte, 8)
	binary.BigEndian.PutUint64(numberBz, l1Number)
	return append(L2OutputOracleKeyPrefix, numberBz...)
}
