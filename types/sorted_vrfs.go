package types

import "bytes"

// For sorting vrf with value
type SortedVrfs []*Vrf

func (sms SortedVrfs) Len() int { return len(sms) }

func (sms SortedVrfs) Less(i, j int) bool {
	if sms[j].Proof == nil {
		return true
	}

	if sms[i].Proof == nil {
		return false
	}
	return bytes.Compare(sms[i].Value, sms[j].Value) == 1
}

func (sms SortedVrfs) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}
