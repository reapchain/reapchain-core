package core

import (
	ctypes "gitlab.reappay.net/sucs-lab/reapchain/rpc/core/types"
	rpctypes "gitlab.reappay.net/sucs-lab/reapchain/rpc/jsonrpc/types"
)

// Health gets node health. Returns empty result (200 OK) on success, no
// response - in case of an error.
// More: https://docs.reapchain.com/master/rpc/#/Info/health
func Health(ctx *rpctypes.Context) (*ctypes.ResultHealth, error) {
	return &ctypes.ResultHealth{}, nil
}
