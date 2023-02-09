package types

import (
	abci "github.com/reapchain/reapchain-core/abci/types"
)

// It communicates with SDK.
// It represents whether the steering member candidate sent a VRF in the consensus round (VRF set).
func NewVrfCheck(vrf *Vrf) *abci.VrfCheck {
	if vrf == nil {
		return nil
	}

	vrfCheck := &abci.VrfCheck{
		SteeringMemberCandidateAddress: vrf.SteeringMemberCandidatePubKey.Address(),
		IsVrfTransmission: true,
	}

	if vrf.Value == nil {
		vrfCheck.IsVrfTransmission = false	
	}

	return vrfCheck
}

// Generate checklist for each steering member candidate whether the steering member candidate sent a VRF in the consensus round (VRF set).
func GetVrfCheckList(vrfSet *VrfSet) abci.VrfCheckList {

	vrfCheckList := make([]*abci.VrfCheck, vrfSet.Size())
	for i, vrf := range vrfSet.Vrfs {
		vrfCheckList[i] = NewVrfCheck(vrf)
	}

	return abci.VrfCheckList {
		VrfCheckList: vrfCheckList,
	}
}
