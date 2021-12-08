package types

import (
	"fmt"

	abci "github.com/reapchain/reapchain-core/abci/types"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

const (
	DefaultConsensusRoundQrnPeriod       = 4
	DefaultConsensusRoundVrfPeriod       = 4
	DefaultConsensusRoundValidatorPeriod = 4
	DefaultConsensusRoundPeriod          = 12
)

type ConsensusRound struct {
	ConsensusStartBlockHeight int64  `json:"consensus_start_block_height"`
	Period                    uint64 `json:"period"`
	QrnPeriod                 uint64 `json:"qrn_period"`
	VrfPeriod                 uint64 `json:"vrf_period"`
	ValidatorPeriod           uint64 `json:"validator_period"`
}

func NewConsensusRound(consensusStartBlockHeight int64, qrnPeriod, vrfPeriod, validatorPeriod uint64) ConsensusRound {
	consensusRound := ConsensusRound{
		ConsensusStartBlockHeight: consensusStartBlockHeight,
		QrnPeriod:                 qrnPeriod,
		VrfPeriod:                 vrfPeriod,
		ValidatorPeriod:           validatorPeriod,
		Period:                    qrnPeriod + vrfPeriod + validatorPeriod,
	}

	if consensusRound.ConsensusStartBlockHeight <= 0 {
		consensusRound.ConsensusStartBlockHeight = 1
	}

	if consensusRound.Period == 0 {
		consensusRound.QrnPeriod = DefaultConsensusRoundQrnPeriod
		consensusRound.VrfPeriod = DefaultConsensusRoundVrfPeriod
		consensusRound.ValidatorPeriod = DefaultConsensusRoundValidatorPeriod
		consensusRound.Period = DefaultConsensusRoundPeriod
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
		QrnPeriod:                 consensusRound.QrnPeriod,
		VrfPeriod:                 consensusRound.VrfPeriod,
		ValidatorPeriod:           consensusRound.ValidatorPeriod,
		Period:                    consensusRound.Period,
	}
}

func ConsensusRoundFromProto(consensusRoundProto tmproto.ConsensusRound) (ConsensusRound, error) {
	return ConsensusRound{
		ConsensusStartBlockHeight: consensusRoundProto.ConsensusStartBlockHeight,
		QrnPeriod:                 consensusRoundProto.QrnPeriod,
		VrfPeriod:                 consensusRoundProto.VrfPeriod,
		ValidatorPeriod:           consensusRoundProto.ValidatorPeriod,
		Period:                    consensusRoundProto.Period,
	}, nil
}

func ValidateConsensusRound(consensusRound tmproto.ConsensusRound, height int64) error {
	if consensusRound.ConsensusStartBlockHeight >= height {
		return fmt.Errorf("QrnPeriod must be greater than height. Got %d, height %d",
			consensusRound.QrnPeriod, height)
	}

	if consensusRound.QrnPeriod <= 0 {
		return fmt.Errorf("QrnPeriod must be greater than 0. Got %d",
			consensusRound.QrnPeriod)
	}

	if consensusRound.VrfPeriod <= 0 {
		return fmt.Errorf("VrfPeriod must be greater than 0. Got %d",
			consensusRound.VrfPeriod)
	}

	if consensusRound.ValidatorPeriod <= 0 {
		return fmt.Errorf("ValidatorPeriod must be greater than 0. Got %d",
			consensusRound.ValidatorPeriod)
	}

	checkPeriod := consensusRound.QrnPeriod + consensusRound.VrfPeriod + consensusRound.ValidatorPeriod
	if checkPeriod != consensusRound.Period {
		return fmt.Errorf("Period must be same QrnPeriod + VrfPeriod + ValidatorPeriod. Expected: %d, Got %d",
			consensusRound.Period, checkPeriod)
	}

	return nil
}

func UpdateConsensusRound(currentConsensusRound tmproto.ConsensusRound, nextConsensusRound *abci.ConsensusRound) tmproto.ConsensusRound {
	res := currentConsensusRound // explicit copy

	if nextConsensusRound == nil {
		return res
	}

	res.QrnPeriod = nextConsensusRound.QrnPeriod
	res.VrfPeriod = nextConsensusRound.VrfPeriod
	res.ValidatorPeriod = nextConsensusRound.ValidatorPeriod
	res.Period = nextConsensusRound.Period

	return res
}
