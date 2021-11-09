package v0

// import (
// 	"errors"
// 	"fmt"
// 	"sync/atomic"
// 	"time"

// 	"github.com/reapchain/reapchain-core/libs/service"
// 	"github.com/reapchain/reapchain-core/p2p"
// 	sm "github.com/reapchain/reapchain-core/state"
// )

// // The callet will verify the state
// func (pool *BlockPool) PeekTwoStates() (firstState *sm.State, secondState *sm.State) {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	if r := pool.stateRequesters[pool.stateHeight]; r != nil {
// 		firstState = r.getState()
// 	}
// 	if r := pool.stateRequesters[pool.stateHeight+1]; r != nil {
// 		secondState = r.getState()
// 	}
// 	return
// }

// // AddBlock validates that the block comes from the peer it was expected from and calls the requester to store it.
// // TODO: ensure that blocks come in order for each peer.
// func (pool *BlockPool) AddState(statePeerID p2p.ID, state *sm.State, blockSize int) {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	// fmt.Println("AddState-", state)
// 	requester := pool.stateRequesters[state.LastBlockHeight]
// 	if requester == nil {
// 		diff := pool.stateHeight - state.LastBlockHeight
// 		if diff < 0 {
// 			diff *= -1
// 		}
// 		if diff > maxDiffBetweenCurrentAndReceivedBlockHeight {
// 			pool.sendError(errors.New("peer sent us a state we didn't expect with a height too far ahead/behind"), statePeerID)
// 		}
// 		return
// 	}

// 	if requester.setState(state, statePeerID) {
// 		atomic.AddInt32(&pool.numStatePending, -1)
// 		peer := pool.statePeers[statePeerID]
// 		if peer != nil {
// 			peer.decrStatePending(blockSize)
// 		}
// 	} else {
// 		fmt.Println("stompesi - state", state)
// 		pool.Logger.Info("invalid state peer", "peer", statePeerID, "blockHeight", state.LastBlockHeight)
// 		pool.sendError(errors.New("invalid state peer"), statePeerID)
// 	}
// }

// func (pool *BlockPool) PopStateRequest() {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	if r := pool.stateRequesters[pool.stateHeight]; r != nil {
// 		if err := r.Stop(); err != nil {
// 			pool.Logger.Error("Error stopping requester", "err", err)
// 		}
// 		delete(pool.stateRequesters, pool.stateHeight)
// 		pool.stateHeight++
// 	} else {
// 		panic(fmt.Sprintf("Expected state requester to pop, got nothing at stateHeight %v", pool.stateHeight))
// 	}
// }

// func (pool *BlockPool) RedoStateRequest(height int64) p2p.ID {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	request := pool.stateRequesters[height]
// 	statePeerID := request.getStatePeerID()
// 	if statePeerID != p2p.ID("") {
// 		pool.removeStatePeer(statePeerID)
// 	}
// 	return statePeerID
// }

// func (pool *BlockPool) RemoveStatePeer(statePeerID p2p.ID) {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	pool.removeStatePeer(statePeerID)
// }

// func (pool *BlockPool) removeStatePeer(statePeerID p2p.ID) {
// 	for _, requester := range pool.stateRequesters {
// 		if requester.getStatePeerID() == statePeerID {
// 			requester.stateRedo(statePeerID)
// 		}
// 	}

// 	peer, ok := pool.statePeers[statePeerID]
// 	if ok {
// 		if peer.timeout != nil {
// 			peer.timeout.Stop()
// 		}

// 		delete(pool.statePeers, statePeerID)

// 		if peer.height == pool.maxStatePeerHeight {
// 			pool.updateMaxStatePeerHeight()
// 		}
// 	}
// }

// func (pool *BlockPool) removeStateTimedoutPeers() {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	for _, peer := range pool.statePeers {
// 		if !peer.didTimeout && peer.numPending > 0 {
// 			curRate := peer.recvMonitor.Status().CurRate
// 			// curRate can be 0 on start
// 			if curRate != 0 && curRate < minRecvRate {
// 				err := errors.New("peer is not sending us data fast enough")
// 				pool.sendError(err, peer.id)
// 				pool.Logger.Error("SendTimeout", "peer", peer.id,
// 					"reason", err,
// 					"curRate", fmt.Sprintf("%d KB/s", curRate/1024),
// 					"minRate", fmt.Sprintf("%d KB/s", minRecvRate/1024))
// 				peer.didTimeout = true
// 			}
// 		}
// 		if peer.didTimeout {
// 			pool.removeStatePeer(peer.id)
// 		}
// 	}
// }

// func (pool *BlockPool) MaxStatePeerHeight() int64 {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()
// 	return pool.maxStatePeerHeight
// }

// func (pool *BlockPool) SetStatePeerRange(statePeerID p2p.ID, base int64, height int64) {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	peer := pool.statePeers[statePeerID]
// 	if peer != nil {
// 		peer.base = base
// 		peer.height = height
// 	} else {
// 		peer = newBPPeer(pool, statePeerID, base, height)
// 		peer.setLogger(pool.Logger.With("peer", statePeerID))
// 		pool.statePeers[statePeerID] = peer
// 	}

// 	if height > pool.maxStatePeerHeight {
// 		pool.maxStatePeerHeight = height
// 	}
// }

// func (pool *BlockPool) updateMaxStatePeerHeight() {
// 	var max int64
// 	for _, peer := range pool.statePeers {
// 		if peer.height > max {
// 			max = peer.height
// 		}
// 	}
// 	pool.maxStatePeerHeight = max
// }

// func (pool *BlockPool) pickIncrAvailableStatePeer(height int64) *bpPeer {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	for _, peer := range pool.statePeers {
// 		if peer.didTimeout {
// 			pool.removePeer(peer.id)
// 			continue
// 		}
// 		if peer.numPending >= maxPendingRequestsPerPeer {
// 			continue
// 		}
// 		if height < peer.base || height > peer.height {
// 			continue
// 		}
// 		peer.incrStatePending()
// 		return peer
// 	}
// 	return nil
// }

// func (pool *BlockPool) makeNextStateRequester() {
// 	pool.mtx.Lock()
// 	defer pool.mtx.Unlock()

// 	nextHeight := pool.stateHeight + pool.stateRequestersLen()
// 	if nextHeight > pool.maxStatePeerHeight {
// 		return
// 	}

// 	request := newStateRequester(pool, nextHeight)

// 	pool.stateRequesters[nextHeight] = request
// 	atomic.AddInt32(&pool.numStatePending, 1)

// 	err := request.Start()
// 	if err != nil {
// 		request.Logger.Error("Error starting request", "err", err)
// 	}
// }

// func (pool *BlockPool) stateRequestersLen() int64 {
// 	return int64(len(pool.stateRequesters))
// }

// func (pool *BlockPool) sendStateRequest(height int64, statePeerID p2p.ID) {
// 	if !pool.IsRunning() {
// 		return
// 	}
// 	pool.stateRequestsCh <- StateRequest{height, statePeerID}
// }

// func (peer *bpPeer) incrStatePending() {
// 	if peer.numPending == 0 {
// 		peer.resetMonitor()
// 		peer.resetTimeout()
// 	}
// 	peer.numPending++
// }

// func (peer *bpPeer) decrStatePending(recvSize int) {
// 	peer.numPending--
// 	if peer.numPending == 0 {
// 		peer.timeout.Stop()
// 	} else {
// 		peer.recvMonitor.Update(recvSize)
// 		peer.resetTimeout()
// 	}
// }

// func newStateRequester(pool *BlockPool, height int64) *bpRequester {
// 	bpr := &bpRequester{
// 		pool:       pool,
// 		height:     height,
// 		gotStateCh: make(chan struct{}, 1),
// 		redoCh:     make(chan p2p.ID, 1),

// 		statePeerID: "",
// 		state:       nil,
// 	}
// 	bpr.BaseService = *service.NewBaseService(nil, "stateRequester", bpr)
// 	return bpr
// }

// func (bpr *bpRequester) resetStates() {
// 	bpr.mtx.Lock()
// 	defer bpr.mtx.Unlock()

// 	if bpr.block != nil {
// 		atomic.AddInt32(&bpr.pool.numStatePending, 1)
// 	}

// 	bpr.statePeerID = ""
// 	bpr.state = nil
// }

// func (bpr *bpRequester) setState(state *sm.State, statePeerID p2p.ID) bool {
// 	bpr.mtx.Lock()

// 	if bpr.state != nil || bpr.statePeerID != statePeerID {
// 		fmt.Println("stompesi - setState", "bpr.statePeerID", bpr.statePeerID, "statePeerID")
// 		bpr.mtx.Unlock()
// 		return false
// 	}
// 	bpr.state = state
// 	bpr.mtx.Unlock()

// 	select {
// 	case bpr.gotStateCh <- struct{}{}:
// 	default:
// 	}
// 	return true
// }

// func (bpr *bpRequester) getState() *sm.State {
// 	bpr.mtx.Lock()
// 	defer bpr.mtx.Unlock()
// 	return bpr.state
// }

// func (bpr *bpRequester) getStatePeerID() p2p.ID {
// 	bpr.mtx.Lock()
// 	defer bpr.mtx.Unlock()
// 	return bpr.statePeerID
// }

// func (bpr *bpRequester) resetState() {
// 	bpr.mtx.Lock()
// 	defer bpr.mtx.Unlock()

// 	if bpr.state != nil {
// 		atomic.AddInt32(&bpr.pool.numStatePending, 1)
// 	}

// 	bpr.statePeerID = ""
// 	bpr.state = nil
// }

// func (bpr *bpRequester) stateRedo(statePeerID p2p.ID) {
// 	select {
// 	case bpr.redoCh <- statePeerID:
// 	default:
// 	}
// }

// func (bpr *bpRequester) requestStateRoutine() {
// OUTER_LOOP:
// 	for {
// 		// Pick a peer to send request to.
// 		var peer *bpPeer
// 	PICK_PEER_LOOP:
// 		for {
// 			if !bpr.IsRunning() || !bpr.pool.IsRunning() {
// 				return
// 			}
// 			peer = bpr.pool.pickIncrAvailableStatePeer(bpr.height)
// 			if peer == nil {
// 				time.Sleep(requestIntervalMS * time.Millisecond)
// 				continue PICK_PEER_LOOP
// 			}
// 			break PICK_PEER_LOOP
// 		}
// 		bpr.mtx.Lock()
// 		bpr.statePeerID = peer.id
// 		bpr.mtx.Unlock()

// 		// Send request and wait.
// 		bpr.pool.sendStateRequest(bpr.height, peer.id)
// 	WAIT_LOOP:
// 		for {
// 			select {
// 			case <-bpr.pool.Quit():
// 				if err := bpr.Stop(); err != nil {
// 					bpr.Logger.Error("Error stopped requester", "err", err)
// 				}
// 				return
// 			case <-bpr.Quit():
// 				return
// 			case statePeerID := <-bpr.redoCh:
// 				if statePeerID == bpr.statePeerID {
// 					bpr.resetState()
// 					continue OUTER_LOOP
// 				} else {
// 					continue WAIT_LOOP
// 				}
// 			case <-bpr.gotStateCh:
// 				// We got a block!
// 				// Continue the for-loop and wait til Quit.
// 				continue WAIT_LOOP
// 			}
// 		}
// 	}
// }

// type StateRequest struct {
// 	Height int64
// 	PeerID p2p.ID
// }
