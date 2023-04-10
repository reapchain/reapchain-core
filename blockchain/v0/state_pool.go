package v0

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/reapchain/reapchain-core/libs/service"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	"github.com/reapchain/reapchain-core/p2p"
	sm "github.com/reapchain/reapchain-core/state"
)

type StatePool struct {
	service.BaseService
	startTime time.Time

	mtx           tmsync.Mutex
	requesters    map[int64]*StatePoolRequester
	height        int64
	peers         map[p2p.ID]*StatePoolPeer
	maxPeerHeight int64
	numPending    int32

	requestsCh chan<- BlockRequest
	errorsCh   chan<- peerError
}
// StatePool keeps track of the fast sync peers, state requests and state responses.
func NewStatePool(start int64, requestsCh chan<- BlockRequest, errorsCh chan<- peerError) *StatePool {
	statePool := &StatePool{
		peers: make(map[p2p.ID]*StatePoolPeer),

		requesters: make(map[int64]*StatePoolRequester),
		height:     start,
		numPending: 0,

		requestsCh: requestsCh,
		errorsCh:   errorsCh,
	}
	statePool.BaseService = *service.NewBaseService(nil, "StatePool", statePool)
	return statePool
}
// OnStart implements service.Service by spawning requesters routine and recording pool's start time.
func (pool *StatePool) OnStart() error {
	go pool.makeRequestersRoutine()
	pool.startTime = time.Now()
	return nil
}
// spawns requesters as needed
func (pool *StatePool) makeRequestersRoutine() {
	for {
		if !pool.IsRunning() {
			break
		}

		_, numPending, lenRequesters := pool.GetStatus()
		switch {
		case numPending >= maxPendingRequests:
			// sleep for a bit.
			time.Sleep(requestIntervalMS * time.Millisecond)
			// check for timed out peers
			pool.removeTimedoutPeers()
		case lenRequesters >= maxTotalRequesters:
			// sleep for a bit.
			time.Sleep(requestIntervalMS * time.Millisecond)
			// check for timed out peers

		default:
			// request for more blocks.
			pool.makeNextRequester()
		}
	}
}

func (pool *StatePool) removeTimedoutPeers() {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	for _, peer := range pool.peers {
		if !peer.didTimeout && peer.numPending > 0 {
			curRate := peer.recvMonitor.Status().CurRate
			// curRate can be 0 on start
			if curRate != 0 && curRate < minRecvRate {
				err := errors.New("peer is not sending us data fast enough")
				pool.sendError(err, peer.id)
				pool.Logger.Error("SendTimeout", "peer", peer.id,
					"reason", err,
					"curRate", fmt.Sprintf("%d KB/s", curRate/1024),
					"minRate", fmt.Sprintf("%d KB/s", minRecvRate/1024))
				peer.didTimeout = true
			}
		}
		if peer.didTimeout {
			pool.removePeer(peer.id)
		}
	}
}
// GetStatus returns pool's height, numPending requests and the number of
// requesters.
func (pool *StatePool) GetStatus() (height int64, numPending int32, lenRequesters int) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	return pool.height, atomic.LoadInt32(&pool.numPending), len(pool.requesters)
}

// IsCaughtUp returns true if this node is caught up, false - otherwise.
// TODO: relax conditions, prevent abuse.
func (pool *StatePool) IsCaughtUp() bool {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	if len(pool.peers) == 0 {
		pool.Logger.Debug("Blockpool has no peers")
		return false
	}

	// Some conditions to determine if we're caught up.
	// Ensures we've either received a block or waited some amount of time,
	// and that we're synced to the highest known height.
	// Note we use maxPeerHeight - 1 because to sync block H requires block H+1
	// to verify the LastCommit.
	receivedBlockOrTimedOut := pool.height > 0 || time.Since(pool.startTime) > 5*time.Second
	ourChainIsLongestAmongPeers := pool.maxPeerHeight == 0 || pool.height >= (pool.maxPeerHeight-1)
	isCaughtUp := receivedBlockOrTimedOut && ourChainIsLongestAmongPeers
	return isCaughtUp
}

// PeekTwoStates returns states at pool.height and pool.height+1.
// We need to see the second state's Commit to validate the first state.
// So we peek two states at a time.
// The caller will verify the commit.
func (pool *StatePool) PeekTwoStates() (firstState *sm.State, secondState *sm.State) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	if r := pool.requesters[pool.height]; r != nil {
		firstState = r.getState()
	}
	if r := pool.requesters[pool.height+1]; r != nil {
		secondState = r.getState()
	}
	return
}

// PopRequest pops the first state at pool.height.
// It must have been validated by 'second'.Commit from PeekTwoStates().
func (pool *StatePool) PopRequest() {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	if r := pool.requesters[pool.height]; r != nil {
		if err := r.Stop(); err != nil {
			pool.Logger.Error("Error stopping requester", "err", err)
		}
		delete(pool.requesters, pool.height)
		pool.height++
	} else {
		panic(fmt.Sprintf("Expected requester to pop, got nothing at height %v", pool.height))
	}
}

// RedoRequest invalidates the state at pool.height,
// Remove the peer and redo request from others.
// Returns the ID of the removed peer.
func (pool *StatePool) RedoRequest(height int64) p2p.ID {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	request := pool.requesters[height]
	peerID := request.getPeerID()
	if peerID != p2p.ID("") {
		pool.removePeer(peerID)
	}
	return peerID
}

// AddState validates that the state comes from the peer it was expected from and calls the requester to store it.
// TODO: ensure that states come in order for each peer.
func (pool *StatePool) AddState(peerID p2p.ID, state *sm.State, blockSize int) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	requester := pool.requesters[state.LastBlockHeight]
	if requester == nil {
		diff := pool.height - state.LastBlockHeight
		if diff < 0 {
			diff *= -1
		}
		if diff > maxDiffBetweenCurrentAndReceivedBlockHeight {
			pool.sendError(errors.New("peer sent us a state we didn't expect with a height too far ahead/behind"), peerID)
		}
		return
	}

	if requester.setState(state, peerID) {
		atomic.AddInt32(&pool.numPending, -1)
		peer := pool.peers[peerID]
		if peer != nil {
			peer.decrPending(blockSize)
		}
	} else {
		pool.Logger.Info("invalid state peer", "peer", peerID, "blockHeight", state.LastBlockHeight)
		pool.sendError(errors.New("invalid state peer"), peerID)
	}
}

// MaxPeerHeight returns the highest reported height.
func (pool *StatePool) MaxPeerHeight() int64 {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()
	return pool.maxPeerHeight
}

// SetPeerRange sets the peer's alleged blockchain base and height.
func (pool *StatePool) SetPeerRange(peerID p2p.ID, base int64, height int64) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	peer := pool.peers[peerID]
	if peer != nil {
		peer.base = base
		peer.height = height
	} else {
		peer = newStatePeer(pool, peerID, base, height)
		peer.setLogger(pool.Logger.With("peer", peerID))
		pool.peers[peerID] = peer
	}

	if height > pool.maxPeerHeight {
		pool.maxPeerHeight = height
	}
}

// RemovePeer removes the peer with peerID from the pool. 
// If there's no peer with peerID, function is a no-op.
func (pool *StatePool) RemovePeer(peerID p2p.ID) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	pool.removePeer(peerID)
}

func (pool *StatePool) removePeer(peerID p2p.ID) {
	for _, requester := range pool.requesters {
		if requester.getPeerID() == peerID {
			requester.redo(peerID)
		}
	}

	peer, ok := pool.peers[peerID]
	if ok {
		if peer.timeout != nil {
			peer.timeout.Stop()
		}

		delete(pool.peers, peerID)

		// Find a new peer with the biggest height and update maxPeerHeight if the
		// peer's height was the biggest.
		if peer.height == pool.maxPeerHeight {
			pool.updateMaxPeerHeight()
		}
	}
}

// If no peers are left, maxPeerHeight is set to 0.
func (pool *StatePool) updateMaxPeerHeight() {
	var max int64
	for _, peer := range pool.peers {
		if peer.height > max {
			max = peer.height
		}
	}
	pool.maxPeerHeight = max
}

// Pick an available peer with the given height available.
// If no peers are available, returns nil.
func (pool *StatePool) pickIncrAvailablePeer(height int64) *StatePoolPeer {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	for _, peer := range pool.peers {
		if peer.didTimeout {
			pool.removePeer(peer.id)
			continue
		}
		if peer.numPending >= maxPendingRequestsPerPeer {
			continue
		}
		if height < peer.base || height > peer.height {
			continue
		}
		peer.incrPending()
		return peer
	}
	return nil
}

func (pool *StatePool) makeNextRequester() {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	nextHeight := pool.height + pool.requestersLen()
	if nextHeight > pool.maxPeerHeight {
		return
	}

	request := NewStatePoolRequester(pool, nextHeight)

	pool.requesters[nextHeight] = request
	atomic.AddInt32(&pool.numPending, 1)

	err := request.Start()
	if err != nil {
		request.Logger.Error("Error starting request", "err", err)
	}
}

func (pool *StatePool) requestersLen() int64 {
	return int64(len(pool.requesters))
}

func (pool *StatePool) sendRequest(height int64, peerID p2p.ID) {
	if !pool.IsRunning() {
		return
	}
	pool.requestsCh <- BlockRequest{height, peerID}
}

func (pool *StatePool) sendError(err error, peerID p2p.ID) {
	if !pool.IsRunning() {
		return
	}
	pool.errorsCh <- peerError{err, peerID}
}

// for debugging purposes nolint:unused
func (pool *StatePool) debug() string {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	str := ""
	nextHeight := pool.height + pool.requestersLen()
	for h := pool.height; h < nextHeight; h++ {
		if pool.requesters[h] == nil {
			str += fmt.Sprintf("H(%v):X ", h)
		} else {
			str += fmt.Sprintf("H(%v):", h)
			str += fmt.Sprintf("B?(%v) ", pool.requesters[h].state != nil)
		}
	}
	return str
}
