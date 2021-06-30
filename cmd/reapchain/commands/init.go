package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	cfg "gitlab.reappay.net/sucs-lab/reapchain/config"
	tmos "gitlab.reappay.net/sucs-lab/reapchain/libs/os"
	tmrand "gitlab.reappay.net/sucs-lab/reapchain/libs/rand"
	"gitlab.reappay.net/sucs-lab/reapchain/p2p"
	"gitlab.reappay.net/sucs-lab/reapchain/privval"
	"gitlab.reappay.net/sucs-lab/reapchain/types"
	tmtime "gitlab.reappay.net/sucs-lab/reapchain/types/time"
)

// InitFilesCmd initialises a fresh Reapchain Core instance.
var InitFilesCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Reapchain",
	RunE:  initFiles,
}

func initFiles(cmd *cobra.Command, args []string) error {
	return initFilesWithConfig(config)
}

func initFilesWithConfig(config *cfg.Config) error {
	// private validator
	privValKeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()
	var pv *privval.FilePV
	if tmos.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		logger.Info("Found private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	} else {
		pv = privval.GenFilePV(privValKeyFile, privValStateFile)
		pv.Save()
		logger.Info("Generated private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	}

	nodeKeyFile := config.NodeKeyFile()
	if tmos.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := config.GenesisFile()
	if tmos.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc := types.GenesisDoc{
			ChainID:         fmt.Sprintf("test-chain-%v", tmrand.Str(6)),
			GenesisTime:     tmtime.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}
		pubKey, err := pv.GetPubKey()
		if err != nil {
			return fmt.Errorf("can't get pubkey: %w", err)
		}
		genDoc.Validators = []types.GenesisValidator{{
			Address: pubKey.Address(),
			PubKey:  pubKey,
			Power:   10,
		}}

		// Default 상임위 의원 추가
		genDoc.StandingMembers = []types.GenesisStandingMember{{
			Address: pubKey.Address(),
			PubKey:  pubKey,
		}}

		// Default 양자난수 추가
		genDoc.Qns = []types.Qn{{
			Address: pubKey.Address(),
			PubKey:  pubKey,
			Value:   tmrand.Uint64(),
		}}

		genDoc.ConsensusRoundInfo = types.NewConsensusRound(0, 0)

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}
