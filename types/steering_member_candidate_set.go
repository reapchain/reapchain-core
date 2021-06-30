package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"gitlab.reappay.net/reapchain/reapchain-core/crypto/merkle"
	tmproto "gitlab.reappay.net/reapchain/reapchain-core/proto/reapchain/types"
)

type SteeringMemberCandidateSet struct {
	SteeringMemberCandidates []*SteeringMemberCandidate `json:"steering_member_candidates"`
}

func NewSteeringMemberCandidateSet(smz []*SteeringMemberCandidate) *SteeringMemberCandidateSet {
	sms := &SteeringMemberCandidateSet{}

	err := sms.updateWithChangeSet(smz)
	if err != nil {
		panic(fmt.Sprintf("Cannot create steering member candidate set: %v", err))
	}

	return sms
}

func (sms *SteeringMemberCandidateSet) UpdateWithChangeSet(smz []*SteeringMemberCandidate) error {
	return sms.updateWithChangeSet(smz)
}

func (sms *SteeringMemberCandidateSet) ValidateBasic() error {
	for idx, sm := range sms.SteeringMemberCandidates {
		if err := sm.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid steering member candidate #%d: %w", idx, err)
		}
	}

	return nil
}

func (sms *SteeringMemberCandidateSet) IsNilOrEmpty() bool {
	return sms == nil || len(sms.SteeringMemberCandidates) == 0
}

func (sms *SteeringMemberCandidateSet) Hash() []byte {
	bzs := make([][]byte, len(sms.SteeringMemberCandidates))
	for i, sm := range sms.SteeringMemberCandidates {
		bzs[i] = sm.Bytes()
	}
	return merkle.HashFromByteSlices(bzs)
}

func SteeringMemberCandidateSetFromProto(vp *tmproto.SteeringMemberCandidateSet) (*SteeringMemberCandidateSet, error) {
	if vp == nil {
		return nil, errors.New("nil steering member candidate set")
	}
	sms := new(SteeringMemberCandidateSet)

	smsProto := make([]*SteeringMemberCandidate, len(vp.SteeringMemberCandidates))
	for i := 0; i < len(vp.SteeringMemberCandidates); i++ {
		v, err := SteeringMemberCandidateFromProto(vp.SteeringMemberCandidates[i])
		if err != nil {
			return nil, err
		}
		smsProto[i] = v
	}
	sms.SteeringMemberCandidates = smsProto

	return sms, sms.ValidateBasic()
}

func (sms *SteeringMemberCandidateSet) ToProto() (*tmproto.SteeringMemberCandidateSet, error) {
	if sms.IsNilOrEmpty() {
		return &tmproto.SteeringMemberCandidateSet{}, nil
	}

	vp := new(tmproto.SteeringMemberCandidateSet)
	smsProto := make([]*tmproto.SteeringMemberCandidate, len(sms.SteeringMemberCandidates))
	for i := 0; i < len(sms.SteeringMemberCandidates); i++ {
		valp, err := sms.SteeringMemberCandidates[i].ToProto()
		if err != nil {
			return nil, err
		}
		smsProto[i] = valp
	}
	vp.SteeringMemberCandidates = smsProto

	return vp, nil
}

func (sms *SteeringMemberCandidateSet) Copy() *SteeringMemberCandidateSet {
	return &SteeringMemberCandidateSet{
		SteeringMemberCandidates: steeringMemberCandidateCopy(sms.SteeringMemberCandidates),
	}
}

func steeringMemberCandidateCopy(sms []*SteeringMemberCandidate) []*SteeringMemberCandidate {
	if sms == nil {
		return nil
	}
	smsCopy := make([]*SteeringMemberCandidate, len(sms))
	for i, sm := range sms {
		smsCopy[i] = sm.Copy()
	}
	return smsCopy
}

func (sms *SteeringMemberCandidateSet) Size() int {
	return len(sms.SteeringMemberCandidates)
}

// 정렬하기 위한 구조체
type SteeringMemberCandidatesByAddress []*SteeringMemberCandidate

func (sms SteeringMemberCandidatesByAddress) Len() int { return len(sms) }

func (sms SteeringMemberCandidatesByAddress) Less(i, j int) bool {
	return bytes.Compare(sms[i].Address, sms[j].Address) == -1
}

func (sms SteeringMemberCandidatesByAddress) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}

func (sms *SteeringMemberCandidateSet) updateWithChangeSet(smz []*SteeringMemberCandidate) error {
	if len(smz) != 0 {
		sort.Sort(SteeringMemberCandidatesByAddress(smz))
		sms.SteeringMemberCandidates = smz[:]
	}

	return nil
}
