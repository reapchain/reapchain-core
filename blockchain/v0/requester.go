package v0

import (
	"github.com/reapchain/reapchain-core/p2p"
)

type Request struct {
	Height int64
	PeerID p2p.ID
}
