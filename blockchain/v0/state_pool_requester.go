package v0

import (
	"sync/atomic"
	"time"

	"github.com/reapchain/reapchain-core/libs/service"
	tmsync "github.com/reapchain/reapchain-core/libs/sync"
	"github.com/reapchain/reapchain-core/p2p"
	sm "github.com/reapchain/reapchain-core/state"
)

type StatePoolRequester struct {
	service.BaseService
	statePool  *StatePool
	height     int64
	gotStateCh chan struct{}
	redoCh     chan p2p.ID

	mtx    tmsync.Mutex
	peerID p2p.ID

	state *sm.State
}

func NewStatePoolRequester(statePool *StatePool, height int64) *StatePoolRequester {
	stateRequester := &StatePoolRequester{
		statePool:  statePool,
		height:     height,
		gotStateCh: make(chan struct{}, 1),
		redoCh:     make(chan p2p.ID, 1),

		peerID: "",
		state:  nil,
	}
	stateRequester.BaseService = *service.NewBaseService(nil, "stateRequester", stateRequester)
	return stateRequester
}

func (stateRequester *StatePoolRequester) OnStart() error {
	go stateRequester.requestRoutine()
	return nil
}

func (stateRequester *StatePoolRequester) setState(state *sm.State, peerID p2p.ID) bool {
	stateRequester.mtx.Lock()
	if stateRequester.state != nil || stateRequester.peerID != peerID {
		stateRequester.mtx.Unlock()
		return false
	}
	stateRequester.state = state
	stateRequester.mtx.Unlock()

	select {
	case stateRequester.gotStateCh <- struct{}{}:
	default:
	}
	return true
}

func (stateRequester *StatePoolRequester) getState() *sm.State {
	stateRequester.mtx.Lock()
	defer stateRequester.mtx.Unlock()
	return stateRequester.state
}

func (stateRequester *StatePoolRequester) getPeerID() p2p.ID {
	stateRequester.mtx.Lock()
	defer stateRequester.mtx.Unlock()
	return stateRequester.peerID
}

func (stateRequester *StatePoolRequester) reset() {
	stateRequester.mtx.Lock()
	defer stateRequester.mtx.Unlock()

	if stateRequester.state != nil {
		atomic.AddInt32(&stateRequester.statePool.numPending, 1)
	}

	stateRequester.peerID = ""
	stateRequester.state = nil
}

func (stateRequester *StatePoolRequester) redo(peerID p2p.ID) {
	select {
	case stateRequester.redoCh <- peerID:
	default:
	}
}

func (stateRequester *StatePoolRequester) requestRoutine() {
OUTER_LOOP:
	for {
		// Pick a peer to send request to.
		var peer *StatePoolPeer
	PICK_PEER_LOOP:
		for {
			if !stateRequester.IsRunning() || !stateRequester.statePool.IsRunning() {
				return
			}
			peer = stateRequester.statePool.pickIncrAvailablePeer(stateRequester.height)
			if peer == nil {
				time.Sleep(requestIntervalMS * time.Millisecond)
				continue PICK_PEER_LOOP
			}
			break PICK_PEER_LOOP
		}
		stateRequester.mtx.Lock()
		stateRequester.peerID = peer.id
		stateRequester.mtx.Unlock()

		// Send request and wait.
		stateRequester.statePool.sendRequest(stateRequester.height, peer.id)
	WAIT_LOOP:
		for {
			select {
			case <-stateRequester.statePool.Quit():
				if err := stateRequester.Stop(); err != nil {
					stateRequester.Logger.Error("Error stopped requester", "err", err)
				}
				return
			case <-stateRequester.Quit():
				return
			case peerID := <-stateRequester.redoCh:
				if peerID == stateRequester.peerID {
					stateRequester.reset()
					continue OUTER_LOOP
				} else {
					continue WAIT_LOOP
				}
			case <-stateRequester.gotStateCh:
				// We got a block!
				// Continue the for-loop and wait til Quit.
				continue WAIT_LOOP
			}
		}
	}
}
