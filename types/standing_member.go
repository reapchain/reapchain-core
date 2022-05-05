package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"

	tmproto "github.com/reapchain/reapchain-core/proto/reapchain-core/types"
)

type StandingMember struct {
	PubKey  crypto.PubKey `json:"pub_key"`
	Address Address       `json:"address"`
	VotingPower int64         `json:"voting_power"`
}

func NewStandingMember(pubKey crypto.PubKey, votingPower int64) *StandingMember {
	return &StandingMember{
		Address: pubKey.Address(),
		PubKey:  pubKey,
		VotingPower:      votingPower,
	}
}

func (standingMember *StandingMember) ValidateBasic() error {
	if standingMember == nil {
		return errors.New("nil standing member")
	}

	if standingMember.PubKey == nil {
		return errors.New("standing member does not have a public key")
	}

	if len(standingMember.Address) != crypto.AddressSize {
		return fmt.Errorf("standing member address is the wrong size: %v", standingMember.Address)
	}

	return nil
}

func (standingMember *StandingMember) Bytes() []byte {
	pubKeyProto, err := ce.PubKeyToProto(standingMember.PubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleStandingMember{
		PubKey: &pubKeyProto,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func StandingMemberFromProto(standingMemberProto *tmproto.StandingMember) (*StandingMember, error) {
	if standingMemberProto == nil {
		return nil, errors.New("nil standing member")
	}

	pubKey, err := ce.PubKeyFromProto(standingMemberProto.PubKey)
	if err != nil {
		return nil, err
	}
	standingMember := new(StandingMember)
	standingMember.Address = standingMemberProto.GetAddress()
	standingMember.PubKey = pubKey

	return standingMember, nil
}

func (standingMember *StandingMember) ToProto() (*tmproto.StandingMember, error) {
	if standingMember == nil {
		return nil, errors.New("nil standing member")
	}

	pubKeyProto, err := ce.PubKeyToProto(standingMember.PubKey)
	if err != nil {
		return nil, err
	}

	standingMemberProto := tmproto.StandingMember{
		Address: standingMember.Address,
		PubKey:  pubKeyProto,
	}

	return &standingMemberProto, nil
}

func (sm *StandingMember) Copy() *StandingMember {
	smCopy := *sm
	return &smCopy
}

func StandingMemberListString(standingMembers []*StandingMember) string {
	chunks := make([]string, len(standingMembers))
	for i, standingMember := range standingMembers {
		chunks[i] = fmt.Sprintf("%s", standingMember.Address)
	}

	return strings.Join(chunks, ",")
}
