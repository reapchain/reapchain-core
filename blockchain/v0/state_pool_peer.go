package v0

import (
	"errors"
	"math"
	"time"

	flow "github.com/reapchain/reapchain-core/libs/flowrate"
	"github.com/reapchain/reapchain-core/libs/log"
	"github.com/reapchain/reapchain-core/p2p"
)

type StatePoolPeer struct {
	didTimeout  bool
	numPending  int32
	height      int64
	base        int64
	pool        *StatePool
	id          p2p.ID
	recvMonitor *flow.Monitor

	timeout *time.Timer

	logger log.Logger
}

func newStatePeer(pool *StatePool, peerID p2p.ID, base int64, height int64) *StatePoolPeer {
	peer := &StatePoolPeer{
		pool:       pool,
		id:         peerID,
		base:       base,
		height:     height,
		numPending: 0,
		logger:     log.NewNopLogger(),
	}
	return peer
}

func (peer *StatePoolPeer) setLogger(l log.Logger) {
	peer.logger = l
}

func (peer *StatePoolPeer) resetMonitor() {
	peer.recvMonitor = flow.New(time.Second, time.Second*40)
	initialValue := float64(minRecvRate) * math.E
	peer.recvMonitor.SetREMA(initialValue)
}

func (peer *StatePoolPeer) resetTimeout() {
	if peer.timeout == nil {
		peer.timeout = time.AfterFunc(peerTimeout, peer.onTimeout)
	} else {
		peer.timeout.Reset(peerTimeout)
	}
}

func (peer *StatePoolPeer) incrPending() {
	if peer.numPending == 0 {
		peer.resetMonitor()
		peer.resetTimeout()
	}
	peer.numPending++
}

func (peer *StatePoolPeer) decrPending(recvSize int) {
	peer.numPending--
	if peer.numPending == 0 {
		peer.timeout.Stop()
	} else {
		peer.recvMonitor.Update(recvSize)
		peer.resetTimeout()
	}
}

func (peer *StatePoolPeer) onTimeout() {
	peer.pool.mtx.Lock()
	defer peer.pool.mtx.Unlock()

	err := errors.New("peer did not send us anything")
	peer.pool.sendError(err, peer.id)
	peer.logger.Error("SendTimeout", "reason", err, "timeout", peerTimeout)
	peer.didTimeout = true
}
