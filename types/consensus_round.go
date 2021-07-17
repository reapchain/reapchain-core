package types

import (
	"fmt"

	abci "github.com/reapchain/reapchain-core/abci/types"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

const (
	DefaultConsensusRoundPeorid = 4
)

type ConsensusRound struct {
	ConsensusStartBlockHeight int64  `json:"consensus_start_block_height"`
	Peorid                    uint64 `json:"peorid"`
}

func NewConsensusRound(consensusStartBlockHeight int64, peorid uint64) ConsensusRound {
	consensusRound := ConsensusRound{
		ConsensusStartBlockHeight: consensusStartBlockHeight,
		Peorid:                    peorid,
	}

	if consensusRound.ConsensusStartBlockHeight < 0 {
		consensusRound.ConsensusStartBlockHeight = 0
	}

	if consensusRound.Peorid == 0 {
		consensusRound.Peorid = DefaultConsensusRoundPeorid
	}

	return consensusRound
}

func (consensusRoundInfo ConsensusRound) ValidateBasic() error {
	if consensusRoundInfo.ConsensusStartBlockHeight < 0 {
		return fmt.Errorf("wrong ConsensusStartBlockHeight (got: %d)", consensusRoundInfo.ConsensusStartBlockHeight)
	}

	if consensusRoundInfo.Peorid < 0 {
		return fmt.Errorf("wrong consensus peorid (got: %d)", consensusRoundInfo.Peorid)
	}

	return nil
}

func (consensusRound *ConsensusRound) ToProto() tmproto.ConsensusRound {
	return tmproto.ConsensusRound{
		ConsensusStartBlockHeight: consensusRound.ConsensusStartBlockHeight,
		Peorid:                    consensusRound.Peorid,
	}
}

func ConsensusRoundFromProto(crProto tmproto.ConsensusRound) ConsensusRound {
	return ConsensusRound{
		ConsensusStartBlockHeight: crProto.ConsensusStartBlockHeight,
		Peorid:                    crProto.Peorid,
	}
}

func UpdateConsensusRound(currentConsensusRound tmproto.ConsensusRound, nextConsensusRound *abci.ConsensusRound) tmproto.ConsensusRound {
	res := currentConsensusRound // explicit copy

	if nextConsensusRound == nil {
		return res
	}

	res.Peorid = nextConsensusRound.Peorid

	return res
}
