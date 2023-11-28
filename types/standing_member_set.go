package types

import (
	"bytes"
	encoding_binary "encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/crypto/tmhash"
	tmproto "github.com/reapchain/reapchain-core/proto/podc/types"
)

// It manages standing member list and coordinator
type StandingMemberSet struct {
	StandingMembers           []*StandingMember `json:"standing_members"`
	Coordinator               *StandingMember   `json:"coordinator"`
	CurrentCoordinatorRanking int64             `json:"current_coordinator_ranking"` // current coordinator index
}

type QrnHash struct {
	Address   Address `json:"address"`
	HashValue uint64  `json:"hash_value"`
}

type QrnHashsByValue []*QrnHash

func NewStandingMemberSet(standingMembers []*StandingMember) *StandingMemberSet {
	standingMemberSet := &StandingMemberSet{}

	err := standingMemberSet.updateWithChangeSet(standingMembers)
	if err != nil {
		panic(fmt.Sprintf("Cannot create standing member set: %v", err))
	}
	return standingMemberSet
}

func (standingMemberSet *StandingMemberSet) UpdateWithChangeSet(standingMembers []*StandingMember) error {
	return standingMemberSet.updateWithChangeSet(standingMembers)
}

// update the standing member set, when a SDK sends the changed standing member information and initialize.
// It adds or removes the standing member from the managed list
func (standingMemberSet *StandingMemberSet) updateWithChangeSet(changes []*StandingMember) error {
	if len(changes) == 0 {
		return nil
	}

	sort.Sort(SortedStandingMembers(changes))
	
	removals := make([]*StandingMember, 0, len(changes))
	updates := make([]*StandingMember, 0, len(changes))
	
	var prevAddr Address
	for _, standingMember := range changes {
		if bytes.Equal(standingMember.Address, prevAddr) {
			err := fmt.Errorf("duplicate entry %v in %v", standingMember, standingMember)
			return err
		}

		if standingMember.VotingPower != 0 {
			// add
			updates = append(updates, standingMember)
		} else {
			// remove
			removals = append(removals, standingMember)
		}
		prevAddr = standingMember.Address
	}

	standingMemberSet.applyUpdates(updates)
	standingMemberSet.applyRemovals(removals)

	sort.Sort(SortedStandingMembers(standingMemberSet.StandingMembers))

	return nil
}

func (standingMemberSet *StandingMemberSet) ValidateBasic() error {
	if standingMemberSet.IsNilOrEmpty() {
		return errors.New("standing member set is nil or empty")
	}

	for idx, standingMember := range standingMemberSet.StandingMembers {
		if err := standingMember.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid standing member #%d: %w", idx, err)
		}
	}

	return nil
}

func (standingMemberSet *StandingMemberSet) IsNilOrEmpty() bool {
	return standingMemberSet == nil || len(standingMemberSet.StandingMembers) == 0
}

// Make standing member set hash to validate and be included in the block header 
func (standingMemberSet *StandingMemberSet) Hash() []byte {
	if standingMemberSet == nil {
		return tmhash.Sum([]byte{})
	}
	
	bytesArray := make([][]byte, len(standingMemberSet.StandingMembers))
	for idx, standingMember := range standingMemberSet.StandingMembers {
		bytesArray[idx] = standingMember.Bytes()
	}
	return merkle.HashFromByteSlices(bytesArray)
}

// Convert the standing member set's proto puffer type to this type to apply the reapchain-core
func StandingMemberSetFromProto(standingMemberSetProto *tmproto.StandingMemberSet) (*StandingMemberSet, error) {
	if standingMemberSetProto == nil {
		return nil, errors.New("nil standing member set")
	}

	standingMemberSet := new(StandingMemberSet)

	standingMembers := make([]*StandingMember, len(standingMemberSetProto.StandingMembers))
	for i := 0; i < len(standingMemberSetProto.StandingMembers); i++ {
		standingMember, err := StandingMemberFromProto(standingMemberSetProto.StandingMembers[i])
		if err != nil {
			return nil, err
		}
		standingMembers[i] = standingMember
	}

	standingMemberSet.StandingMembers = standingMembers
	standingMemberSet.CurrentCoordinatorRanking = standingMemberSetProto.CurrentCoordinatorRanking

	return standingMemberSet, standingMemberSet.ValidateBasic()
}

// Convert the type to proto puffer type to send the type to other peer or SDK
func (standingMemberSet *StandingMemberSet) ToProto() (*tmproto.StandingMemberSet, error) {
	if standingMemberSet.IsNilOrEmpty() == true {
		return &tmproto.StandingMemberSet{}, nil
	}

	standingMemberSetProto := new(tmproto.StandingMemberSet)
	standingMembersProto := make([]*tmproto.StandingMember, len(standingMemberSet.StandingMembers))

	for i := 0; i < len(standingMemberSet.StandingMembers); i++ {
		standingMemberProto, err := standingMemberSet.StandingMembers[i].ToProto()
		if err != nil {
			return nil, err
		}
		standingMembersProto[i] = standingMemberProto
	}

	standingMemberSetProto.StandingMembers = standingMembersProto
	standingMemberSetProto.CurrentCoordinatorRanking = standingMemberSet.CurrentCoordinatorRanking

	return standingMemberSetProto, nil
}

func (standingMemberSet *StandingMemberSet) Copy() *StandingMemberSet {
	var standingMembers []*StandingMember

	if standingMemberSet.StandingMembers == nil {
		standingMembers = nil
	} else {
		standingMembers = make([]*StandingMember, len(standingMemberSet.StandingMembers))
		for i, standingMember := range standingMemberSet.StandingMembers {
			standingMembers[i] = standingMember.Copy()
		}
	}

	return &StandingMemberSet{
		Coordinator:               standingMemberSet.Coordinator,
		StandingMembers:           standingMembers,
		CurrentCoordinatorRanking: standingMemberSet.CurrentCoordinatorRanking,
	}
}

func (standingMemberSet *StandingMemberSet) Size() int {
	if standingMemberSet == nil {
		return 0
	}
	return len(standingMemberSet.StandingMembers)
}
// Check whether the standing member's address is included in the standing member set
func (standingMemberSet *StandingMemberSet) HasAddress(address []byte) bool {
	for _, standingMember := range standingMemberSet.StandingMembers {
		if bytes.Equal(standingMember.Address, address) {
			return true
		}
	}
	return false
}

func (standingMemberSet *StandingMemberSet) GetStandingMemberByIdx(idx int32) (address []byte, standingMember *StandingMember) {
	if idx < 0 || int(idx) >= len(standingMemberSet.StandingMembers) {
		return nil, nil
	}
	standingMember = standingMemberSet.StandingMembers[idx]
	return standingMember.Address, standingMember.Copy()
}

func (standingMemberSet *StandingMemberSet) GetStandingMemberByAddress(address []byte) (idx int32, standingMember *StandingMember) {
	if standingMemberSet == nil {
		return -1, nil
	}
	for idx, standingMember := range standingMemberSet.StandingMembers {
		if bytes.Equal(standingMember.Address, address) {
			return int32(idx), standingMember.Copy()
		}
	}
	return -1, nil
}

// Merges the standing member list with the updates list.
// When two elements with same address are seen, the one from updates is selected.
// Expects updates to be a list of updates sorted by address with no duplicates or errors.
func (standingMemberSet *StandingMemberSet) applyUpdates(updates []*StandingMember) {
	existing := standingMemberSet.StandingMembers
	sort.Sort(SortedStandingMembers(existing))

	merged := make([]*StandingMember, 0, len(existing)+len(updates))
	i := 0

	for len(existing) > 0 && len(updates) > 0 {
		if bytes.Compare(existing[0].Address, updates[0].Address) < 0 { // unchanged standing member
			merged = append(merged, existing[0])
			existing = existing[1:]
		} else {
			// Apply add or update.
			merged = append(merged, updates[0])
			if bytes.Equal(existing[0].Address, updates[0].Address) {
				// StandingMember is present in both, advance existing.
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

	standingMemberSet.StandingMembers = merged[:i]
}

// Removes the standing member specified in 'deletes'.
// Should not fail as verification has been done before.
// Expects vals to be sorted by address (done by applyUpdates).
func (standingMemberSet *StandingMemberSet) applyRemovals(deletes []*StandingMember) {
	existing := standingMemberSet.StandingMembers
	merged := make([]*StandingMember, 0, len(existing))
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

	standingMemberSet.StandingMembers = merged[:i]
}

func StandingMemberSetFromExistingStandingMembers(standingMembers []*StandingMember) (*StandingMemberSet, error) {
	sort.Sort(SortedStandingMembers(standingMembers))

	if len(standingMembers) == 0 {
		return nil, errors.New("standing member set is empty")
	}
	for _, standingMember := range standingMembers {
		err := standingMember.ValidateBasic()
		if err != nil {
			return nil, fmt.Errorf("can't create standing member set: %w", err)
		}
	}

	standingMemberSet := &StandingMemberSet{
		StandingMembers: standingMembers[:],
	}
	
	return standingMemberSet, nil
}

// ---- Sorting with qrn hash ----
func (qrnHash QrnHashsByValue) Len() int { return len(qrnHash) }

func (qrnHash QrnHashsByValue) Less(i, j int) bool {
	return qrnHash[i].HashValue > qrnHash[j].HashValue
}

func (qrnHash QrnHashsByValue) Swap(i, j int) {
	qrnHash[i], qrnHash[j] = qrnHash[j], qrnHash[i]
}
// --------------------------------


// Set the coordinator with qrnSet
func (standingMemberSet *StandingMemberSet) SetCoordinator(qrnSet *QrnSet) {
	standingMemberSet.Coordinator = nil

	qrnHashs := make([]*QrnHash, qrnSet.Size())
	qrnSetHash := qrnSet.Hash()

	for i, qrn := range qrnSet.Qrns {
		qrnHash := make([][]byte, 2)
		qrnHash[0] = qrnSetHash
		qrnHash[1] = qrn.GetQrnBytes()

		address := qrn.StandingMemberPubKey.Address()
		index, _ := standingMemberSet.GetStandingMemberByAddress(address)

		qrnHashs[i] = &QrnHash{
			Address:   qrn.StandingMemberPubKey.Address(),
			HashValue: 0,
		}

		if index != -1 && qrn.Signature != nil {
			qrnHashs[i].HashValue = encoding_binary.LittleEndian.Uint64(merkle.HashFromByteSlices(qrnHash))
		}
	}
	
	sort.Sort(QrnHashsByValue(qrnHashs))
	_, standingMember := standingMemberSet.GetStandingMemberByAddress(qrnHashs[standingMemberSet.CurrentCoordinatorRanking].Address)

	fmt.Println("==================================================================")
	for i, hash := range qrnHashs {
		fmt.Println(i, "HashValue", hash.HashValue, "address", hash.Address)
	}
	fmt.Println("standingMemberSet.CurrentCoordinatorRanking", standingMemberSet.CurrentCoordinatorRanking)
	fmt.Println("standingMember", standingMember)
	fmt.Println("==================================================================")

	for standingMember == nil {
		standingMemberSet.CurrentCoordinatorRanking++
		if standingMemberSet.CurrentCoordinatorRanking == int64(qrnSet.Size()) {
			standingMemberSet.CurrentCoordinatorRanking = 0
		}
		_, standingMember = standingMemberSet.GetStandingMemberByAddress(qrnHashs[standingMemberSet.CurrentCoordinatorRanking].Address)
	}
	standingMemberSet.Coordinator = standingMember
}