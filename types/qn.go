package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reapchain/reapchain/crypto"
	ce "github.com/reapchain/reapchain/crypto/encoding"
	tmproto "github.com/reapchain/reapchain/proto/reapchain/types"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
type Qn struct {
	Address Address       `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
	Value   uint64        `json:"value"`
}

func NewQn(pubKey crypto.PubKey, value uint64) *Qn {
	return &Qn{
		Address: pubKey.Address(),
		PubKey:  pubKey,
		Value:   value,
	}
}

func (qn *Qn) ValidateBasic() error {
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

func (qn *Qn) Bytes() []byte {
	pk, err := ce.PubKeyToProto(qn.PubKey)
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleQn{
		PubKey: &pk,
		Value:  qn.Value,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func QnListString(qns []*Qn) string {
	chunks := make([]string, len(qns))
	for i, qn := range qns {
		chunks[i] = fmt.Sprintf("%s:%d", qn.Address, qn.Value)
	}

	return strings.Join(chunks, ",")
}
