package state

import (
	"errors"
	"fmt"

)

// Rollback overwrites the current ReapchainCore state (height n) with the most
// recent previous state (height n - 1).
// Note that this function does not affect application state.
func Rollback(bs BlockStore, ss Store) (int64, []byte, error) {
	rollbackState, err := ss.LoadRollback()
	if err != nil {
		return -1, nil, err
	}
	if rollbackState.IsEmpty() {
		return -1, nil, errors.New("no state found")
	}

	// persist the new state. This overrides the invalid one. NOTE: this will also
	// persist the validator set and consensus params over the existing structures,
	// but both should be the same
	if err := ss.Save(rollbackState); err != nil {
		return -1, nil, fmt.Errorf("failed to save rolled back state: %w", err)
	}

	bs.SaveRollbackBlock(rollbackState.LastBlockHeight)
	// bs.saveState()

	// SaveBlockStoreState(&bss, bs.db)


	return rollbackState.LastBlockHeight, rollbackState.AppHash, nil
}
