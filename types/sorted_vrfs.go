package types

import "bytes"

// 정렬하기 위한 구조체
type SortedVrfs []*Vrf

func (sms SortedVrfs) Len() int { return len(sms) }

func (sms SortedVrfs) Less(i, j int) bool {
	return bytes.Compare(sms[i].Seed, sms[j].Seed) == 1
}

func (sms SortedVrfs) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}
