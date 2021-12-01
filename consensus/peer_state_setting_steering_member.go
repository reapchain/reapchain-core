package consensus

import (
	"fmt"

	"github.com/reapchain/reapchain-core/types"
)

func (ps *PeerState) SetHasSettingSteeringMember(settingSteeringMember *types.SettingSteeringMember) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.NextConsensusStartBlockHeight == settingSteeringMember.Height {
		ps.DidSendSettingSteeringMembers = true
	}
}

func (ps *PeerState) SendSettingSteeringMember(settingSteeringMember *types.SettingSteeringMember) bool {
	if settingSteeringMember == nil {
		return false
	}

	if ps.NextConsensusStartBlockHeight == settingSteeringMember.Height {
		if ps.DidSendSettingSteeringMembers == false {
			msg := &SettingSteeringMemberMessage{settingSteeringMember}

			if ps.peer.Send(SettingSteeringMemberChannel, MustEncode(msg)) {
				ps.SetHasSettingSteeringMember(settingSteeringMember)
				return true
			} else {
				fmt.Println("Faile to send")
			}

		}
	}

	return false
}
