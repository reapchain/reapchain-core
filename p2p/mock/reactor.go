package mock

import (
	"github.com/reapchain/reapchain-core/libs/log"
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/p2p/conn"
)

type Reactor struct {
	p2p.BaseReactor

	Channels []*conn.ChannelDescriptor
}

func NewReactor() *Reactor {
	r := &Reactor{}
	r.BaseReactor = *p2p.NewBaseReactor("Mock-PEX", r)
	r.SetLogger(log.TestingLogger())
	return r
}

func (r *Reactor) GetChannels() []*conn.ChannelDescriptor            { return r.Channels }
func (r *Reactor) AddPeer(peer p2p.Peer)                             {}
func (r *Reactor) RemovePeer(peer p2p.Peer, reason interface{})      {}
func (r *Reactor) Receive(chID byte, peer p2p.Peer, msgBytes []byte) {}
