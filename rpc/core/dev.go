package core

import (
	ctypes "gitlab.reappay.net/reapchain/reapchain-core/rpc/core/types"
	rpctypes "gitlab.reappay.net/reapchain/reapchain-core/rpc/jsonrpc/types"
)

// UnsafeFlushMempool removes all transactions from the mempool.
func UnsafeFlushMempool(ctx *rpctypes.Context) (*ctypes.ResultUnsafeFlushMempool, error) {
	env.Mempool.Flush()
	return &ctypes.ResultUnsafeFlushMempool{}, nil
}
