package types

import (
	"bytes"
	encoding_binary "encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type StandingMemberSet struct {
	StandingMembers []*StandingMember `json:"standing_members"`
	Coordinator     *StandingMember   `json:"coordinator"`
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

func (standingMemberSet *StandingMemberSet) SetCoordinator(qrnSet *QrnSet) {
	standingMemberSet.Coordinator = nil

	maxValue := uint64(0)
	qrnSetHash := qrnSet.Hash()

	for _, qrn := range qrnSet.Qrns {
		qrnHash := make([][]byte, 2)
		qrnHash[0] = qrnSetHash
		qrnHash[1] = qrn.GetQrnBytes()

		result := merkle.HashFromByteSlices(qrnHash)
		if uint64(maxValue) < encoding_binary.LittleEndian.Uint64(result) {
			maxValue = encoding_binary.LittleEndian.Uint64(result)
			_, standingMember := standingMemberSet.GetStandingMemberByIdx(qrn.StandingMemberIndex)
			if standingMember != nil {
				standingMemberSet.Coordinator = standingMember
			}
		}
	}

	fmt.Println("cordinator is ", standingMemberSet.Coordinator.Address)
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
		Coordinator:     standingMemberSet.Coordinator,
		StandingMembers: standingMembers,
	}
}

func (standingMemberSet *StandingMemberSet) Size() int {
	if standingMemberSet == nil {
		return 0
	}
	return len(standingMemberSet.StandingMembers)
}

func (standingMemberSet *StandingMemberSet) updateWithChangeSet(standingMembers []*StandingMember) error {
	if len(standingMembers) != 0 {
		sort.Sort(SortedStandingMembers(standingMembers))
		standingMemberSet.StandingMembers = standingMembers[:]
	}

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

// func (standingMemberSet *StandingMemberSet) GetCoordinator() (coordinator *StandingMember) {
// 	if len(standingMemberSet.StandingMembers) == 0 {
// 		return nil
// 	}
// 	if standingMemberSet.Coordinator == nil {
// 		standingMemberSet.Coordinator = standingMemberSet.findCoordinator()
// 	}
// 	return standingMemberSet.Coordinator.Copy()
// }

// func (standingMemberSet *StandingMemberSet) findCoordinator() *StandingMember {
// 	var coordinator *StandingMember
// 	for _, standingMember := range standingMemberSet.StandingMembers {
// 		if coordinator == nil || !bytes.Equal(standingMember.Address, coordinator.Address) {
// 			coordinator = coordinator.CompareCoordinatorPriority(standingMember)
// 		}
// 	}
// 	return coordinator
// }
