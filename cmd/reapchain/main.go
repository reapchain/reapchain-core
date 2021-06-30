package main

import (
	"os"
	"path/filepath"

	cmd "gitlab.reappay.net/sucs-lab//reapchain/cmd/reapchain/commands"
	"gitlab.reappay.net/sucs-lab//reapchain/cmd/reapchain/commands/debug"
	cfg "gitlab.reappay.net/sucs-lab//reapchain/config"
	"gitlab.reappay.net/sucs-lab//reapchain/libs/cli"
	nm "gitlab.reappay.net/sucs-lab//reapchain/node"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.GenValidatorCmd,
		cmd.InitFilesCmd,
		cmd.ProbeUpnpCmd,
		cmd.LightCmd,
		cmd.ReplayCmd,
		cmd.ReplayConsoleCmd,
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ShowValidatorCmd,
		cmd.TestnetFilesCmd,
		cmd.ShowNodeIDCmd,
		cmd.GenNodeKeyCmd,
		cmd.VersionCmd,
		debug.DebugCmd,
		cli.NewCompletionCmd(rootCmd, true),
	)

	// NOTE:
	// Users wishing to:
	//	* Use an external signer for their validators
	//	* Supply an in-proc abci app
	//	* Supply a genesis doc file from another source
	//	* Provide their own DB implementation
	// can copy this file and use something other than the
	// DefaultNewNode function
	//fmt.Println("stompesi-start-0000")
	nodeFunc := nm.DefaultNewNode

	// Create & start node
	rootCmd.AddCommand(cmd.NewRunNodeCmd(nodeFunc))

	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("$HOME", cfg.DefaultReapchainDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
