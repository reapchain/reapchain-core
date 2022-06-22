package types

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/crypto/tmhash"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

type SteeringMemberCandidateSet struct {
	mtx     tmsync.Mutex
	
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

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) UpdateWithChangeSet(changes []*SteeringMemberCandidate) error {
	steeringMemberCandidateSet.mtx.Lock()
	defer steeringMemberCandidateSet.mtx.Unlock()
	
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

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) applyUpdates(updates []*SteeringMemberCandidate) {
	existing := steeringMemberCandidateSet.SteeringMemberCandidates
	sort.Sort(SortedSteeringMemberCandidates(existing))

	merged := make([]*SteeringMemberCandidate, len(existing)+len(updates))
	i := 0

	for len(existing) > 0 && len(updates) > 0 {
		if bytes.Compare(existing[0].Address, updates[0].Address) < 0 { // unchanged validator
			merged[i] = existing[0]
			existing = existing[1:]
		} else {
			// Apply add or update.
			merged[i] = updates[0]
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
		merged[i] = existing[j]
		i++
	}
	// OR add updates which are left.
	for j := 0; j < len(updates); j++ {
		merged[i] = updates[j]
		i++
	}

	steeringMemberCandidateSet.SteeringMemberCandidates = merged[:i]
}

func (steeringMemberCandidateSet *SteeringMemberCandidateSet) applyRemovals(deletes []*SteeringMemberCandidate) {
	existing := steeringMemberCandidateSet.SteeringMemberCandidates

	fmt.Println("applyRemovals", "existing", len(existing), "deletes", len(deletes))
	
	merged := make([]*SteeringMemberCandidate, len(existing) - len(deletes))
	i := 0

	// Loop over deletes until we removed all of them.
	for len(deletes) > 0 {
		if bytes.Equal(existing[0].Address, deletes[0].Address) {
			deletes = deletes[1:]
		} else { // Leave it in the resulting slice.
			merged[i] = existing[0]
			i++
		}
		existing = existing[1:]
	}

	// Add the elements which are left.
	for j := 0; j < len(existing); j++ {
		merged[i] = existing[j]
		i++
	}

	steeringMemberCandidateSet.SteeringMemberCandidates = merged[:i]
}