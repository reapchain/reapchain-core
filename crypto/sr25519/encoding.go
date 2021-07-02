package sr25519

import (
	"github.com/reapchain/reapchain-core/crypto"
	tmjson "github.com/reapchain/reapchain-core/libs/json"
)

var _ crypto.PrivKey = PrivKey{}

const (
	PrivKeyName = "reapchain/PrivKeySr25519"
	PubKeyName  = "reapchain/PubKeySr25519"

	// SignatureSize is the size of an Edwards25519 signature. Namely the size of a compressed
	// Sr25519 point, and a field element. Both of which are 32 bytes.
	SignatureSize = 64
)

func init() {

	tmjson.RegisterType(PubKey{}, PubKeyName)
	tmjson.RegisterType(PrivKey{}, PrivKeyName)
}
