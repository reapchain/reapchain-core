package types

import "bytes"

type SortedStandingMembers []*StandingMember

func (sms SortedStandingMembers) Len() int { return len(sms) }

func (sms SortedStandingMembers) Less(i, j int) bool {
	return bytes.Compare(sms[i].Address, sms[j].Address) == -1
}

func (sms SortedStandingMembers) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}
