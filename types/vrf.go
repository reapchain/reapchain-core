package types

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.reappay.net/sucs-lab/reapchain/crypto"
	ce "gitlab.reappay.net/sucs-lab/reapchain/crypto/encoding"
	"gitlab.reappay.net/sucs-lab/reapchain/crypto/merkle"
	tmproto "gitlab.reappay.net/sucs-lab/reapchain/proto/reapchain/types"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
type Vrfs []Vrf

func (qns Vrfs) Hash() []byte {
	bzs := make([][]byte, len(qns))
	for i, qn := range qns {
		bzs[i] = qn.Bytes()
	}
	return merkle.HashFromByteSlices(bzs)
}

type Vrf struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
	Value   uint64        `json:"value"`
}

func NewVrf(pubKey crypto.PubKey, value uint64) *Vrf {
	return &Vrf{
		Address: pubKey.Address(),
		PubKey:  pubKey,
		Value:   value,
	}
}

func (qn *Vrf) ValidateBasic() error {
	if qn == nil {
		return errors.New("nil qn")
	}

	if qn.PubKey == nil {
		return errors.New("qn does not have a public key")
	}

	if len(qn.Address) != crypto.AddressSize {
		return fmt.Errorf("qn of member address is the wrong size: %v", qn.Address)
	}

	//TODO: stompesi -
	return nil
}

func (qn *Vrf) Bytes() []byte {
	pk, err := ce.PubKeyToProto(qn.PubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleVrf{
		PubKey: &pk,
		Value:  qn.Value,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func VrfListString(qns []*Vrf) string {
	chunks := make([]string, len(qns))
	for i, qn := range qns {
		chunks[i] = fmt.Sprintf("%s:%d", qn.Address, qn.Value)
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
