package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto/merkle"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type StandingMemberSet struct {
	Coordinator     *StandingMember   `json:"coordinator"`
	StandingMembers []*StandingMember `json:"standing_members"`
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

	for idx, sm := range standingMemberSet.StandingMembers {
		if err := sm.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid standing member #%d: %w", idx, err)
		}
	}

	// if err := standingMemberSet.Coordinator.ValidateBasic(); err != nil {
	// 	return fmt.Errorf("coordinator failed validate basic, error: %w", err)
	// }

	return nil
}

func (standingMemberSet *StandingMemberSet) HasAddress(address []byte) bool {
	for _, val := range standingMemberSet.StandingMembers {
		if bytes.Equal(val.Address, address) {
			return true
		}
	}
	return false
}

func (standingMemberSet *StandingMemberSet) IsNilOrEmpty() bool {
	return standingMemberSet == nil || len(standingMemberSet.StandingMembers) == 0
}

func (standingMemberSet *StandingMemberSet) Hash() []byte {
	bytesArray := make([][]byte, len(standingMemberSet.StandingMembers))
	for i, sm := range standingMemberSet.StandingMembers {
		bytesArray[i] = sm.Bytes()
	}
	return merkle.HashFromByteSlices(bytesArray)
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

	coordinator, _ := StandingMemberFromProto(standingMemberSetProto.GetCoordinator())

	standingMemberSet.Coordinator = coordinator
	standingMemberSet.StandingMembers = standingMembers

	return standingMemberSet, standingMemberSet.ValidateBasic()
}

func (standingMemberSet *StandingMemberSet) ToProto() (*tmproto.StandingMemberSet, error) {
	if standingMemberSet.IsNilOrEmpty() {
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
	standingMemberCoordinator, _ := standingMemberSet.Coordinator.ToProto()

	standingMemberSetProto.Coordinator = standingMemberCoordinator
	standingMemberSetProto.StandingMembers = standingMembersProto

	return standingMemberSetProto, nil
}

func (standingMemberSet *StandingMemberSet) Copy() *StandingMemberSet {
	return &StandingMemberSet{
		Coordinator:     standingMemberSet.Coordinator,
		StandingMembers: standingMemberCopy(standingMemberSet.StandingMembers),
	}
}

func standingMemberCopy(standingMembers []*StandingMember) []*StandingMember {
	if standingMembers == nil {
		return nil
	}
	standingMembersCopy := make([]*StandingMember, len(standingMembers))
	for i, standingMember := range standingMembers {
		standingMembersCopy[i] = standingMember.Copy()
	}
	return standingMembersCopy
}

func (standingMemberSet *StandingMemberSet) Size() int {
	if standingMemberSet == nil {
		return 0
	}
	return len(standingMemberSet.StandingMembers)
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

func (standingMemberSet *StandingMemberSet) GetByAddress(address []byte) (index int32, standingMember *StandingMember) {
	for idx, standingMember := range standingMemberSet.StandingMembers {
		if bytes.Equal(standingMember.Address, address) {
			return int32(idx), standingMember.Copy()
		}
	}
	return -1, nil
}

//??
func (standingMemberSet *StandingMemberSet) GetCoordinator() (coordinator *StandingMember) {
	if len(standingMemberSet.StandingMembers) == 0 {
		return nil
	}
	if standingMemberSet.Coordinator == nil {
		standingMemberSet.Coordinator = standingMemberSet.findCoordinator()
	}
	return standingMemberSet.Coordinator.Copy()
}

func (standingMemberSet *StandingMemberSet) findCoordinator() *StandingMember {
	var coordinator *StandingMember
	for _, standingMember := range standingMemberSet.StandingMembers {
		if coordinator == nil || !bytes.Equal(standingMember.Address, coordinator.Address) {
			coordinator = coordinator.CompareCoordinatorPriority(standingMember)
		}
	}
	return coordinator
}

// h := sha256.Sum256(bz)
