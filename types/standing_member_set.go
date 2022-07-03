package types

import (
	"bytes"
	encoding_binary "encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

type StandingMemberSet struct {
	StandingMembers           []*StandingMember `json:"standing_members"`
	Coordinator               *StandingMember   `json:"coordinator"`
	CurrentCoordinatorRanking int64             `json:"current_coordinator_ranking"`
}

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

func (standingMemberSet *StandingMemberSet) HasAddress(address []byte) bool {
	for _, standingMember := range standingMemberSet.StandingMembers {
		if bytes.Equal(standingMember.Address, address) {
			return true
		}
	}
	return false
}

func (standingMemberSet *StandingMemberSet) Hash() []byte {
	bytesArray := make([][]byte, len(standingMemberSet.StandingMembers))
	for idx, standingMember := range standingMemberSet.StandingMembers {
		bytesArray[idx] = standingMember.Bytes()
	}
	return merkle.HashFromByteSlices(bytesArray)
}

type QrnHash struct {
	Address   Address `json:"address"`
	HashValue uint64  `json:"hash_value"`
}

type QrnHashsByValue []*QrnHash

func (qrnHash QrnHashsByValue) Len() int { return len(qrnHash) }

func (qrnHash QrnHashsByValue) Less(i, j int) bool {
	return qrnHash[i].HashValue > qrnHash[j].HashValue
}

func (qrnHash QrnHashsByValue) Swap(i, j int) {
	qrnHash[i], qrnHash[j] = qrnHash[j], qrnHash[i]
}

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

	for standingMember == nil {
		standingMemberSet.CurrentCoordinatorRanking++
		if standingMemberSet.CurrentCoordinatorRanking == int64(qrnSet.Size()) {
			standingMemberSet.CurrentCoordinatorRanking = 0
		}
		_, standingMember = standingMemberSet.GetStandingMemberByAddress(qrnHashs[standingMemberSet.CurrentCoordinatorRanking].Address)
	}
	standingMemberSet.Coordinator = standingMember
}

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

func (standingMemberSet *StandingMemberSet) applyUpdates(updates []*StandingMember) {
	existing := standingMemberSet.StandingMembers
	sort.Sort(SortedStandingMembers(existing))

	merged := make([]*StandingMember, len(existing)+len(updates))
	i := 0

	for len(existing) > 0 && len(updates) > 0 {
		if bytes.Compare(existing[0].Address, updates[0].Address) < 0 { // unchanged validator
			merged[i] = existing[0]
			existing = existing[1:]
		} else {
			// Apply add or update.
			merged[i] = updates[0]
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
		merged[i] = existing[j]
		i++
	}
	// OR add updates which are left.
	for j := 0; j < len(updates); j++ {
		merged[i] = updates[j]
		i++
	}

	standingMemberSet.StandingMembers = merged[:i]
}

func (standingMemberSet *StandingMemberSet) applyRemovals(deletes []*StandingMember) {
	existing := standingMemberSet.StandingMembers

	merged := make([]*StandingMember, len(existing)-len(deletes))
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

	standingMemberSet.StandingMembers = merged[:i]
}