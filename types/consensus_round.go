package types

import (
	"fmt"

	abci "github.com/reapchain/reapchain-core/abci/types"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

const (
	DefaultConsensusRoundQrnPeorid       = 4
	DefaultConsensusRoundVrfPeorid       = 4
	DefaultConsensusRoundValidatorPeorid = 4
	DefaultConsensusRoundPeorid          = 12
)

type ConsensusRound struct {
	ConsensusStartBlockHeight int64  `json:"consensus_start_block_height"`
	Peorid                    uint64 `json:"peorid"`
	QrnPeorid                 uint64 `json:"qrn_peorid"`
	VrfPeorid                 uint64 `json:"vrf_peorid"`
	ValidatorPeorid           uint64 `json:"validator_peorid"`
}

func NewConsensusRound(consensusStartBlockHeight int64, qrnPeorid, vrfPeorid, validatorPeorid uint64) ConsensusRound {
	consensusRound := ConsensusRound{
		ConsensusStartBlockHeight: consensusStartBlockHeight,
		QrnPeorid:                 qrnPeorid,
		VrfPeorid:                 vrfPeorid,
		ValidatorPeorid:           validatorPeorid,
		Peorid:                    qrnPeorid + vrfPeorid + validatorPeorid,
	}

	if consensusRound.ConsensusStartBlockHeight <= 0 {
		consensusRound.ConsensusStartBlockHeight = 1
	}

	if consensusRound.Peorid == 0 {
		consensusRound.QrnPeorid = DefaultConsensusRoundQrnPeorid
		consensusRound.VrfPeorid = DefaultConsensusRoundVrfPeorid
		consensusRound.ValidatorPeorid = DefaultConsensusRoundValidatorPeorid
		consensusRound.Peorid = DefaultConsensusRoundPeorid
	}

	return consensusRound
}

func (consensusRoundInfo ConsensusRound) ValidateBasic() error {
	if consensusRoundInfo.ConsensusStartBlockHeight < 0 {
		return fmt.Errorf("wrong ConsensusStartBlockHeight (got: %d)", consensusRoundInfo.ConsensusStartBlockHeight)
	}

	return nil
}

func (consensusRound *ConsensusRound) ToProto() tmproto.ConsensusRound {
	return tmproto.ConsensusRound{
		ConsensusStartBlockHeight: consensusRound.ConsensusStartBlockHeight,
		QrnPeorid:                 consensusRound.QrnPeorid,
		VrfPeorid:                 consensusRound.VrfPeorid,
		ValidatorPeorid:           consensusRound.ValidatorPeorid,
		Peorid:                    consensusRound.Peorid,
	}
}

func ConsensusRoundFromProto(consensusRoundProto tmproto.ConsensusRound) (ConsensusRound, error) {
	return ConsensusRound{
		ConsensusStartBlockHeight: consensusRoundProto.ConsensusStartBlockHeight,
		QrnPeorid:                 consensusRoundProto.QrnPeorid,
		VrfPeorid:                 consensusRoundProto.VrfPeorid,
		ValidatorPeorid:           consensusRoundProto.ValidatorPeorid,
		Peorid:                    consensusRoundProto.Peorid,
	}, nil
}

func UpdateConsensusRound(currentConsensusRound tmproto.ConsensusRound, nextConsensusRound *abci.ConsensusRound) tmproto.ConsensusRound {
	res := currentConsensusRound // explicit copy

	if nextConsensusRound == nil {
		return res
	}

	res.QrnPeorid = nextConsensusRound.QrnPeorid
	res.VrfPeorid = nextConsensusRound.VrfPeorid
	res.ValidatorPeorid = nextConsensusRound.ValidatorPeorid
	res.Peorid = nextConsensusRound.Peorid

	return res
}
