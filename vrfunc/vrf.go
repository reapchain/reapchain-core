// Package ed25519 implements a verifiable random function using the Edwards form
// of Curve25519, SHA512 and the Elligator map.
//
//     E is Curve25519 (in Edwards coordinates), h is SHA512.
//     f is the elligator map (bytes->E) that covers half of E.
//     8 is the cofactor of E, the group order is 8*l for prime l.
//     Setup : the prover publicly commits to a public key (P : E)
//     H : names -> E
//         H(n) = f(h(n))^8
//     VRF : keys -> names -> vrfs
//         VRF_x(n) = h(n, H(n)^x))
//     Prove : keys -> names -> proofs
//         Prove_x(n) = tuple(c=h(n, g^r, H(n)^r), t=r-c*x, ii=H(n)^x)
//             where r = h(x, n) is used as a source of randomness
//     Check : E -> names -> vrfs -> proofs -> bool
//         Check(P, n, vrf, (c,t,ii)) = vrf == h(n, ii)
//                                     && c == h(n, g^t*P^c, H(n)^t*ii^c)
package vrfunc

import (
	"bytes"
	"crypto/sha512"
	"errors"

	"github.com/reapchain/reapchain-core/crypto/edwards25519"
	"github.com/reapchain/reapchain-core/crypto/extra25519"

	"golang.org/x/crypto/ed25519"
)

const (
	PublicKeySize    = ed25519.PublicKeySize
	PrivateKeySize   = ed25519.PrivateKeySize
	Size             = 32
	intermediateSize = ed25519.PublicKeySize
	ProofSize        = 32 + 32 + intermediateSize
)

var (
	ErrGetPubKey = errors.New("[vrf] Couldn't get corresponding public-key from private-key")
)

type PrivateKey []byte
type PublicKey []byte

func (sk PrivateKey) expandSecret() (x, skhr *[32]byte) {
	x, skhr = new([32]byte), new([32]byte)
	skh := sha512.Sum512(sk[:32])
	copy(x[:], skh[:])
	copy(skhr[:], skh[32:])
	x[0] &= 248
	x[31] &= 127
	x[31] |= 64
	return
}

func hashToCurve(m []byte) *edwards25519.ExtendedGroupElement {
	// H(n) = (f(h(n))^8)
	hmbH := sha512.Sum512(m)
	var hmb [32]byte
	copy(hmb[:], hmbH[:])
	var hm edwards25519.ExtendedGroupElement
	extra25519.HashToEdwards(&hm, &hmb)
	edwards25519.GeDouble(&hm, &hm)
	edwards25519.GeDouble(&hm, &hm)
	edwards25519.GeDouble(&hm, &hm)
	return &hm
}

// Prove returns the vrf value and a proof such that
// Verify(m, vrf, proof) == true. The vrf value is the
// same as returned by Compute(m).
func (sk PrivateKey) Prove(m []byte) (vrf, proof []byte) {
	x, skhr := sk.expandSecret()
	var sH, rH [64]byte
	var r, s, minusS, t, gB, grB, hrB, hxB, hB [32]byte
	var ii, gr, hr edwards25519.ExtendedGroupElement

	h := hashToCurve(m)
	h.ToBytes(&hB)
	edwards25519.GeScalarMult(&ii, x, h)
	ii.ToBytes(&hxB)

	// use hash of private-, public-key and msg as randomness source:
	hash := sha512.New()
	hash.Write(skhr[:])
	hash.Write(sk[32:]) // public key, as in ed25519
	hash.Write(m)
	hash.Sum(rH[:0])
	hash.Reset()
	edwards25519.ScReduce(&r, &rH)

	edwards25519.GeScalarMultBase(&gr, &r)
	edwards25519.GeScalarMult(&hr, &r, h)
	gr.ToBytes(&grB)
	hr.ToBytes(&hrB)
	gB = edwards25519.BaseBytes

	// H2(g, h, g^x, h^x, g^r, h^r, m)
	hash.Write(gB[:])
	hash.Write(hB[:])
	hash.Write(sk[32:]) // ed25519 public-key
	hash.Write(hxB[:])
	hash.Write(grB[:])
	hash.Write(hrB[:])
	hash.Write(m)
	hash.Sum(sH[:0])
	hash.Reset()
	edwards25519.ScReduce(&s, &sH)

	edwards25519.ScNeg(&minusS, &s)
	edwards25519.ScMulAdd(&t, x, &minusS, &r)

	proof = make([]byte, ProofSize)
	copy(proof[:32], s[:])
	copy(proof[32:64], t[:])
	copy(proof[64:96], hxB[:])

	hash.Reset()
	hash.Write(hxB[:])
	hash.Write(m)
	vrf = hash.Sum(nil)[:Size]
	return
}

// Verify returns true iff vrf=Compute(m) for the sk that
// corresponds to pk.
func (pkBytes PublicKey) Verify(m, vrfBytes, proof []byte) bool {
	if len(proof) != ProofSize || len(vrfBytes) != Size || len(pkBytes) != PublicKeySize {
		return false
	}
	var pk, s, sRef, t, vrf, hxB, hB, gB, ABytes, BBytes [32]byte
	copy(vrf[:], vrfBytes)
	copy(pk[:], pkBytes[:])
	copy(s[:32], proof[:32])
	copy(t[:32], proof[32:64])
	copy(hxB[:], proof[64:96])

	hash := sha512.New()
	hash.Write(hxB[:]) // const length
	hash.Write(m)
	if !bytes.Equal(hash.Sum(nil)[:Size], vrf[:Size]) {
		return false
	}
	hash.Reset()

	var P, B, ii, iic edwards25519.ExtendedGroupElement
	var A, hmtP, iicP edwards25519.ProjectiveGroupElement
	if !P.FromBytesBaseGroup(&pk) {
		return false
	}
	if !ii.FromBytesBaseGroup(&hxB) {
		return false
	}
	edwards25519.GeDoubleScalarMultVartime(&A, &s, &P, &t)
	A.ToBytes(&ABytes)
	gB = edwards25519.BaseBytes

	h := hashToCurve(m) // h = H1(m)
	h.ToBytes(&hB)
	edwards25519.GeDoubleScalarMultVartime(&hmtP, &t, h, &[32]byte{})
	edwards25519.GeDoubleScalarMultVartime(&iicP, &s, &ii, &[32]byte{})
	iicP.ToExtended(&iic)
	hmtP.ToExtended(&B)
	edwards25519.GeAdd(&B, &B, &iic)
	B.ToBytes(&BBytes)

	var sH [64]byte
	// sRef = H2(g, h, g^x, v, g^t路G^s,H1(m)^t路v^s, m), with v=H1(m)^x=h^x
	hash.Write(gB[:])
	hash.Write(hB[:])
	hash.Write(pkBytes)
	hash.Write(hxB[:])
	hash.Write(ABytes[:]) // const length (g^t*G^s)
	hash.Write(BBytes[:]) // const length (H1(m)^t*v^s)
	hash.Write(m)
	hash.Sum(sH[:0])

	edwards25519.ScReduce(&sRef, &sH)
	return sRef == s
}
