package store

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	DelegateType   = 0
	UndelegateType = iota
)

type Staker struct {
	StakingAmount uint64 `json:"stakingAmount"`
	Address       string `json:"address"`
}

type QuorumNode struct {
	FpBtcPk      string   `json:"fpBtcPk"`
	FpVoteWeight uint64   `json:"fpVoteWeight"`
	IsSign       bool     `json:"isSign"`
	Staker       []Staker `json:"staker"`
}

type StakeDetails struct {
	BatchID           uint64       `json:"batchId"`
	TotalBTCVote      uint64       `json:"totalBtcVote"`
	BabylonBlock      uint64       `json:"babylonBlock"`
	StateRoot         string       `json:"stateRoot"`
	EthBlock          uint64       `json:"ethBlock"`
	BitcoinQuorum     []QuorumNode `json:"bitcoinQuorum"`
	SymbioticSignNode []string     `json:"symbioticSignNode"`
}

func (s *Storage) SetStakeDetails(msg CreateBTCDelegation, stakeType int8) error {
	if stakeType == DelegateType {
		stakeDB, err := s.db.Get(getStakeDetailsKey(), nil)
		if err != nil {
			if errors.Is(err, leveldb.ErrNotFound) {
				var sD = StakeDetails{
					TotalBTCVote: uint64(msg.CBD.StakingValue),
					BitcoinQuorum: []QuorumNode{
						{
							FpBtcPk:      msg.CBD.FpBtcPkList[0].MarshalHex(),
							FpVoteWeight: uint64(msg.CBD.StakingValue),
							Staker: []Staker{
								{
									StakingAmount: uint64(msg.CBD.StakingValue),
									Address:       msg.CBD.StakerAddr,
								},
							},
						},
					},
				}
				bz, err := json.Marshal(sD)
				if err != nil {
					return err
				}
				return s.db.Put(getStakeDetailsKey(), bz, nil)
			} else {
				return err
			}
		}

		var sD StakeDetails
		if err = json.Unmarshal(stakeDB, &sD); err != nil {
			return err
		}
		var existFpBtcPk bool
		var existStaker bool
		for _, quorum := range sD.BitcoinQuorum {
			if quorum.FpBtcPk == msg.CBD.FpBtcPkList[0].MarshalHex() {
				quorum.FpVoteWeight += uint64(msg.CBD.StakingValue)
				existFpBtcPk = true
				for _, staker := range quorum.Staker {
					if staker.Address == msg.CBD.StakerAddr {
						existStaker = true
						staker.StakingAmount += uint64(msg.CBD.StakingValue)
					}
				}
			}
		}

		if !existFpBtcPk {
			var quorum = QuorumNode{
				FpBtcPk:      msg.CBD.FpBtcPkList[0].MarshalHex(),
				FpVoteWeight: uint64(msg.CBD.StakingValue),
				Staker: []Staker{
					{
						StakingAmount: uint64(msg.CBD.StakingValue),
						Address:       msg.CBD.StakerAddr,
					},
				},
			}
			sD.BitcoinQuorum = append(sD.BitcoinQuorum, quorum)
		} else {
			if !existStaker {
				for _, quorum := range sD.BitcoinQuorum {
					if quorum.FpBtcPk == msg.CBD.FpBtcPkList[0].MarshalHex() {
						var staker = Staker{
							StakingAmount: uint64(msg.CBD.StakingValue),
							Address:       msg.CBD.StakerAddr,
						}
						quorum.Staker = append(quorum.Staker, staker)
					}
				}
			}
		}

		bsD, err := json.Marshal(sD)
		if err != nil {
			return err
		}
		return s.db.Put(getStakeDetailsKey(), bsD, nil)
	} else if stakeType == UndelegateType {
		stakeDB, err := s.db.Get(getStakeDetailsKey(), nil)
		if err != nil {
			return err
		}
		var sD StakeDetails
		if err = json.Unmarshal(stakeDB, &sD); err != nil {
			return err
		}
		for _, quorum := range sD.BitcoinQuorum {
			if quorum.FpBtcPk == msg.CBD.FpBtcPkList[0].MarshalHex() {
				for _, staker := range quorum.Staker {
					staker.StakingAmount -= uint64(msg.CBD.StakingValue)
				}
			}
		}
		bsD, err := json.Marshal(sD)
		if err != nil {
			return err
		}
		return s.db.Put(getStakeDetailsKey(), bsD, nil)
	} else {
		return errors.New("unknown stake type")
	}
}

func (s *Storage) GetStakeDetails() (StakeDetails, error) {
	sDB, err := s.db.Get(getStakeDetailsKey(), nil)
	if err != nil {
		return handleError(StakeDetails{}, err)
	}
	var sD StakeDetails
	if err = json.Unmarshal(sDB, &sD); err != nil {
		return StakeDetails{}, err
	}
	return sD, nil
}

func (s *Storage) SetBatchStakeDetails(batchID uint64, fpSignCache map[string]string, stateRoot string, babylonBlockHeight uint64, ethBlockHeight uint64) error {
	sDB, err := s.db.Get(getStakeDetailsKey(), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return errors.New("the database does not have stake data")
		}
		return err
	}

	var sD StakeDetails
	if err = json.Unmarshal(sDB, &sD); err != nil {
		return err
	}

	var sF SymbioticFpIds
	sFB, err := s.db.Get(getSymbioticFpIdsKey(batchID), nil)
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return err
		}
	} else {
		if err = json.Unmarshal(sFB, &sF); err != nil {
			return err
		}
		for _, sR := range sF.SignRequests {
			sD.SymbioticSignNode = append(sD.SymbioticSignNode, sR.SignAddress)
		}
	}

	sD.BabylonBlock = babylonBlockHeight
	sD.StateRoot = stateRoot
	sD.EthBlock = ethBlockHeight

	for fpPubkeyHex, sR := range fpSignCache {
		for i, quorum := range sD.BitcoinQuorum {
			if strings.ToLower(fpPubkeyHex) == strings.ToLower(quorum.FpBtcPk) && sR == stateRoot {
				sD.BitcoinQuorum[i].IsSign = true
			}
		}
	}

	bsD, err := json.Marshal(sD)
	if err != nil {
		return err
	}

	return s.db.Put(getBatchStakeDetailsKey(batchID), bsD, nil)
}

func (s *Storage) GetBatchStakeDetails(batchID uint64) (StakeDetails, error) {
	bSDB, err := s.db.Get(getBatchStakeDetailsKey(batchID), nil)
	if err != nil {
		return handleError(StakeDetails{}, err)
	}

	var sD StakeDetails
	if err = json.Unmarshal(bSDB, &sD); err != nil {
		return StakeDetails{}, err
	}
	sD.BatchID = batchID

	return sD, nil
}

func (s *Storage) GetBatchTotalBabylonStakeAmount(batchID uint64) (uint64, error) {
	bSDB, err := s.db.Get(getBatchStakeDetailsKey(batchID), nil)
	if err != nil {
		return 0, err
	}

	var sD StakeDetails
	if err = json.Unmarshal(bSDB, &sD); err != nil {
		return 0, err
	}

	var totalAmount uint64
	for _, quorum := range sD.BitcoinQuorum {
		if quorum.IsSign {
			for _, staker := range quorum.Staker {
				totalAmount += staker.StakingAmount
			}
		}
	}

	return totalAmount, nil
}
