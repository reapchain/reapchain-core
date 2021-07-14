package consensus

import (
	"fmt"

	"github.com/reapchain/reapchain-core/types"
)

func (cs *State) tryUpdateQrn(qrn *types.Qrn) (bool, error) {
	fmt.Println("한종빈1")
	updated, _ := cs.updateQrn(qrn)

	return updated, nil
}

func (cs *State) updateQrn(qrn *types.Qrn) (updated bool, err error) {
	fmt.Println("한종빈2", cs.state.QrnSet.Height, qrn.Height, cs.Height)
	if qrn.Height+1 != cs.Height {
		return
	}

	//TODO: stompesi
	fmt.Println("한종빈3", cs.state.QrnSet.StandingMemberSet)
	updated, err = cs.state.QrnSet.UpdateQrn(qrn)
	if !updated {
		return
	}

	return updated, err
}
