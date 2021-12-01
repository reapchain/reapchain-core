package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reapchain/reapchain-core/crypto"
	ce "github.com/reapchain/reapchain-core/crypto/encoding"

	tmproto "github.com/reapchain/reapchain-core/proto/reapchain/types"
)

type SteeringMemberCandidate struct {
	PubKey  crypto.PubKey `json:"pub_key"`
	Address Address       `json:"address"`
}

func NewSteeringMemberCandidate(pubKey crypto.PubKey) *SteeringMemberCandidate {
	return &SteeringMemberCandidate{
		Address: pubKey.Address(),
		PubKey:  pubKey,
	}
}

func (steeringMemberCandidate *SteeringMemberCandidate) ValidateBasic() error {
	if steeringMemberCandidate == nil {
		return errors.New("nil steering member candidate")
	}

	if steeringMemberCandidate.PubKey == nil {
		return errors.New("steering member candidate does not have a public key")
	}

	if len(steeringMemberCandidate.Address) != crypto.AddressSize {
		return fmt.Errorf("steering member candidate address is the wrong size: %v", steeringMemberCandidate.Address)
	}

	return nil
}

func (steeringMemberCandidate *SteeringMemberCandidate) Bytes() []byte {
	pubKeyProto, err := ce.PubKeyToProto(steeringMemberCandidate.PubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SteeringMemberCandidate{
		PubKey: pubKeyProto,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func (steeringMemberCandidate *SteeringMemberCandidate) ToProto() (*tmproto.SteeringMemberCandidate, error) {
	if steeringMemberCandidate == nil {
		return nil, errors.New("nil steering member candidate")
	}

	pubKeyProto, err := ce.PubKeyToProto(steeringMemberCandidate.PubKey)
	if err != nil {
		return nil, err
	}

	steeringMemberCandidateProto := tmproto.SteeringMemberCandidate{
		Address: steeringMemberCandidate.Address,
		PubKey:  pubKeyProto,
	}

	return &steeringMemberCandidateProto, nil
}

func (steeringMemberCandidate *SteeringMemberCandidate) Copy() *SteeringMemberCandidate {
	steeringMemberCandidateCopy := *steeringMemberCandidate
	return &steeringMemberCandidateCopy
}

func SteeringMemberCandidateListString(steeringMemberCandidates []*SteeringMemberCandidate) string {
	chunks := make([]string, len(steeringMemberCandidates))
	for i, steeringMemberCandidate := range steeringMemberCandidates {
		chunks[i] = fmt.Sprintf("%s", steeringMemberCandidate.Address)
	}

	return strings.Join(chunks, ",")
}

func SteeringMemberCandidateFromProto(steeringMemberCandidateProto *tmproto.SteeringMemberCandidate) (*SteeringMemberCandidate, error) {
	if steeringMemberCandidateProto == nil {
		return nil, errors.New("nil steering member candidate")
	}

	pubKey, err := ce.PubKeyFromProto(steeringMemberCandidateProto.PubKey)
	if err != nil {
		return nil, err
	}
	steeringMemberCandidate := new(SteeringMemberCandidate)
	steeringMemberCandidate.Address = steeringMemberCandidateProto.GetAddress()
	steeringMemberCandidate.PubKey = pubKey

	return steeringMemberCandidate, nil
}
