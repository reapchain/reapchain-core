package v0

import (
	"fmt"

	"github.com/reapchain/reapchain-core/p2p"
)

type peerError struct {
	err    error
	peerID p2p.ID
}

func (e peerError) Error() string {
	return fmt.Sprintf("error with peer %v: %s", e.peerID, e.err.Error())
}
