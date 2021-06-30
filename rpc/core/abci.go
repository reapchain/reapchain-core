package core

import (
	abci "gitlab.reappay.net/sucs-lab//reapchain/abci/types"
	"gitlab.reappay.net/sucs-lab//reapchain/libs/bytes"
	"gitlab.reappay.net/sucs-lab//reapchain/proxy"
	ctypes "gitlab.reappay.net/sucs-lab//reapchain/rpc/core/types"
	rpctypes "gitlab.reappay.net/sucs-lab//reapchain/rpc/jsonrpc/types"
)

// ABCIQuery queries the application for some information.
// More: https://docs.reapchain.com/master/rpc/#/ABCI/abci_query
func ABCIQuery(
	ctx *rpctypes.Context,
	path string,
	data bytes.HexBytes,
	height int64,
	prove bool,
) (*ctypes.ResultABCIQuery, error) {
	resQuery, err := env.ProxyAppQuery.QuerySync(abci.RequestQuery{
		Path:   path,
		Data:   data,
		Height: height,
		Prove:  prove,
	})
	if err != nil {
		return nil, err
	}
	env.Logger.Info("ABCIQuery", "path", path, "data", data, "result", resQuery)
	return &ctypes.ResultABCIQuery{Response: *resQuery}, nil
}

// ABCIInfo gets some info about the application.
// More: https://docs.reapchain.com/master/rpc/#/ABCI/abci_info
func ABCIInfo(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error) {
	resInfo, err := env.ProxyAppQuery.InfoSync(proxy.RequestInfo)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultABCIInfo{Response: *resInfo}, nil
}
