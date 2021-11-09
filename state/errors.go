package state

import "fmt"

type (
	ErrInvalidBlock error
	ErrProxyAppConn error

	ErrUnknownBlock struct {
		Height int64
	}

	ErrBlockHashMismatch struct {
		CoreHash []byte
		AppHash  []byte
		Height   int64
	}

	ErrAppBlockHeightTooHigh struct {
		CoreHeight int64
		AppHeight  int64
	}

	ErrAppBlockHeightTooLow struct {
		AppHeight int64
		StoreBase int64
	}

	ErrLastStateMismatch struct {
		Height int64
		Core   []byte
		App    []byte
	}

	ErrStateMismatch struct {
		Got      *State
		Expected *State
	}

	ErrNoValSetForHeight struct {
		Height int64
	}

	ErrNoStandingMemberSetForHeight struct {
		Height int64
	}

	ErrNoSteeringMemberSetForHeight struct {
		Height int64
	}

	ErrNoQrnSetForHeight struct {
		Height int64
	}

	ErrNoSettingSteeringMemberForHeight struct {
		Height int64
	}

	ErrNoVrfSetForHeight struct {
		Height int64
	}

	ErrNoConsensusParamsForHeight struct {
		Height int64
	}

	ErrNoABCIResponsesForHeight struct {
		Height int64
	}

	ErrNoConsensusRoundForHeight struct {
		Height int64
	}
)

func (e ErrUnknownBlock) Error() string {
	return fmt.Sprintf("could not find block #%d", e.Height)
}

func (e ErrBlockHashMismatch) Error() string {
	return fmt.Sprintf(
		"app block hash (%X) does not match core block hash (%X) for height %d",
		e.AppHash,
		e.CoreHash,
		e.Height,
	)
}

func (e ErrAppBlockHeightTooHigh) Error() string {
	return fmt.Sprintf("app block height (%d) is higher than core (%d)", e.AppHeight, e.CoreHeight)
}

func (e ErrAppBlockHeightTooLow) Error() string {
	return fmt.Sprintf("app block height (%d) is too far below block store base (%d)", e.AppHeight, e.StoreBase)
}

func (e ErrLastStateMismatch) Error() string {
	return fmt.Sprintf(
		"latest reapchain block (%d) LastAppHash (%X) does not match app's AppHash (%X)",
		e.Height,
		e.Core,
		e.App,
	)
}

func (e ErrStateMismatch) Error() string {
	return fmt.Sprintf(
		"state after replay does not match saved state. Got ----\n%v\nExpected ----\n%v\n",
		e.Got,
		e.Expected,
	)
}

func (e ErrNoValSetForHeight) Error() string {
	return fmt.Sprintf("could not find validator set for height #%d", e.Height)
}

func (e ErrNoStandingMemberSetForHeight) Error() string {
	return fmt.Sprintf("could not find standing member set for height #%d", e.Height)
}

func (e ErrNoSteeringMemberSetForHeight) Error() string {
	return fmt.Sprintf("could not find steering member candidate set for height #%d", e.Height)
}

func (e ErrNoQrnSetForHeight) Error() string {
	return fmt.Sprintf("could not find qrn set for height #%d", e.Height)
}

func (e ErrNoSettingSteeringMemberForHeight) Error() string {
	return fmt.Sprintf("could not find qrn set for height #%d", e.Height)
}

func (e ErrNoVrfSetForHeight) Error() string {
	return fmt.Sprintf("could not find vrf set for height #%d", e.Height)
}

func (e ErrNoConsensusParamsForHeight) Error() string {
	return fmt.Sprintf("could not find consensus params for height #%d", e.Height)
}

func (e ErrNoABCIResponsesForHeight) Error() string {
	return fmt.Sprintf("could not find results for height #%d", e.Height)
}

func (e ErrNoConsensusRoundForHeight) Error() string {
	return fmt.Sprintf("could not find consensus round for height #%d", e.Height)
}
