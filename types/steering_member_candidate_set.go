package types

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/crypto/tmhash"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type SteeringMemberCandidateSet struct {
	SteeringMemberCandidates []*SteeringMemberCandidate `json:"steering_member_candidates"`
}

func NewSteeringMemberCandidateSet(steeringMemberCandidates []*SteeringMemberCandidate) *SteeringMemberCandidateSet {
	steeringMemberCandidateSet := &SteeringMemberCandidateSet{}

	err := steeringMemberCandidateSet.UpdateWithChangeSet(steeringMemberCandidates)
	if err != nil {
		panic(fmt.Sprintf("Cannot create steering member candidate set: %v", err))
	}

	return steeringMemberCandidateSet
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) ValidateBasic() error {
	// if steeringMemberCandidateSet.IsNilOrEmpty() {
	// 	return errors.New("steering member candidate set is nil or empty")
	// 	// return nil
	// }

	for idx, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		if err := steeringMemberCandidate.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid steering member candidate #%d: %w", idx, err)
		}
	}

	return nil
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) IsNilOrEmpty() bool {
	return steeringMemberCandidateSet == nil || len(steeringMemberCandidateSet.SteeringMemberCandidates) == 0
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) Hash() []byte {
	if steeringMemberCandidateSet == nil {
		return tmhash.Sum([]byte{})
	}

	bytesArray := make([][]byte, len(steeringMemberCandidateSet.SteeringMemberCandidates))
	for idx, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		bytesArray[idx] = steeringMemberCandidate.Bytes()
	}
	return merkle.HashFromByteSlices(bytesArray)
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) ToProto() (*tmproto.SteeringMemberCandidateSet, error) {
	if steeringMemberCandidateSet.IsNilOrEmpty() == true {
		return &tmproto.SteeringMemberCandidateSet{}, nil
	}

	steeringMemberCandidateSetProto := new(tmproto.SteeringMemberCandidateSet)
	steeringMemberCandidatesProto := make([]*tmproto.SteeringMemberCandidate, len(steeringMemberCandidateSet.SteeringMemberCandidates))

	for i := 0; i < len(steeringMemberCandidateSet.SteeringMemberCandidates); i++ {
		steeringMemberCandidateProto, err := steeringMemberCandidateSet.SteeringMemberCandidates[i].ToProto()
		if err != nil {
			return nil, err
		}
		steeringMemberCandidatesProto[i] = steeringMemberCandidateProto
	}

	steeringMemberCandidateSetProto.SteeringMemberCandidates = steeringMemberCandidatesProto

	return steeringMemberCandidateSetProto, nil
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) Copy() *SteeringMemberCandidateSet {
	if steeringMemberCandidateSet == nil {
		return nil
	}

	var steeringMemberCandidates []*SteeringMemberCandidate

	if steeringMemberCandidateSet.SteeringMemberCandidates == nil {
		steeringMemberCandidates = nil
	} else {
		steeringMemberCandidates = make([]*SteeringMemberCandidate, len(steeringMemberCandidateSet.SteeringMemberCandidates))
		for i, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
			steeringMemberCandidates[i] = steeringMemberCandidate.Copy()
		}
	}

	return &SteeringMemberCandidateSet{
		SteeringMemberCandidates: steeringMemberCandidates,
	}
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) Size() int {
	if steeringMemberCandidateSet == nil {
		return 0
	}
	return len(steeringMemberCandidateSet.SteeringMemberCandidates)
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) UpdateWithChangeSet(steeringMemberCandidates []*SteeringMemberCandidate) error {
	if len(steeringMemberCandidates) != 0 {
		sort.Sort(SortedSteeringMemberCandidates(steeringMemberCandidates))
		steeringMemberCandidateSet.SteeringMemberCandidates = steeringMemberCandidates[:]
	}

	return nil
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) HasAddress(address []byte) bool {
	for _, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		if bytes.Equal(steeringMemberCandidate.Address, address) {
			return true
		}
	}
	return false
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) GetSteeringMemberCandidateByIdx(idx int32) (address []byte, steeringMemberCandidate *SteeringMemberCandidate) {
	if idx < 0 || int(idx) >= len(steeringMemberCandidateSet.SteeringMemberCandidates) {
		return nil, nil
	}
	steeringMemberCandidate = steeringMemberCandidateSet.SteeringMemberCandidates[idx]
	return steeringMemberCandidate.Address, steeringMemberCandidate.Copy()
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) GetSteeringMemberCandidateByAddress(address []byte) (idx int32, steeringMemberCandidate *SteeringMemberCandidate) {
	if steeringMemberCandidateSet == nil {
		return -1, nil
	}
	for idx, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		if bytes.Equal(steeringMemberCandidate.Address, address) {
			return int32(idx), steeringMemberCandidate.Copy()
		}
	}
	return -1, nil
}

func SteeringMemberCandidateSetFromProto(steeringMemberCandidateSetProto *tmproto.SteeringMemberCandidateSet) (*SteeringMemberCandidateSet, error) {
	steeringMemberCandidateSet := new(SteeringMemberCandidateSet)
	if steeringMemberCandidateSetProto == nil {
		return steeringMemberCandidateSet, nil
	}

	steeringMemberCandidates := make([]*SteeringMemberCandidate, len(steeringMemberCandidateSetProto.SteeringMemberCandidates))
	for i := 0; i < len(steeringMemberCandidateSetProto.SteeringMemberCandidates); i++ {
		steeringMemberCandidate, err := SteeringMemberCandidateFromProto(steeringMemberCandidateSetProto.SteeringMemberCandidates[i])
		if err != nil {
			return nil, err
		}
		steeringMemberCandidates[i] = steeringMemberCandidate
	}

	steeringMemberCandidateSet.SteeringMemberCandidates = steeringMemberCandidates

	return steeringMemberCandidateSet, steeringMemberCandidateSet.ValidateBasic()
}
