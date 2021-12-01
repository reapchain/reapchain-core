package consensus

import (
	"errors"
	"fmt"

	"github.com/reapchain/reapchain-core/types"
)

type VrfMessage struct {
	Vrf *types.Vrf
}

// ValidateBasic performs basic validation.
func (vrfMessage *VrfMessage) ValidateBasic() error {
	return vrfMessage.Vrf.ValidateBasic()
}

// String returns a string representation.
func (vrfMessage *VrfMessage) String() string {
	return fmt.Sprintf("[Vrf %v]", vrfMessage.Vrf)
}

type HasVrfMessage struct {
	Height int64
	Index  int32
}

// ValidateBasic performs basic validation.
func (m *HasVrfMessage) ValidateBasic() error {
	if m.Height < 0 {
		return errors.New("negative Height")
	}

	if m.Index < 0 {
		return errors.New("negative Index")
	}
	return nil
}

// String returns a string representation.
func (m *HasVrfMessage) String() string {
	return fmt.Sprintf("[HasVrf VI:%v V:{%v}]", m.Index, m.Height)
}
