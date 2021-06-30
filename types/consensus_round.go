package types

import (
	"fmt"

	tmproto "gitlab.reappay.net/sucs-lab/reapchain/proto/reapchain/types"
)

const (
	DefaultConsensusRoundPeorid = 4
)

type ConsensusRound struct {
	ConsensusStartBlockHeight int64 `json:"consensus_start_block_height"`
	Peorid                    int64 `json:"peorid"`
}

func NewConsensusRound(consensusStartBlockHeight, peorid int64) ConsensusRound {
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

func (consensusRoundInfo ConsensusRound) ToProto() tmproto.ConsensusRound {
	return tmproto.ConsensusRound{
		ConsensusStartBlockHeight: consensusRoundInfo.ConsensusStartBlockHeight,
		Peorid:                    consensusRoundInfo.Peorid,
	}
}

func ConsensusRoundFromProto(crProto tmproto.ConsensusRound) ConsensusRound {
	return ConsensusRound{
		ConsensusStartBlockHeight: crProto.ConsensusStartBlockHeight,
		Peorid:                    crProto.Peorid,
	}
}
