package light_test

import (
	"context"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"testing"
	"time"

	dbm "github.com/tendermint/tm-db"

	"gitlab.reappay.net/sucs-lab//reapchain/abci/example/kvstore"
	"gitlab.reappay.net/sucs-lab//reapchain/libs/log"
	"gitlab.reappay.net/sucs-lab//reapchain/light"
	"gitlab.reappay.net/sucs-lab//reapchain/light/provider"
	httpp "gitlab.reappay.net/sucs-lab//reapchain/light/provider/http"
	dbs "gitlab.reappay.net/sucs-lab//reapchain/light/store/db"
	rpctest "gitlab.reappay.net/sucs-lab//reapchain/rpc/test"
)

// Automatically getting new headers and verifying them.
func ExampleClient_Update() {
	// give Reapchain time to generate some blocks
	time.Sleep(5 * time.Second)

	dbDir, err := ioutil.TempDir("", "light-client-example")
	if err != nil {
		stdlog.Fatal(err)
	}
	defer os.RemoveAll(dbDir)

	var (
		config  = rpctest.GetConfig()
		chainID = config.ChainID()
	)

	primary, err := httpp.New(chainID, config.RPC.ListenAddress)
	if err != nil {
		stdlog.Fatal(err)
	}

	block, err := primary.LightBlock(context.Background(), 2)
	if err != nil {
		stdlog.Fatal(err)
	}

	db, err := dbm.NewGoLevelDB("light-client-db", dbDir)
	if err != nil {
		stdlog.Fatal(err)
	}

	c, err := light.NewClient(
		context.Background(),
		chainID,
		light.TrustOptions{
			Period: 504 * time.Hour, // 21 days
			Height: 2,
			Hash:   block.Hash(),
		},
		primary,
		[]provider.Provider{primary}, // NOTE: primary should not be used here
		dbs.New(db, chainID),
		light.Logger(log.TestingLogger()),
	)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer func() {
		if err := c.Cleanup(); err != nil {
			stdlog.Fatal(err)
		}
	}()

	time.Sleep(2 * time.Second)

	h, err := c.Update(context.Background(), time.Now())
	if err != nil {
		stdlog.Fatal(err)
	}

	if h != nil && h.Height > 2 {
		fmt.Println("successful update")
	} else {
		fmt.Println("update failed")
	}
	// Output: successful update
}

// Manually getting light blocks and verifying them.
func ExampleClient_VerifyLightBlockAtHeight() {
	// give Reapchain time to generate some blocks
	time.Sleep(5 * time.Second)

	dbDir, err := ioutil.TempDir("", "light-client-example")
	if err != nil {
		stdlog.Fatal(err)
	}
	defer os.RemoveAll(dbDir)

	var (
		config  = rpctest.GetConfig()
		chainID = config.ChainID()
	)

	primary, err := httpp.New(chainID, config.RPC.ListenAddress)
	if err != nil {
		stdlog.Fatal(err)
	}

	block, err := primary.LightBlock(context.Background(), 2)
	if err != nil {
		stdlog.Fatal(err)
	}

	db, err := dbm.NewGoLevelDB("light-client-db", dbDir)
	if err != nil {
		stdlog.Fatal(err)
	}

	c, err := light.NewClient(
		context.Background(),
		chainID,
		light.TrustOptions{
			Period: 504 * time.Hour, // 21 days
			Height: 2,
			Hash:   block.Hash(),
		},
		primary,
		[]provider.Provider{primary}, // NOTE: primary should not be used here
		dbs.New(db, chainID),
		light.Logger(log.TestingLogger()),
	)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer func() {
		if err := c.Cleanup(); err != nil {
			stdlog.Fatal(err)
		}
	}()

	_, err = c.VerifyLightBlockAtHeight(context.Background(), 3, time.Now())
	if err != nil {
		stdlog.Fatal(err)
	}

	h, err := c.TrustedLightBlock(3)
	if err != nil {
		stdlog.Fatal(err)
	}

	fmt.Println("got header", h.Height)
	// Output: got header 3
}

func TestMain(m *testing.M) {
	// start a reapchain node (and kvstore) in the background to test against
	app := kvstore.NewApplication()
	node := rpctest.StartReapchain(app, rpctest.SuppressStdout)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopReapchain(node)
	os.Exit(code)
}
