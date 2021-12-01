package consensus

import (
	"errors"
	"fmt"

	"github.com/reapchain/reapchain-core/types"
)

type SettingSteeringMemberMessage struct {
	SettingSteeringMember *types.SettingSteeringMember
}

// ValidateBasic performs basic validation.
func (settingSteeringMemberMessage *SettingSteeringMemberMessage) ValidateBasic() error {
	return settingSteeringMemberMessage.SettingSteeringMember.ValidateBasic()
}

// String returns a string representation.
func (settingSteeringMemberMessage *SettingSteeringMemberMessage) String() string {
	return fmt.Sprintf("[SettingSteeringMember %v]", settingSteeringMemberMessage.SettingSteeringMember)
}

type HasSettingSteeringMemberMessage struct {
	Height int64
	Index  int32
}

// ValidateBasic performs basic validation.
func (m *HasSettingSteeringMemberMessage) ValidateBasic() error {
	if m.Height < 0 {
		return errors.New("negative Height")
	}

	if m.Index < 0 {
		return errors.New("negative Index")
	}
	return nil
}

// String returns a string representation.
func (m *HasSettingSteeringMemberMessage) String() string {
	return fmt.Sprintf("[HasSettingSteeringMember VI:%v V:{%v}]", m.Index, m.Height)
}
