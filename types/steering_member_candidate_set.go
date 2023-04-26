package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/crypto/tmhash"
	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
)

// It manages steering member candidate list
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

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) UpdateWithChangeSet(standingMemberCandiates []*SteeringMemberCandidate) error {
	return steeringMemberCandidateSet.updateWithChangeSet(standingMemberCandiates)
}

// update the steering member candiate set, when a SDK sends the changed steering member candiate information and initialize.
// It adds or removes the steering member candiate from the managed list
func (steeringMemberCandidateSet *SteeringMemberCandidateSet) updateWithChangeSet(changes []*SteeringMemberCandidate) error {
	if len(changes) == 0 {
		return nil
	}

	sort.Sort(SortedSteeringMemberCandidates(changes))
	
	removals := make([]*SteeringMemberCandidate, 0, len(changes))
	updates := make([]*SteeringMemberCandidate, 0, len(changes))
	
	var prevAddr Address
	for _, steeringMemberCandidate := range changes {
		if bytes.Equal(steeringMemberCandidate.Address, prevAddr) {
			err := fmt.Errorf("duplicate entry %v in %v", steeringMemberCandidate, steeringMemberCandidate)
			return err
		}

		if steeringMemberCandidate.VotingPower != 0 {
			// add
			updates = append(updates, steeringMemberCandidate)
		} else {
			// remove
			removals = append(removals, steeringMemberCandidate)
		}
		prevAddr = steeringMemberCandidate.Address
	}

	steeringMemberCandidateSet.applyUpdates(updates)
	steeringMemberCandidateSet.applyRemovals(removals)

	sort.Sort(SortedSteeringMemberCandidates(steeringMemberCandidateSet.SteeringMemberCandidates))

	return nil
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) ValidateBasic() error {
	if steeringMemberCandidateSet.IsNil() {
		return errors.New("steering member candidate set is nil or empty")
	}

	for idx, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		if err := steeringMemberCandidate.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid steering member candidate #%d: %w", idx, err)
		}
	}

	return nil
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) IsNil() bool {
	return steeringMemberCandidateSet == nil
}

// Make steering member candiate set hash to validate and be included in the block header
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

// Convert the steering member candiate set's proto puffer type to this type to apply the reapchain-core
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

// Convert the type to proto puffer type to send the type to other peer or SDK
func (steeringMemberCandidateSet *SteeringMemberCandidateSet) ToProto() (*tmproto.SteeringMemberCandidateSet, error) {
	if steeringMemberCandidateSet.IsNil() == true {
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

// Check whether the steering member candiate's address is included in the steering member candidate set
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

// Merges the steering member candidate list with the updates list.
// When two elements with same address are seen, the one from updates is selected.
// Expects updates to be a list of updates sorted by address with no duplicates or errors.
func (steeringMemberCandidateSet *SteeringMemberCandidateSet) applyUpdates(updates []*SteeringMemberCandidate) {
	existing := steeringMemberCandidateSet.SteeringMemberCandidates
	sort.Sort(SortedSteeringMemberCandidates(existing))

	merged := make([]*SteeringMemberCandidate, 0, len(existing)+len(updates))
	i := 0

	for len(existing) > 0 && len(updates) > 0 {
		if bytes.Compare(existing[0].Address, updates[0].Address) < 0 { // unchanged steering member candidate
			merged = append(merged, existing[0])
			existing = existing[1:]
		} else {
			// Apply add or update.
			merged = append(merged, updates[0])
			if bytes.Equal(existing[0].Address, updates[0].Address) {
				// SteeringMemberCandidate is present in both, advance existing.
				existing = existing[1:]
			}
			updates = updates[1:]
		}
		i++
	}

	// Add the elements which are left.
	for j := 0; j < len(existing); j++ {
		merged = append(merged, existing[j])
		i++
	}
	// OR add updates which are left.
	for j := 0; j < len(updates); j++ {
		merged = append(merged, updates[j])
		i++
	}

	steeringMemberCandidateSet.SteeringMemberCandidates = merged[:i]
}

// Removes the steering member candiate specified in 'deletes'.
// Should not fail as verification has been done before.
// Expects vals to be sorted by address (done by applyUpdates).
func (steeringMemberCandidateSet *SteeringMemberCandidateSet) applyRemovals(deletes []*SteeringMemberCandidate) {
	existing := steeringMemberCandidateSet.SteeringMemberCandidates	
	merged := make([]*SteeringMemberCandidate, 0, len(existing))
	i := 0

	for len(existing) > 0 {
		j := 0
		deleteLen := len(deletes)
		for deleteLen > j {
			if bytes.Equal(existing[0].Address, deletes[j].Address) {
				deletes = deletes[j+1:]
				break
			} 
			j++
		} 

		if (deleteLen == j) {
			merged = append(merged, existing[0])
			i++
		}
		existing = existing[1:]
	}

	steeringMemberCandidateSet.SteeringMemberCandidates = merged[:i]
}

// It replace the steering member candiate with new steering member candate list,
// and return the new steering member candidate set.
func SteeringMemberCandidateSetFromExistingSteeringMemberCandidates(steeringMemberCandidates []*SteeringMemberCandidate) (*SteeringMemberCandidateSet, error) {
	if len(steeringMemberCandidates) == 0 {
		return nil, errors.New("steering member candidate set is empty")
	}
	for _, val := range steeringMemberCandidates {
		err := val.ValidateBasic()
		if err != nil {
			return nil, fmt.Errorf("can't create steering member candidate set: %w", err)
		}
	}

	steeringMemberCandidateSet := &SteeringMemberCandidateSet{
		SteeringMemberCandidates: steeringMemberCandidates,
	}
	sort.Sort(SortedSteeringMemberCandidates(steeringMemberCandidateSet.SteeringMemberCandidates))
	return steeringMemberCandidateSet, nil
}