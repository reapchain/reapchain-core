package vrf

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/reapchain/reapchain-core/privval"
)

//privateKey := PrivateKey(pv.Key.PrivKey.Bytes())
//publicKey, err := privateKey.Public()
func TestHonestComplete(t *testing.T) {
	pv := privval.GenFilePV("priv_validator_key.json", "priv_validator_state.json")
	sk := PrivateKey(pv.Key.PrivKey.Bytes())
	pk, err := sk.Public()
	if err != true {

	}

	alice := []byte("alice")
	aliceVRF := sk.Compute(alice)
	aliceVRFFromProof, aliceProof := sk.Prove(alice, false)

	fmt.Printf("pk:           %X\n", pk)
	fmt.Printf("sk:           %X\n", sk)
	fmt.Printf("alice(bytes): %X\n", alice)
	fmt.Printf("aliceVRF:     %X\n", aliceVRF)
	fmt.Printf("aliceProof:   %X\n", aliceProof)
	fmt.Printf("aliceVRFFromProof:   %X\n", aliceVRFFromProof)

	if !pk.Verify(alice, aliceVRF, aliceProof) {
		t.Error("Gen -> Compute -> Prove -> Verify -> FALSE")
	}
	if !bytes.Equal(aliceVRF, aliceVRFFromProof) {
		t.Error("Compute != Prove")
	}
}
