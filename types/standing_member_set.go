package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"gitlab.reappay.net/reapchain/reapchain-core/crypto/merkle"
	tmproto "gitlab.reappay.net/reapchain/reapchain-core/proto/reapchain/types"
)

type StandingMemberSet struct {
	StandingMembers []*StandingMember `json:"standing_members"`
}

func NewStandingMemberSet(smz []*StandingMember) *StandingMemberSet {
	sms := &StandingMemberSet{}

	err := sms.updateWithChangeSet(smz)
	if err != nil {
		panic(fmt.Sprintf("Cannot create standing member set: %v", err))
	}

	return sms
}

func (sms *StandingMemberSet) UpdateWithChangeSet(smz []*StandingMember) error {
	return sms.updateWithChangeSet(smz)
}

func (sms *StandingMemberSet) ValidateBasic() error {
	if sms.IsNilOrEmpty() {
		return errors.New("standing member set is nil or empty")
	}

	for idx, sm := range sms.StandingMembers {
		if err := sm.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid standing member #%d: %w", idx, err)
		}
	}

	return nil
}

func (sms *StandingMemberSet) IsNilOrEmpty() bool {
	return sms == nil || len(sms.StandingMembers) == 0
}

func (sms *StandingMemberSet) Hash() []byte {
	bzs := make([][]byte, len(sms.StandingMembers))
	for i, sm := range sms.StandingMembers {
		bzs[i] = sm.Bytes()
	}
	return merkle.HashFromByteSlices(bzs)
}

func StandingMemberSetFromProto(vp *tmproto.StandingMemberSet) (*StandingMemberSet, error) {
	if vp == nil {
		return nil, errors.New("nil standing member set")
	}
	sms := new(StandingMemberSet)

	smsProto := make([]*StandingMember, len(vp.StandingMembers))
	for i := 0; i < len(vp.StandingMembers); i++ {
		v, err := StandingMemberFromProto(vp.StandingMembers[i])
		if err != nil {
			return nil, err
		}
		smsProto[i] = v
	}
	sms.StandingMembers = smsProto

	return sms, sms.ValidateBasic()
}

func (sms *StandingMemberSet) ToProto() (*tmproto.StandingMemberSet, error) {
	if sms.IsNilOrEmpty() {
		return &tmproto.StandingMemberSet{}, nil
	}

	vp := new(tmproto.StandingMemberSet)
	smsProto := make([]*tmproto.StandingMember, len(sms.StandingMembers))
	for i := 0; i < len(sms.StandingMembers); i++ {
		valp, err := sms.StandingMembers[i].ToProto()
		if err != nil {
			return nil, err
		}
		smsProto[i] = valp
	}
	vp.StandingMembers = smsProto

	return vp, nil
}

func (sms *StandingMemberSet) Copy() *StandingMemberSet {
	return &StandingMemberSet{
		StandingMembers: standingMemberCopy(sms.StandingMembers),
	}
}

func standingMemberCopy(sms []*StandingMember) []*StandingMember {
	if sms == nil {
		return nil
	}
	smsCopy := make([]*StandingMember, len(sms))
	for i, sm := range sms {
		smsCopy[i] = sm.Copy()
	}
	return smsCopy
}

func (sms *StandingMemberSet) Size() int {
	return len(sms.StandingMembers)
}

// 정렬하기 위한 구조체
type StandingMembersByAddress []*StandingMember

func (sms StandingMembersByAddress) Len() int { return len(sms) }

func (sms StandingMembersByAddress) Less(i, j int) bool {
	return bytes.Compare(sms[i].Address, sms[j].Address) == -1
}

func (sms StandingMembersByAddress) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}

func (sms *StandingMemberSet) updateWithChangeSet(smz []*StandingMember) error {
	if len(smz) != 0 {
		sort.Sort(StandingMembersByAddress(smz))
		sms.StandingMembers = smz[:]
	}

	return nil
}

func (standingMemberSet *StandingMemberSet) GetByIndex(index int32) (address []byte, standingMember *StandingMember) {
	if index < 0 || int(index) >= len(standingMemberSet.StandingMembers) {
		return nil, nil
	}
	standingMember = standingMemberSet.StandingMembers[index]
	return standingMember.Address, standingMember.Copy()
}
