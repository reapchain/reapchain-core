package types

import (
	abci "github.com/reapchain/reapchain-core/abci/types"
)


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
func GetVrfCheckList(vrfSet *VrfSet) abci.VrfCheckList {

	vrfCheckList := make([]*abci.VrfCheck, vrfSet.Size())
	for i, vrf := range vrfSet.Vrfs {
		vrfCheckList[i] = NewVrfCheck(vrf)
	}

	return abci.VrfCheckList {
		VrfCheckList: vrfCheckList,
	}
}
