package v0

import (
	"fmt"
	"reflect"
	"time"

	bc "github.com/reapchain/reapchain-core/blockchain"
	"github.com/reapchain/reapchain-core/libs/log"
	"github.com/reapchain/reapchain-core/p2p"
	bcproto "github.com/reapchain/reapchain-core/proto/reapchain/blockchain"
	sm "github.com/reapchain/reapchain-core/state"
	"github.com/reapchain/reapchain-core/store"
	"github.com/reapchain/reapchain-core/types"
)

const (
	// BlockchainChannel is a channel for blocks and status updates (`BlockStore` height)
	BlockchainChannel = byte(0x40)

	trySyncIntervalMS = 10

	// stop syncing when last block's time is
	// within this much of the system time.
	// stopSyncingDurationMinutes = 10

	// ask for best height every 10s
	statusUpdateIntervalSeconds = 10
	// check if we should switch to consensus reactor
	switchToConsensusIntervalSeconds = 1
)

type consensusReactor interface {
	// for when we switch from blockchain reactor and fast sync to
	// the consensus machine
	SwitchToConsensus(state sm.State, skipWAL bool)
}

// BlockchainReactor handles long-term catchup syncing.
type BlockchainReactor struct {
	p2p.BaseReactor

	// immutable
	initialState sm.State

	blockExec  *sm.BlockExecutor
	store      *store.BlockStore
	stateStore sm.Store
	blockPool  *BlockPool
	statePool  *StatePool
	fastSync   bool

	requestsCh      <-chan Request
	stateRequestsCh <-chan Request
	errorsCh        <-chan peerError
}

// NewBlockchainReactor returns new reactor instance.
func NewBlockchainReactor(state sm.State, blockExec *sm.BlockExecutor, store *store.BlockStore,
	stateStore sm.Store, fastSync bool) *BlockchainReactor {

	if state.LastBlockHeight != store.Height() {
		panic(fmt.Sprintf("state (%v) and store (%v) height mismatch", state.LastBlockHeight,
			store.Height()))
	}

	requestsCh := make(chan Request, maxTotalRequesters)
	stateRequestsCh := make(chan Request, maxTotalRequesters)

	const capacity = 1000                      // must be bigger than peers count
	errorsCh := make(chan peerError, capacity) // so we don't block in #Receive#pool.AddBlock

	startHeight := store.Height() + 1
	if startHeight == 1 {
		startHeight = state.InitialHeight
	}
	blockPool := NewBlockPool(startHeight, requestsCh, errorsCh)
	statePool := NewStatePool(startHeight, stateRequestsCh, errorsCh)

	bcR := &BlockchainReactor{
		initialState:    state,
		blockExec:       blockExec,
		store:           store,
		stateStore:      stateStore,
		blockPool:       blockPool,
		statePool:       statePool,
		fastSync:        fastSync,
		requestsCh:      requestsCh,
		stateRequestsCh: stateRequestsCh,
		errorsCh:        errorsCh,
	}

	bcR.BaseReactor = *p2p.NewBaseReactor("BlockchainReactor", bcR)
	return bcR
}

// SetLogger implements service.Service by setting the logger on reactor and pool.
func (bcR *BlockchainReactor) SetLogger(l log.Logger) {
	bcR.BaseService.Logger = l
	bcR.blockPool.Logger = l
	bcR.statePool.Logger = l
}

// OnStart implements service.Service.
func (bcR *BlockchainReactor) OnStart() error {
	if bcR.fastSync {
		err := bcR.blockPool.Start()
		if err != nil {
			return err
		}
		err = bcR.statePool.Start()
		if err != nil {
			return err
		}

		go bcR.poolRoutine(false)
	}
	return nil
}

// SwitchToFastSync is called by the state sync reactor when switching to fast sync.
func (bcR *BlockchainReactor) SwitchToFastSync(state sm.State) error {
	bcR.fastSync = true
	bcR.initialState = state

	bcR.blockPool.height = state.LastBlockHeight + 1
	bcR.statePool.height = state.LastBlockHeight + 1
	err := bcR.blockPool.Start()
	if err != nil {
		return err
	}
	err = bcR.statePool.Start()
	if err != nil {
		return err
	}

	go bcR.poolRoutine(true)
	return nil
}

// OnStop implements service.Service.
func (bcR *BlockchainReactor) OnStop() {
	if bcR.fastSync {
		if err := bcR.blockPool.Stop(); err != nil {
			bcR.Logger.Error("Error stopping pool", "err", err)
		}
		if err := bcR.statePool.Stop(); err != nil {
			bcR.Logger.Error("Error stopping pool", "err", err)
		}
	}
}

// GetChannels implements Reactor
func (bcR *BlockchainReactor) GetChannels() []*p2p.ChannelDescriptor {
	return []*p2p.ChannelDescriptor{
		{
			ID:                  BlockchainChannel,
			Priority:            5,
			SendQueueCapacity:   1000,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: bc.MaxMsgSize,
		},
	}
}

// AddPeer implements Reactor by sending our state to peer.
func (bcR *BlockchainReactor) AddPeer(peer p2p.Peer) {
	msgBytes, err := bc.EncodeMsg(&bcproto.StatusResponse{
		Base:   bcR.store.Base(),
		Height: bcR.store.Height()})
	if err != nil {
		bcR.Logger.Error("could not convert msg to protobuf", "err", err)
		return
	}

	peer.Send(BlockchainChannel, msgBytes)
	// it's OK if send fails. will try later in poolRoutine

	// peer is added to the pool once we receive the first
	// bcStatusResponseMessage from the peer and call pool.SetPeerRange
}

// RemovePeer implements Reactor by removing peer from the pool.
func (bcR *BlockchainReactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	bcR.blockPool.RemovePeer(peer.ID())
	bcR.statePool.RemovePeer(peer.ID())
}

// respondToPeer loads a block and sends it to the requesting peer,
// if we have it. Otherwise, we'll respond saying we don't have it.
func (bcR *BlockchainReactor) respondToPeer(msg *bcproto.BlockRequest,
	src p2p.Peer) (queued bool) {

	fmt.Println("respondToPeer - ", msg.Height)

	block := bcR.store.LoadBlock(msg.Height)
	if block != nil {
		bl, err := block.ToProto()
		if err != nil {
			bcR.Logger.Error("could not convert msg to protobuf", "err", err)
			return false
		}

		msgBytes, err := bc.EncodeMsg(&bcproto.BlockResponse{Block: bl})
		if err != nil {
			bcR.Logger.Error("could not marshal msg", "err", err)
			return false
		}

		return src.TrySend(BlockchainChannel, msgBytes)
	}

	bcR.Logger.Info("Peer asking for a block we don't have", "src", src, "height", msg.Height)

	msgBytes, err := bc.EncodeMsg(&bcproto.NoBlockResponse{Height: msg.Height})
	if err != nil {
		bcR.Logger.Error("could not convert msg to protobuf", "err", err)
		return false
	}

	return src.TrySend(BlockchainChannel, msgBytes)
}

// Receive implements Reactor by handling 4 types of messages (look below).
func (bcR *BlockchainReactor) Receive(chID byte, src p2p.Peer, msgBytes []byte) {
	msg, err := bc.DecodeMsg(msgBytes)
	if err != nil {
		bcR.Logger.Error("Error decoding message", "src", src, "chId", chID, "err", err)
		bcR.Switch.StopPeerForError(src, err)
		return
	}

	if err = bc.ValidateMsg(msg); err != nil {
		bcR.Logger.Error("Peer sent us invalid msg3", "peer", src, "msg", msg, "err", err)
		bcR.Switch.StopPeerForError(src, err)
		return
	}

	// bcR.Logger.Error("Receive", "src", src, "chID", chID, "msg", msg)

	switch msg := msg.(type) {
	case *bcproto.BlockRequest:
		bcR.respondToPeer(msg, src)

	case *bcproto.BlockResponse:
		bi, err := types.BlockFromProto(msg.Block)
		// fmt.Println("bcproto.BlockResponse", bi.LastCommit.Height)

		if err != nil {
			bcR.Logger.Error("Block content is invalid", "err", err)
			return
		}
		bcR.blockPool.AddBlock(src.ID(), bi, len(msgBytes))

	case *bcproto.StateRequest:
		bcR.respondStateToPeer(msg, src)

	case *bcproto.StateResponse:
		state, err := sm.SyncStateFromProto(msg.State)
		fmt.Println("stompesi-state", state.LastBlockHeight)

		if err != nil {
			bcR.Logger.Error("State content is invalid", "err", err)
			return
		}
		bcR.statePool.AddState(src.ID(), state, len(msgBytes))

	case *bcproto.StatusRequest:
		// Send peer our state.
		msgBytes, err := bc.EncodeMsg(&bcproto.StatusResponse{
			Height: bcR.store.Height(),
			Base:   bcR.store.Base(),
		})
		if err != nil {
			bcR.Logger.Error("could not convert msg to protobut", "err", err)
			return
		}
		src.TrySend(BlockchainChannel, msgBytes)
	case *bcproto.StatusResponse:
		// Got a peer status. Unverified.
		bcR.blockPool.SetPeerRange(src.ID(), msg.Base, msg.Height)
		bcR.statePool.SetPeerRange(src.ID(), msg.Base, msg.Height)
	case *bcproto.NoBlockResponse:
		bcR.Logger.Debug("Peer does not have requested block", "peer", src, "height", msg.Height)
	default:
		bcR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
	}
}

// Handle messages from the poolReactor telling the reactor what to do.
// NOTE: Don't sleep in the FOR_LOOP or otherwise slow it down!
func (bcR *BlockchainReactor) poolRoutine(stateSynced bool) {

	trySyncTicker := time.NewTicker(trySyncIntervalMS * time.Millisecond)
	defer trySyncTicker.Stop()

	statusUpdateTicker := time.NewTicker(statusUpdateIntervalSeconds * time.Second)
	defer statusUpdateTicker.Stop()

	switchToConsensusTicker := time.NewTicker(switchToConsensusIntervalSeconds * time.Second)
	defer switchToConsensusTicker.Stop()

	blocksSynced := uint64(0)

	chainID := bcR.initialState.ChainID
	state := bcR.initialState

	// debug.PrintStack()

	lastHundred := time.Now()
	lastRate := 0.0

	didProcessCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-bcR.Quit():
				return
			// stompesi
			case <-bcR.blockPool.Quit():
				return
			case request := <-bcR.requestsCh:
				peer := bcR.Switch.Peers().Get(request.PeerID)
				if peer == nil {
					continue
				}
				msgBytes, err := bc.EncodeMsg(&bcproto.BlockRequest{Height: request.Height})
				if err != nil {
					bcR.Logger.Error("could not convert msg to proto", "err", err)
					continue
				}

				queued := peer.TrySend(BlockchainChannel, msgBytes)
				if !queued {
					bcR.Logger.Debug("Send queue is full, drop block request", "peer", peer.ID(), "height", request.Height)
				}

			case request := <-bcR.stateRequestsCh:
				peer := bcR.Switch.Peers().Get(request.PeerID)
				if peer == nil {
					continue
				}
				msgBytes, err := bc.EncodeMsg(&bcproto.StateRequest{Height: request.Height})
				if err != nil {
					bcR.Logger.Error("could not convert msg to proto", "err", err)
					continue
				}

				queued := peer.TrySend(BlockchainChannel, msgBytes)
				if !queued {
					bcR.Logger.Debug("Send queue is full, drop block request", "peer", peer.ID(), "height", request.Height)
				}

			case err := <-bcR.errorsCh:
				peer := bcR.Switch.Peers().Get(err.peerID)
				if peer != nil {
					bcR.Switch.StopPeerForError(peer, err)
				}

			case <-statusUpdateTicker.C:
				// ask for status updates
				go bcR.BroadcastStatusRequest() // nolint: errcheck

			}
		}
	}()

FOR_LOOP:
	for {
		select {
		case <-switchToConsensusTicker.C:
			height, numPending, lenRequesters := bcR.blockPool.GetStatus()
			_, numStatePending, lenStateRequesters := bcR.statePool.GetStatus()

			outbound, inbound, _ := bcR.Switch.NumPeers()

			bcR.Logger.Debug("Consensus ticker", "numPending", numPending, "numStatePending", numStatePending, "total", lenRequesters, "lenStateRequesters", lenStateRequesters,
				"outbound", outbound, "inbound", inbound)

			if bcR.blockPool.IsCaughtUp() {
				bcR.Logger.Info("Time to switch to consensus reactor!", "height", height)
				if err := bcR.blockPool.Stop(); err != nil {
					bcR.Logger.Error("Error stopping pool", "err", err)
				}
				conR, ok := bcR.Switch.Reactor("CONSENSUS").(consensusReactor)
				if ok {
					conR.SwitchToConsensus(state, blocksSynced > 0 || stateSynced)
				}

				break FOR_LOOP
			}

			if bcR.statePool.IsCaughtUp() {
				bcR.Logger.Info("Time to switch to consensus reactor!", "height", height)
				if err := bcR.statePool.Stop(); err != nil {
					bcR.Logger.Error("Error stopping pool", "err", err)
				}
				conR, ok := bcR.Switch.Reactor("CONSENSUS").(consensusReactor)
				if ok {
					conR.SwitchToConsensus(state, blocksSynced > 0 || stateSynced)
				}

				break FOR_LOOP
			}

		case <-trySyncTicker.C: // chan time
			select {
			case didProcessCh <- struct{}{}:
			default:
			}

		case <-didProcessCh:
			// NOTE: It is a subtle mistake to process more than a single block
			// at a time (e.g. 10) here, because we only TrySend 1 request per
			// loop.  The ratio mismatch can result in starving of blocks, a
			// sudden burst of requests and responses, and repeat.
			// Consequently, it is better to split these routines rather than
			// coupling them as it's written here.  TODO uncouple from request
			// routine.

			// See if there are any blocks to sync.
			first, second := bcR.blockPool.PeekTwoBlocks()
			// bcR.Logger.Info("TrySync peeked", "first", first, "second", second)
			if first == nil || second == nil {
				// We need both to sync the first block.
				didProcessCh <- struct{}{}
				continue FOR_LOOP
			}

			firstState, secondState := bcR.statePool.PeekTwoStates()
			if firstState == nil || secondState == nil {
				didProcessCh <- struct{}{}
				continue FOR_LOOP
			}

			firstParts := first.MakePartSet(types.BlockPartSizeBytes)
			firstPartSetHeader := firstParts.Header()
			firstID := types.BlockID{Hash: first.Hash(), PartSetHeader: firstPartSetHeader}
			// Finally, verify the first block using the second's commit
			// NOTE: we can probably make this more efficient, but note that calling
			// first.Hash() doesn't verify the tx contents, so MakePartSet() is
			// currently necessary.
			err := state.Validators.VerifyCommitLight(chainID, firstID, first.Height, second.LastCommit)
			if err != nil {
				fmt.Println("stompesi-step3")
				bcR.Logger.Error("Error in validation", "err", err)
				peerID := bcR.blockPool.RedoRequest(first.Height)
				peer := bcR.Switch.Peers().Get(peerID)
				if peer != nil {
					// NOTE: we've already removed the peer's request, but we
					// still need to clean up the rest.
					bcR.Switch.StopPeerForError(peer, fmt.Errorf("blockchainReactor validation error: %v", err))
				}
				peerID2 := bcR.blockPool.RedoRequest(second.Height)
				peer2 := bcR.Switch.Peers().Get(peerID2)
				if peer2 != nil && peer2 != peer {
					// NOTE: we've already removed the peer's request, but we
					// still need to clean up the rest.
					bcR.Switch.StopPeerForError(peer2, fmt.Errorf("blockchainReactor validation error: %v", err))
				}
				continue FOR_LOOP
			} else {
				bcR.blockPool.PopRequest()
				bcR.statePool.PopRequest()

				// TODO: validate

				fmt.Println("first.Height", first.Height)

				state.QrnSet = firstState.QrnSet
				state.VrfSet = firstState.VrfSet
				state.SettingSteeringMember = firstState.SettingSteeringMember
				state.NextQrnSet = firstState.NextQrnSet
				state.NextVrfSet = firstState.NextVrfSet

				fmt.Println("state.NextQrnSet", state.NextQrnSet)
				fmt.Println("state.NextVrfSet", state.NextVrfSet)

				// TODO: batch saves so we dont persist to disk every block
				bcR.store.SaveBlock(first, firstParts, second.LastCommit)

				// TODO: same thing for app - but we would need a way to
				// get the hash without persisting the state
				var err error

				if state.SettingSteeringMember != nil {
					fmt.Println("poolRoutine-SteeringMemberIndexes", state.SettingSteeringMember.SteeringMemberIndexes)
				} else {
					fmt.Println("poolRoutine-SteeringMemberIndexes-nil")
				}

				state, _, err = bcR.blockExec.ApplyBlock(state, firstID, first)
				if err != nil {
					// TODO This is bad, are we zombie?
					panic(fmt.Sprintf("Failed to process committed block (%d:%X): %v", first.Height, first.Hash(), err))
				}
				blocksSynced++

				if blocksSynced%100 == 0 {
					lastRate = 0.9*lastRate + 0.1*(100/time.Since(lastHundred).Seconds())
					lastHundred = time.Now()
				}
			}
			continue FOR_LOOP

		case <-bcR.Quit():
			break FOR_LOOP
		}
	}
}

// BroadcastStatusRequest broadcasts `BlockStore` base and height.
func (bcR *BlockchainReactor) BroadcastStatusRequest() error {
	bm, err := bc.EncodeMsg(&bcproto.StatusRequest{})
	if err != nil {
		bcR.Logger.Error("could not convert msg to proto", "err", err)
		return fmt.Errorf("could not convert msg to proto: %w", err)
	}

	bcR.Switch.Broadcast(BlockchainChannel, bm)

	return nil
}

func (bcR *BlockchainReactor) respondStateToPeer(msg *bcproto.StateRequest,
	src p2p.Peer) (queued bool) {
	fmt.Println("respondStateToPeer - ", msg.Height)
	qrnSet, err := bcR.stateStore.LoadQrnSet(msg.Height)
	if err != nil {
		fmt.Println("respondStateToPeer1 - qnr load 실패", msg.Height)
		return false
	}

	if qrnSet != nil {

		vrfSet, err := bcR.stateStore.LoadVrfSet(msg.Height)
		if err != nil {
			fmt.Println("respondStateToPeer2 - ", msg.Height)
			return false
		}

		settingSteeringMember, err := bcR.stateStore.LoadSettingSteeringMember(msg.Height)
		if err != nil {
			fmt.Println("respondStateToPeer3 - ", msg.Height)
			return false
		}

		nextVrfSet, err := bcR.stateStore.LoadNextVrfSet(msg.Height)
		if err != nil {
			fmt.Println("respondStateToPeer2 - ", msg.Height)
			return false
		}

		nextQrnSet, err := bcR.stateStore.LoadNextQrnSet(msg.Height)
		if err != nil {
			fmt.Println("respondStateToPeer2 - ", msg.Height)
			return false
		}

		state := &sm.State{
			LastBlockHeight:       msg.Height,
			SettingSteeringMember: settingSteeringMember,
			VrfSet:                vrfSet,
			QrnSet:                qrnSet,
			NextVrfSet:            nextVrfSet,
			NextQrnSet:            nextQrnSet,
		}

		stateProto, err := state.ToProto()
		if err != nil {
			fmt.Println("respondStateToPeer4 - ", msg.Height)
			return false
		}

		msgBytes, err := bc.EncodeMsg(&bcproto.StateResponse{State: stateProto})

		if err != nil {
			fmt.Println("respondStateToPeer5 - ", msg.Height)
			bcR.Logger.Error("could not marshal msg", "err", err)
			return false
		}

		fmt.Println("respondStateToPeer1 - 성공", msg.Height)

		return src.TrySend(BlockchainChannel, msgBytes)
	}
	fmt.Println("respondStateToPeer1 - qnr 없음", msg.Height)

	msgBytes, err := bc.EncodeMsg(&bcproto.NoStateResponse{Height: msg.Height})
	if err != nil {
		fmt.Println("respondStateToPeer6 - ", msg.Height)
		bcR.Logger.Error("could not convert msg to protobuf", "err", err)
		return false
	}

	fmt.Println("respondStateToPeer1 - 없어", msg.Height)
	return src.TrySend(BlockchainChannel, msgBytes)
}
