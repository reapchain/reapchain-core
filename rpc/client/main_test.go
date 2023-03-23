package client_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/reapchain/reapchain-core/abci/example/kvstore"
	nm "github.com/reapchain/reapchain-core/node"
	rpctest "github.com/reapchain/reapchain-core/rpc/test"
)

var node *nm.Node

func TestMain(m *testing.M) {
	// start a reapchain-core node (and kvstore) in the background to test against
	dir, err := ioutil.TempDir("/tmp", "rpc-client-test")
	if err != nil {
		panic(err)
	}

	app := kvstore.NewPersistentKVStoreApplication(dir)
	node = rpctest.StartReapchainCore(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopReapchainCore(node)
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
