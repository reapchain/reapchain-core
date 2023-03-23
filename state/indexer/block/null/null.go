package null

import (
	"context"
	"errors"

	"github.com/reapchain/reapchain-core/libs/pubsub/query"
	"github.com/reapchain/reapchain-core/state/indexer"
	"github.com/reapchain/reapchain-core/types"
)

var _ indexer.BlockIndexer = (*BlockerIndexer)(nil)

// TxIndex implements a no-op block indexer.
type BlockerIndexer struct{}

func (idx *BlockerIndexer) Has(height int64) (bool, error) {
	return false, errors.New(`indexing is disabled (set 'tx_index = "kv"' in config)`)
}

func (idx *BlockerIndexer) Index(types.EventDataNewBlockHeader) error {
	return nil
}

func (idx *BlockerIndexer) Search(ctx context.Context, q *query.Query) ([]int64, error) {
	return []int64{}, nil
}
