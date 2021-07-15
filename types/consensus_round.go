package types

import (
	"fmt"

	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

const (
	DefaultConsensusRoundPeorid = 4
)

type ConsensusRound struct {
	ConsensusStartBlockHeight int64 `json:"consensus_start_block_height"`
	Peorid                    int64 `json:"peorid"`
}

func NewConsensusRound(consensusStartBlockHeight, peorid int64) *ConsensusRound {
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

	return &consensusRound
}

func (consensusRound ConsensusRound) ValidateBasic() error {
	if consensusRound.ConsensusStartBlockHeight < 0 {
		return fmt.Errorf("wrong ConsensusStartBlockHeight (got: %d)", consensusRound.ConsensusStartBlockHeight)
	}

	if consensusRound.Peorid < 0 {
		return fmt.Errorf("wrong consensus peorid (got: %d)", consensusRound.Peorid)
	}

	return nil
}

func (consensusRound ConsensusRound) ToProto() tmproto.ConsensusRound {
	return tmproto.ConsensusRound{
		ConsensusStartBlockHeight: consensusRound.ConsensusStartBlockHeight,
		Peorid:                    consensusRound.Peorid,
	}
}

func ConsensusRoundFromProto(crProto tmproto.ConsensusRound) *ConsensusRound {
	return &ConsensusRound{
		ConsensusStartBlockHeight: crProto.ConsensusStartBlockHeight,
		Peorid:                    crProto.Peorid,
	}
}
