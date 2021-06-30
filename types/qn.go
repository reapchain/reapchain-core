package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"gitlab.reappay.net/sucs-lab/reapchain/crypto"
	tmbytes "gitlab.reappay.net/sucs-lab/reapchain/libs/bytes"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
// type Qns []Qn

// func (qns Qns) Hash() []byte {
// 	bzs := make([][]byte, len(qns))
// 	for i, qn := range qns {
// 		bzs[i] = qn.Bytes()
// 	}
// 	return merkle.HashFromByteSlices(bzs)
// }
type QnData struct {

	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Qns QnSet `json:"qns"`

	// Volatile
	hash tmbytes.HexBytes
}

func (qnData *QnData) Hash() tmbytes.HexBytes {
	if qnData == nil {
		return ([Qns{}).Hash()
	}
	if qnData.hash == nil {
		qnData.hash = qnData.Qns.Hash() // NOTE: leaves of merkle tree are TxIDs
	}
	return qnData.hash
}

type Qn struct {
	Height                int64     `json:"height"`
	Timestamp             time.Time `json:"timestamp"`
	StandingMemberAddress Address   `json:"standing_member_address"`
	StandingMemberIndex   int32     `json:"standing_member_index"`
	Value                 uint64    `json:"value"`
	Signature             []byte    `json:"signature"`
}

// func NewQn(pubKey crypto.PubKey, value uint64, height int64, signature []byte) *Qn {
// 	height
// 	BlockID
// 	Timestamp
// 	StandingMemberAddress
// 	StandingMemberIndex
// 	Value
// 	Signature

// 	return &Qn{
// 		Address:   pubKey.Address(),
// 		PubKey:    pubKey,
// 		Value:     value,
// 		Height:    height,
// 		Signature: signature,
// 	}
// }

func (qn *Qn) ValidateBasic() error {
	if qn.Height < 0 {
		return errors.New("negative Height")
	}

	if len(qn.StandingMemberAddress) != crypto.AddressSize {
		return fmt.Errorf("expected StandingMemberAddress size to be %d bytes, got %d bytes",
			crypto.AddressSize,
			len(qn.StandingMemberAddress),
		)
	}

	if qn.StandingMemberIndex < 0 {
		return errors.New("negative StandingMemberIndex")
	}

	if len(qn.Signature) == 0 {
		return errors.New("signature is missing")
	}

	if len(qn.Signature) > MaxSignatureSize {
		return fmt.Errorf("signature is too big (max: %d)", MaxSignatureSize)
	}

	//TODO: stompesi -
	return nil
}

// func (qn *Qn) Bytes() []byte {
// 	pk, err := ce.PubKeyToProto(qn.PubKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	pbv := tmproto.SimpleQn{
// 		PubKey: &pk,
// 		Value:  qn.Value,
// 	}

// 	bz, err := pbv.Marshal()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return bz
// }

// func QnListString(qns []*Qn) string {
// 	chunks := make([]string, len(qns))
// 	for i, qn := range qns {
// 		chunks[i] = fmt.Sprintf("%s:%d", qn.Address, qn.Value)
// 	}

// 	return strings.Join(chunks, ",")
// }

func (qn *Qn) Verify(pubKey crypto.PubKey) error {
	if !bytes.Equal(pubKey.Address(), qn.StandingMemberAddress) {
		return ErrQnInvalidStandingMemberAddress
	}
	qnProto := qn.ToProto()
	if !pubKey.VerifySignature(QnValueToBytes(qnProto.Value), qn.Signature) {
		return ErrQnInvalidSignature
	}
	return nil
}

func QnValueToBytes(value uint64) []byte {
	qnValueBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(qnValueBytes, value)

	return qnValueBytes
}
