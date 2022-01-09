package psql

import (
	"github.com/reapchain/reapchain-core/state/indexer"
	"github.com/reapchain/reapchain-core/state/txindex"
)

var (
	_ indexer.BlockIndexer = BackportBlockIndexer{}
	_ txindex.TxIndexer    = BackportTxIndexer{}
)
