package types

import (
	tmbytes "github.com/reapchain/reapchain-core/libs/bytes"
)

// Tx is an arbitrary byte array.
// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// Might we want types here ?
// type Qrns []Qrn

// func (qrns Qrns) Hash() []byte {
// 	bzs := make([][]byte, len(qrns))
// 	for i, qrn := range qrns {
// 		bzs[i] = qrn.Bytes()
// 	}
// 	return merkle.HashFromByteSlices(bzs)
// }
type QrnData struct {

	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Qrns QrnSet `json:"qrns"`

	// Volatile
	hash tmbytes.HexBytes
}

func (qrnData *QrnData) Hash() tmbytes.HexBytes {
	if qrnData == nil {
		return (&QrnSet{}).Hash()
	}
	if qrnData.hash == nil {
		qrnData.hash = qrnData.Qrns.Hash() // NOTE: leaves of merkle tree are TxIDs
	}
	return qrnData.hash
}
