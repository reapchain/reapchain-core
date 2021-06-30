package types

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.reappay.net/sucs-lab//reapchain/crypto"
	ce "gitlab.reappay.net/sucs-lab//reapchain/crypto/encoding"

	tmproto "gitlab.reappay.net/sucs-lab//reapchain/proto/reapchain/types"
)

// Volatile state for each SteeringMemberCandidate
// NOTE: The ProposerPriority is not included in SteeringMemberCandidate.Hash();
// make sure to update that method if changes are made here
type SteeringMemberCandidate struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
}

// NewSteeringMemberCandidate returns a new steering member candidate with the given pubkey and voting power.
func NewSteeringMemberCandidate(pubKey crypto.PubKey) *SteeringMemberCandidate {
	return &SteeringMemberCandidate{
		Address: pubKey.Address(),
		PubKey:  pubKey,
	}
}

// ValidateBasic performs basic validation.
func (sm *SteeringMemberCandidate) ValidateBasic() error {
	if sm == nil {
		return errors.New("nil steering member candidate")
	}
	if sm.PubKey == nil {
		return errors.New("steering member candidate does not have a public key")
	}

	if len(sm.Address) != crypto.AddressSize {
		return fmt.Errorf("steering member candidate address is the wrong size: %v", sm.Address)
	}

	return nil
}

func (sm *SteeringMemberCandidate) Bytes() []byte {
	pk, err := ce.PubKeyToProto(sm.PubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleSteeringMemberCandidate{
		PubKey: &pk,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func SteeringMemberCandidateFromProto(vp *tmproto.SteeringMemberCandidate) (*SteeringMemberCandidate, error) {
	if vp == nil {
		return nil, errors.New("nil steering member candidate")
	}

	pk, err := ce.PubKeyFromProto(vp.PubKey)
	if err != nil {
		return nil, err
	}
	v := new(SteeringMemberCandidate)
	v.Address = vp.GetAddress()
	v.PubKey = pk

	return v, nil
}

func (sm *SteeringMemberCandidate) ToProto() (*tmproto.SteeringMemberCandidate, error) {
	if sm == nil {
		return nil, errors.New("nil steering member candidate")
	}

	pk, err := ce.PubKeyToProto(sm.PubKey)
	if err != nil {
		return nil, err
	}

	vp := tmproto.SteeringMemberCandidate{
		Address: sm.Address,
		PubKey:  pk,
	}

	return &vp, nil
}

func (sm *SteeringMemberCandidate) Copy() *SteeringMemberCandidate {
	smCopy := *sm
	return &smCopy
}

func SteeringMemberCandidateListString(vals []*SteeringMemberCandidate) string {
	chunks := make([]string, len(vals))
	for i, val := range vals {
		chunks[i] = fmt.Sprintf("%s", val.Address)
	}

	return strings.Join(chunks, ",")
}
