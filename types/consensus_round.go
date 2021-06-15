package types

import ()

const (
	DefaultConsensusRoundPeorid = 4
)
type ConsensusRound struct {
	ConsensusStartBlockHeight   int64         `json:"consensus_start_block_height"`
	Peorid 											int64         `json:"peorid"`
}

func NewConsensusRound(consensusStartBlockHeight, peorid int64) ConsensusRound {
	consensusRound := ConsensusRound{
		ConsensusStartBlockHeight: consensusStartBlockHeight, 
		Peorid: peorid,
	}

	if(consensusRound.ConsensusStartBlockHeight < 0) {
		consensusRound.ConsensusStartBlockHeight = 0
	}

	if(consensusRound.Peorid == 0) {
		consensusRound.Peorid = DefaultConsensusRoundPeorid
	}
	
	return consensusRound
}