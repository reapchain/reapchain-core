package consensus

import (
	"github.com/reapchain/reapchain-core/types"
)

func (ps *PeerState) SetHasSettingSteeringMember(height int64) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.NextConsensusStartBlockHeight == height {
		ps.DidSendSettingSteeringMember = true
	}
}

func (ps *PeerState) PickSendSettingSteeringMember(settingSteeringMember *types.SettingSteeringMember) bool {
	if settingSteeringMember == nil {
		return false
	}

	if ps.NextConsensusStartBlockHeight == settingSteeringMember.Height {
		if ps.DidSendSettingSteeringMember == false {
			msg := &SettingSteeringMemberMessage{settingSteeringMember.Copy()}
			if ps.peer.Send(SettingSteeringMemberChannel, MustEncode(msg)) {
				ps.SetHasSettingSteeringMember(settingSteeringMember.Height)
				return true
			} else {
				ps.logger.Debug("SendSettingSteeringMember: Faile to send")
			}
		}
	}

	return false
}

func (ps *PeerState) RequestSettingSteeringMember(height int64) bool {
	
	msg := &RequestSettingSteeringMemberMessage{Height: height}
	
	if !ps.peer.Send(SettingSteeringMemberChannel, MustEncode(msg)) {
		ps.logger.Debug("RequestSettingSteeringMember: Faile to send")
	}

	return false
}

func (ps *PeerState) ReSendSettingSteeringMember(settingSteeringMember *types.SettingSteeringMember) bool {
	if settingSteeringMember == nil {
		return false
	}
	msg := &ResponseSettingSteeringMemberMessage{settingSteeringMember.Copy()}
	
	if ps.peer.Send(SettingSteeringMemberChannel, MustEncode(msg)) {
		ps.SetHasSettingSteeringMember(settingSteeringMember.Height)
		return true
	} else {
		ps.logger.Debug("SendSettingSteeringMember: Faile to send")
	}

	return false
}