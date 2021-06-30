package client_test

import (
	"io/ioutil"
	"os"
	"testing"

	"gitlab.reappay.net/sucs-lab//reapchain/abci/example/kvstore"
	nm "gitlab.reappay.net/sucs-lab//reapchain/node"
	rpctest "gitlab.reappay.net/sucs-lab//reapchain/rpc/test"
)

var node *nm.Node

func TestMain(m *testing.M) {
	// start a reapchain node (and kvstore) in the background to test against
	dir, err := ioutil.TempDir("/tmp", "rpc-client-test")
	if err != nil {
		panic(err)
	}

	app := kvstore.NewPersistentKVStoreApplication(dir)
	node = rpctest.StartReapchain(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopReapchain(node)
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
