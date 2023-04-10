package types

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/reapchain/reapchain-core/crypto"
	"github.com/reapchain/reapchain-core/crypto/merkle"
	"github.com/reapchain/reapchain-core/libs/bits"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

const MAXIMUM_STEERING_MEMBER_CANDIDATES = 30

// It manages the status of VRF-related lists and blocks heights for consensus.
type VrfSet struct {
	Height int64

	mtx          tmsync.Mutex
	Vrfs         []*Vrf
	VrfsBitArray *bits.BitArray // It is used to check existings of VRF 
}

func NewVrfSet(height int64, steeringMemberCandidateSet *SteeringMemberCandidateSet, vrfs []*Vrf) *VrfSet {

	if vrfs == nil || len(vrfs) == 0 {
		vrfs = make([]*Vrf, steeringMemberCandidateSet.Size())
		for i, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
			vrfs[i] = NewVrfAsEmpty(height, steeringMemberCandidate.PubKey)
			vrfs[i].VrfIndex = int32(i)
		}
	}

	return &VrfSet{
		Height:                     height,
		VrfsBitArray:               bits.NewBitArray(steeringMemberCandidateSet.Size()),
		Vrfs:                       vrfs,
	}
}
// Check all vrfs in the VRF list
func (vrfSet *VrfSet) ValidateBasic() error {
	for idx, vrf := range vrfSet.Vrfs {
		if err := vrf.ValidateBasic(); err != nil {
			return fmt.Errorf("Invalid vrf #%d: %w", idx, err)
		}
	}

	return nil
}

// Get the height of the blockchain of the current VRF set
func (vrfSet *VrfSet) GetHeight() int64 {
	if vrfSet == nil {
		return 0
	}
	return vrfSet.Height
}

// Get vrfs size in the VrfSet
func (vrfSet *VrfSet) Size() int {
	if vrfSet == nil {
		return 0
	}

	return len(vrfSet.Vrfs)
}

func (vrfSet *VrfSet) IsNilOrEmpty() bool {
	return vrfSet == nil || len(vrfSet.Vrfs) == 0
}

// Function to transform from ProtocolBuffer to the corresponding type
func VrfSetFromProto(vrfSetProto *tmproto.VrfSet) (*VrfSet, error) {
	if vrfSetProto == nil {
		return nil, errors.New("nil vrf set")
	}

	vrfSet := new(VrfSet)
	vrfs := make([]*Vrf, len(vrfSetProto.Vrfs))

	for i, vrfProto := range vrfSetProto.Vrfs {
		vrf := VrfFromProto(vrfProto)
		if vrf == nil {
			return nil, errors.New(nilVrfStr)
		}
		vrfs[i] = vrf
	}
	vrfSet.Height = vrfSetProto.Height
	vrfSet.VrfsBitArray = bits.NewBitArray(len(vrfs))
	vrfSet.Vrfs = vrfs

	return vrfSet, vrfSet.ValidateBasic()
}

// Add vrf in the set
func (vrfSet *VrfSet) AddVrf(vrf *Vrf) bool {
	vrfSet.mtx.Lock()
	defer func() {
		vrfSet.mtx.Unlock()
	}()


	if vrf == nil {
		return false
	}

	if vrf.Value != nil {
		if vrf.Verify() == false {
			return false
		}
	
		if vrfSet.Height != vrf.Height {
			return false
		}
	}

	vrfIndex := vrfSet.GetVrfIndexByAddress(vrf.SteeringMemberCandidatePubKey.Address())
	
	if vrfIndex == -1 {
		return false
	}

	if vrfSet.VrfsBitArray.GetIndex(int(vrfIndex)) == false {
		vrf.VrfIndex = vrfIndex
		vrfSet.Vrfs[vrfIndex] = vrf.Copy()
		
		// Set the vrf flag (true == existing)
		vrfSet.VrfsBitArray.SetIndex(int(vrfIndex), true)
		return true
	}

	return false
}

// Get vrf with steering member candidate public key
func (vrfSet *VrfSet) GetVrf(steeringMemberCandidatePubKey crypto.PubKey) (vrf *Vrf) {
	vrfIndex := vrfSet.GetVrfIndexByAddress(steeringMemberCandidatePubKey.Address())

	if vrfIndex != -1 {
		return vrfSet.Vrfs[vrfIndex]
	}

	return nil
}

// Generate vrf list Hash
func (vrfSet *VrfSet) Hash() []byte {
	vrfBytesArray := make([][]byte, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf != nil && vrf.Proof != nil {
			vrfBytesArray[i] = vrf.GetVrfBytes()
		}
	}
	return merkle.HashFromByteSlices(vrfBytesArray)
}

// Function to transform from current set (type) to the corresponding ProtocolBuffer
func (vrfSet *VrfSet) ToProto() (*tmproto.VrfSet, error) {
	if vrfSet.IsNilOrEmpty() {
		return &tmproto.VrfSet{}, nil
	}

	vrfSetProto := new(tmproto.VrfSet)
	vrfsProto := make([]*tmproto.Vrf, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		vrfProto := vrf.ToProto()
		if vrfProto != nil {
			vrfsProto[i] = vrfProto
		}
	}
	vrfSetProto.Height = vrfSet.Height
	vrfSetProto.Vrfs = vrfsProto

	return vrfSetProto, nil
}

func (vrfSet *VrfSet) Copy() *VrfSet {
	if vrfSet == nil {
		return nil
	}

	vrfsCopy := make([]*Vrf, len(vrfSet.Vrfs))
	for i, vrf := range vrfSet.Vrfs {
		if vrf != nil {
			vrfsCopy[i] = vrf.Copy()
		}
	}

	return &VrfSet{
		Height:                     vrfSet.Height,
		Vrfs:                       vrfsCopy,
		VrfsBitArray:               vrfSet.VrfsBitArray.Copy(),
	}
}

func (vrfSet *VrfSet) GetByIndex(vrfIndex int32) *Vrf {
	if vrfSet == nil {
		return nil
	}

	return vrfSet.Vrfs[vrfIndex]
}

func (vrfSet *VrfSet) GetVrfIndexByAddress(address []byte) (int32) {
	if vrfSet == nil {
		return -1
	}
	for idx, vrf := range vrfSet.Vrfs {
		if bytes.Equal(vrf.SteeringMemberCandidatePubKey.Address(), address) {
			return int32(idx)
		}
	}
	return -1
}

func (vrfSet *VrfSet) BitArray() *bits.BitArray {
	if vrfSet == nil {
		return nil
	}

	vrfSet.mtx.Lock()
	defer func() {
		vrfSet.mtx.Unlock()
	}()

	return vrfSet.VrfsBitArray.Copy()
}

type VrfSetReader interface {
	GetHeight() int64
	Size() int
	BitArray() *bits.BitArray
	GetByIndex(int32) *Vrf
}

// update the vrf set, when a SDK sends the changed steering member candidate information and initialize.
// It adds or removes the vrf from the managed list
func (vrfSet *VrfSet) UpdateWithChangeSet(steeringMemberCandidateSet *SteeringMemberCandidateSet) error {
	vrfSet.mtx.Lock()
	defer func() {
		vrfSet.mtx.Unlock()
	}()

	vrfs := make([]*Vrf, 0, len(steeringMemberCandidateSet.SteeringMemberCandidates))

	for _, steeringMemberCandidate := range steeringMemberCandidateSet.SteeringMemberCandidates {
		vrf := vrfSet.GetVrf(steeringMemberCandidate.PubKey)

		if vrf != nil {
			vrfs = append(vrfs, vrf.Copy())
		}
	}

	vrfsBitArray := bits.NewBitArray(len(vrfs))
	
	for i, vrf := range vrfs {
		vrfIndex := vrfSet.GetVrfIndexByAddress(vrf.SteeringMemberCandidatePubKey.Address())
		vrfsBitArray.SetIndex(i, vrfSet.VrfsBitArray.GetIndex(int(vrfIndex)))
	}

	vrfSet.Vrfs = vrfs[:]
	vrfSet.VrfsBitArray = vrfsBitArray.Copy()

	return nil
}

// Check an address is included in vrf set at the consensus round.
// It means that the address is whether it is steering member candidate or not steering member candidate.
func (vrfSet *VrfSet) HasAddress(address []byte) (bool) {
	if vrfSet == nil {
		return false
	}
	for _, vrf := range vrfSet.Vrfs {
		if bytes.Equal(vrf.SteeringMemberCandidatePubKey.Address(), address) {
			return true
		}
	}
	return false
}

// Get steering member addresses selected as validators for the next consensus round.
func (vrfSet *VrfSet) GetSteeringMemberAddresses() *SettingSteeringMember {
	if len(vrfSet.Vrfs) != 0 {
		sort.Sort(SortedVrfs(vrfSet.Vrfs))
		var steeringMemberSize int
		if len(vrfSet.Vrfs) < MAXIMUM_STEERING_MEMBER_CANDIDATES {
			steeringMemberSize = len(vrfSet.Vrfs)
		} else {
			steeringMemberSize = MAXIMUM_STEERING_MEMBER_CANDIDATES
		}

		settingSteeringMember := NewSettingSteeringMember(steeringMemberSize)

		for i := 0; i < steeringMemberSize; i++ {
			if (vrfSet.Vrfs[i].Value != nil) {
				settingSteeringMember.SteeringMemberAddresses = append(settingSteeringMember.SteeringMemberAddresses, vrfSet.Vrfs[i].SteeringMemberCandidatePubKey.Address())
			}
		}

		return settingSteeringMember
	}
	return nil
}

func VrfSetFromExistingVrfs(vrfs []*Vrf) (*VrfSet, error) {
	if len(vrfs) == 0 {
		return nil, errors.New("vrf set is empty")
	}
	for _, val := range vrfs {
		err := val.ValidateBasic()
		if err != nil {
			return nil, fmt.Errorf("can't create vrf set: %w", err)
		}
	}

	vrfSet := &VrfSet{
		Vrfs: vrfs,
	}
	return vrfSet, nil
}