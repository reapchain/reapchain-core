package types

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.reappay.net/sucs-lab//reapchain/crypto"
	ce "gitlab.reappay.net/sucs-lab//reapchain/crypto/encoding"

	tmproto "gitlab.reappay.net/sucs-lab//reapchain/proto/reapchain/types"
)

// Volatile state for each StandingMember
// NOTE: The ProposerPriority is not included in StandingMember.Hash();
// make sure to update that method if changes are made here
type StandingMember struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
}

// NewStandingMember returns a new standing member with the given pubkey and voting power.
func NewStandingMember(pubKey crypto.PubKey) *StandingMember {
	return &StandingMember{
		Address: pubKey.Address(),
		PubKey:  pubKey,
	}
}

// ValidateBasic performs basic validation.
func (sm *StandingMember) ValidateBasic() error {
	if sm == nil {
		return errors.New("nil standing member")
	}
	if sm.PubKey == nil {
		return errors.New("standing member does not have a public key")
	}

	if len(sm.Address) != crypto.AddressSize {
		return fmt.Errorf("standing member address is the wrong size: %v", sm.Address)
	}

	return nil
}

func (sm *StandingMember) Bytes() []byte {
	pk, err := ce.PubKeyToProto(sm.PubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleStandingMember{
		PubKey: &pk,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func StandingMemberFromProto(vp *tmproto.StandingMember) (*StandingMember, error) {
	if vp == nil {
		return nil, errors.New("nil standing member")
	}

	pk, err := ce.PubKeyFromProto(vp.PubKey)
	if err != nil {
		return nil, err
	}
	v := new(StandingMember)
	v.Address = vp.GetAddress()
	v.PubKey = pk

	return v, nil
}

func (sm *StandingMember) ToProto() (*tmproto.StandingMember, error) {
	if sm == nil {
		return nil, errors.New("nil standing member")
	}

	pk, err := ce.PubKeyToProto(sm.PubKey)
	if err != nil {
		return nil, err
	}

	vp := tmproto.StandingMember{
		Address: sm.Address,
		PubKey:  pk,
	}

	return &vp, nil
}

func (sm *StandingMember) Copy() *StandingMember {
	smCopy := *sm
	return &smCopy
}

func StandingMemberListString(vals []*StandingMember) string {
	chunks := make([]string, len(vals))
	for i, val := range vals {
		chunks[i] = fmt.Sprintf("%s", val.Address)
	}

	return strings.Join(chunks, ",")
}
