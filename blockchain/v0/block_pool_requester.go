package v0

import (
	"sync/atomic"
	"time"

	"github.com/reapchain/reapchain-core/libs/service"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/types"
)

//bpRequest performs the role of requesting other nodes to validate the check
type bpRequester struct {
	service.BaseService
	pool       *BlockPool
	height     int64
	gotBlockCh chan struct{}
	redoCh     chan p2p.ID // redo may send multitime, add peerId to identify repeat

	mtx    tmsync.Mutex
	peerID p2p.ID

	block *types.Block
}

func newBPRequester(pool *BlockPool, height int64) *bpRequester {
	bpr := &bpRequester{
		pool:       pool,
		height:     height,
		gotBlockCh: make(chan struct{}, 1),
		redoCh:     make(chan p2p.ID, 1),

		peerID: "",
		block:  nil,
	}
	bpr.BaseService = *service.NewBaseService(nil, "bpRequester", bpr)
	return bpr
}

func (bpr *bpRequester) OnStart() error {
	go bpr.requestRoutine()
	return nil
}

// Returns true if the peer matches and block doesn't already exist.
func (bpr *bpRequester) setBlock(block *types.Block, peerID p2p.ID) bool {
	bpr.mtx.Lock()
	if bpr.block != nil || bpr.peerID != peerID {
		bpr.mtx.Unlock()
		return false
	}
	bpr.block = block
	bpr.mtx.Unlock()

	select {
	case bpr.gotBlockCh <- struct{}{}:
	default:
	}
	return true
}

func (bpr *bpRequester) getBlock() *types.Block {
	bpr.mtx.Lock()
	defer bpr.mtx.Unlock()
	return bpr.block
}

func (bpr *bpRequester) getPeerID() p2p.ID {
	bpr.mtx.Lock()
	defer bpr.mtx.Unlock()
	return bpr.peerID
}

// This is called from the requestRoutine, upon redo().
func (bpr *bpRequester) reset() {
	bpr.mtx.Lock()
	defer bpr.mtx.Unlock()

	if bpr.block != nil {
		atomic.AddInt32(&bpr.pool.numPending, 1)
	}

	bpr.peerID = ""
	bpr.block = nil
}

// Tells bpRequester to pick another peer and try again.
// NOTE: Nonblocking, and does nothing if another redo
// was already requested.
func (bpr *bpRequester) redo(peerID p2p.ID) {
	select {
	case bpr.redoCh <- peerID:
	default:
	}
}

// Responsible for making more requests as necessary
// Returns only when a block is found (e.g. AddBlock() is called)
func (bpr *bpRequester) requestRoutine() {
OUTER_LOOP:
	for {
		// Pick a peer to send request to.
		var peer *bpPeer
	PICK_PEER_LOOP:
		for {
			if !bpr.IsRunning() || !bpr.pool.IsRunning() {
				return
			}
			peer = bpr.pool.pickIncrAvailablePeer(bpr.height)
			if peer == nil {
				// log.Info("No peers available", "height", height)
				time.Sleep(requestIntervalMS * time.Millisecond)
				continue PICK_PEER_LOOP
			}
			break PICK_PEER_LOOP
		}
		bpr.mtx.Lock()
		bpr.peerID = peer.id
		bpr.mtx.Unlock()

		// Send request and wait.
		bpr.pool.sendRequest(bpr.height, peer.id)
	WAIT_LOOP:
		for {
			select {
			case <-bpr.pool.Quit():
				if err := bpr.Stop(); err != nil {
					bpr.Logger.Error("Error stopped requester", "err", err)
				}
				return
			case <-bpr.Quit():
				return
			case peerID := <-bpr.redoCh:
				if peerID == bpr.peerID {
					bpr.reset()
					continue OUTER_LOOP
				} else {
					continue WAIT_LOOP
				}
			case <-bpr.gotBlockCh:
				// We got a block!
				// Continue the for-loop and wait til Quit.
				continue WAIT_LOOP
			}
		}
	}
}
