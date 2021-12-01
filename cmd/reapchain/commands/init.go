package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	cfg "github.com/reapchain/reapchain-core/config"
	tmos "github.com/reapchain/reapchain-core/libs/os"
	tmrand "github.com/reapchain/reapchain-core/libs/rand"
	"github.com/reapchain/reapchain-core/p2p"
	"github.com/reapchain/reapchain-core/privval"
	"github.com/reapchain/reapchain-core/types"
	tmtime "github.com/reapchain/reapchain-core/types/time"
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
	privValidator := privval.LoadOrGenFilePV(privValKeyFile, privValStateFile)

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

		genDoc.StandingMembers = []types.GenesisMember{{
			Address: pubKey.Address(),
			PubKey:  pubKey,
		}}

		genDoc.ConsensusRound = types.NewConsensusRound(1, 0, 0, 0)

		qrnValue := tmrand.Uint64()
		qrn := types.NewQrn(0, pubKey, qrnValue)
		qrn.Timestamp = genDoc.GenesisTime

		err = privValidator.SignQrn(qrn)
		if err != nil {
			logger.Error("Can't sign qrn", "err", err)
		}

		if qrn.VerifySign() == false {
			logger.Error("Is invalid sign of qrn")
		}

		genDoc.Qrns = []types.Qrn{*qrn}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}

		genDoc.SteeringMemberCandidates = []types.GenesisMember{}

		genDoc.Vrfs = []types.Vrf{}

		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}