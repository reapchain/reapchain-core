package v0

import (
	"github.com/reapchain/reapchain-core/p2p"
)

// BlockRequest stores a block request identified by the block Height and the PeerID responsible for
// delivering the block

type BlockRequest struct {
	Height int64
	PeerID p2p.ID
}
