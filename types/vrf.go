package types

import (
	"errors"
	"fmt"
	"time"

	abci "github.com/reapchain/reapchain-core/abci/types"
	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	tmbytes "github.com/reapchain/reapchain-core/libs/bytes"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
	"github.com/reapchain/reapchain-core/vrfunc"
)

const (
	nilVrfStr string = "nil-Vrf"
)

type Vrf struct {
	Height                        int64         `json:"height"`
	Timestamp                     time.Time     `json:"timestamp"`
	SteeringMemberCandidatePubKey crypto.PubKey `json:"steering_member_candidate_pub_key"`
	SteeringMemberCandidateIndex  int32         `json:"steering_member_candidate_index"`
	Value                         []byte        `json:"value"`
	Proof                         []byte        `json:"proof"`
	Seed                          []byte        `json:"seed"`
}

func NewVrf(height int64, steeringMemberCandidatePubKey crypto.PubKey, seed []byte) *Vrf {
	vrf := Vrf{
		Height:                        height,
		Timestamp:                     time.Now(),
		SteeringMemberCandidatePubKey: steeringMemberCandidatePubKey,
		Seed:                          seed,
	}

	return &vrf
}

func NewVrfAsEmpty(height int64, steeringMemberCandidatePubKey crypto.PubKey) *Vrf {
	vrf := Vrf{
		Height:                        height,
		Timestamp:                     time.Unix(0, 0),
		SteeringMemberCandidatePubKey: steeringMemberCandidatePubKey,
	}

	return &vrf
}

func (vrf *Vrf) ValidateBasic() error {
	if vrf.Height < 0 {
		return errors.New("negative Height")
	}

	steeringMemberCandidateAddress := vrf.SteeringMemberCandidatePubKey.Address()
	if len(steeringMemberCandidateAddress) != crypto.AddressSize {
		return fmt.Errorf("expected SteeringMemberCandidateAddress size to be %d bytes, got %d bytes",
			crypto.AddressSize,
			len(steeringMemberCandidateAddress),
		)
	}

	return nil
}

func (vrf *Vrf) Copy() *Vrf {
	vrfCopy := *vrf
	return &vrfCopy
}

func (vrf *Vrf) String() string {
	if vrf == nil {
		return nilVrfStr
	}
	return fmt.Sprintf("Vrf{%X %v %X @ %s}",
		tmbytes.Fingerprint(vrf.SteeringMemberCandidatePubKey.Address()),
		vrf.Height,
		tmbytes.Fingerprint(vrf.Proof),
		CanonicalTime(vrf.Timestamp),
	)
}

func (vrf *Vrf) GetVrfBytesForSign() []byte {
	if vrf == nil {
		return nil
	}

	pubKeyProto, err := ce.PubKeyToProto(vrf.SteeringMemberCandidatePubKey)
	if err != nil {
		panic(err)
	}

	vrfProto := tmproto.Vrf{
		Height:                        vrf.Height,
		Timestamp:                     vrf.Timestamp,
		SteeringMemberCandidatePubKey: pubKeyProto,
		Value:                         vrf.Value,
	}

	VrfSignBytes, err := vrfProto.Marshal()
	if err != nil {
		panic(err)
	}
	return VrfSignBytes
}

func (vrf *Vrf) GetVrfBytes() []byte {
	if vrf == nil {
		return nil
	}

	pubKeyProto, err := ce.PubKeyToProto(vrf.SteeringMemberCandidatePubKey)
	if err != nil {
		panic(err)
	}

	vrfProto := tmproto.Vrf{
		Height:                        vrf.Height,
		Timestamp:                     vrf.Timestamp,
		SteeringMemberCandidatePubKey: pubKeyProto,
		Value:                         vrf.Value,
		Proof:                         vrf.Proof,
	}

	vrfSignBytes, err := vrfProto.Marshal()
	if err != nil {
		panic(err)
	}
	return vrfSignBytes
}

func (vrf *Vrf) Verify() bool {
	publicKey := vrfunc.PublicKey(vrf.SteeringMemberCandidatePubKey.Bytes())

	return publicKey.Verify(vrf.Seed, vrf.Value, vrf.Proof)
}

func (vrf *Vrf) ToProto() *tmproto.Vrf {
	if vrf == nil {
		return nil
	}

	pubKey, err := ce.PubKeyToProto(vrf.SteeringMemberCandidatePubKey)
	if err != nil {
		return nil
	}

	vrfProto := tmproto.Vrf{
		Height:                        vrf.Height,
		Timestamp:                     vrf.Timestamp,
		SteeringMemberCandidatePubKey: pubKey,
		SteeringMemberCandidateIndex:  vrf.SteeringMemberCandidateIndex,
		Value:                         vrf.Value,
		Proof:                         vrf.Proof,
		Seed:                          vrf.Seed,
	}

	return &vrfProto
}

func VrfFromProto(vrfProto *tmproto.Vrf) *Vrf {
	if vrfProto == nil {
		return nil
	}

	pubKey, err := ce.PubKeyFromProto(vrfProto.SteeringMemberCandidatePubKey)
	if err != nil {
		return nil
	}

	vrf := new(Vrf)
	vrf.Height = vrfProto.Height
	vrf.Timestamp = vrfProto.Timestamp
	vrf.SteeringMemberCandidatePubKey = pubKey
	vrf.SteeringMemberCandidateIndex = vrfProto.SteeringMemberCandidateIndex
	vrf.Value = vrfProto.Value
	vrf.Proof = vrfProto.Proof
	vrf.Seed = vrfProto.Seed

	return vrf
}

func VrfFromAbci(vrfUpdate *abci.VrfUpdate) *Vrf {
	if vrfUpdate == nil {
		return nil
	}

	pubKey, err := ce.PubKeyFromProto(vrfUpdate.SteeringMemberCandidatePubKey)
	if err != nil {
		return nil
	}

	vrf := new(Vrf)
	vrf.Height = vrfUpdate.Height
	vrf.Timestamp = vrfUpdate.Timestamp
	vrf.SteeringMemberCandidatePubKey = pubKey
	vrf.Value = vrfUpdate.Value
	vrf.Proof = vrfUpdate.Proof

	return vrf
}
