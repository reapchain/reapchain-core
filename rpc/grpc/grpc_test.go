package coregrpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.reappay.net/sucs-lab/reapchain/abci/example/kvstore"
	core_grpc "gitlab.reappay.net/sucs-lab/reapchain/rpc/grpc"
	rpctest "gitlab.reappay.net/sucs-lab/reapchain/rpc/test"
)

func TestMain(m *testing.M) {
	// start a reapchain node in the background to test against
	app := kvstore.NewApplication()
	node := rpctest.StartReapchain(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopReapchain(node)
	os.Exit(code)
}

func TestBroadcastTx(t *testing.T) {
	res, err := rpctest.GetGRPCClient().BroadcastTx(
		context.Background(),
		&core_grpc.RequestBroadcastTx{Tx: []byte("this is a tx")},
	)
	require.NoError(t, err)
	require.EqualValues(t, 0, res.CheckTx.Code)
	require.EqualValues(t, 0, res.DeliverTx.Code)
}
