package types

import (
	"errors"
	"fmt"

	"github.com/reapchain/reapchain/crypto/merkle"
	tmproto "github.com/reapchain/reapchain/proto/reapchain/types"
)

type StandingMemberSet struct {
	StandingMembers []*StandingMember `json:"standing_members"`
}

func NewStandingMemberSet(smz []*StandingMember) *StandingMemberSet {
	sms := &StandingMemberSet{}

	if len(smz) == 0 {
		sms.StandingMembers = make([]*StandingMember, 0)
	} else {
		sms.StandingMembers = smz
	}

	return sms
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
		StandingMembers: standingMemberListCopy(sms.StandingMembers),
	}
}

func standingMemberListCopy(sms []*StandingMember) []*StandingMember {
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
