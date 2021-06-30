package core

import (
	ctypes "gitlab.reappay.net/sucs-lab/reapchain/rpc/core/types"
	rpctypes "gitlab.reappay.net/sucs-lab/reapchain/rpc/jsonrpc/types"
)

// UnsafeFlushMempool removes all transactions from the mempool.
func UnsafeFlushMempool(ctx *rpctypes.Context) (*ctypes.ResultUnsafeFlushMempool, error) {
	env.Mempool.Flush()
	return &ctypes.ResultUnsafeFlushMempool{}, nil
}
