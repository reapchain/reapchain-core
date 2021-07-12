package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"
	"github.com/reapchain/reapchain-core/crypto/merkle"
	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
type Vrfs []Vrf

type Vrf struct {
	SteeringMemberCandidatePubKey crypto.PubKey `json:"steering_member_candidate_pub_key"`
	Value                         uint64        `json:"value"`
	Proof                         []byte        `json:"proof"`
}

func (vrfs Vrfs) Hash() []byte {
	vrfz := make([][]byte, len(vrfs))
	for i, vrf := range vrfs {
		vrfz[i] = vrf.Bytes()
	}
	return merkle.HashFromByteSlices(vrfz)
}

//TODO: stompesi
func NewVrf(pubKey crypto.PubKey, value uint64) *Vrf {
	return &Vrf{
		SteeringMemberCandidatePubKey: pubKey,
		Value:                         value,
	}
}

func (qrn *Vrf) ValidateBasic() error {
	steeringMemberCandidateAddress := qrn.SteeringMemberCandidatePubKey.Address()
	if qrn == nil {
		return errors.New("nil qrn")
	}

	if qrn.SteeringMemberCandidatePubKey == nil {
		return errors.New("qrn does not have a public key")
	}

	if len(steeringMemberCandidateAddress) != crypto.AddressSize {
		return fmt.Errorf("qrn of member address is the wrong size: %v", steeringMemberCandidateAddress)
	}

	//TODO: stompesi -
	return nil
}

func (qrn *Vrf) Bytes() []byte {
	pk, err := ce.PubKeyToProto(qrn.SteeringMemberCandidatePubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleVrf{
		PubKey: &pk,
		Value:  qrn.Value,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func VrfListString(qrns []*Vrf) string {
	chunks := make([]string, len(qrns))
	for i, qrn := range qrns {
		chunks[i] = fmt.Sprintf("%s:%d", qrn.SteeringMemberCandidatePubKey.Address(), qrn.Value)
	}

	return strings.Join(chunks, ",")
}

func VrfFromProto(pv *tmproto.Vrf) (*Vrf, error) {
	if pv == nil {
		return nil, errors.New("nil vrf")
	}

	vrf := new(Vrf)

	return vrf, vrf.ValidateBasic()
}
