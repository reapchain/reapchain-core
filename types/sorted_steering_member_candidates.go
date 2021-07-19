package types

import "bytes"

// 정렬하기 위한 구조체
type SortedSteeringMemberCandidates []*SteeringMemberCandidate

func (sms SortedSteeringMemberCandidates) Len() int { return len(sms) }

func (sms SortedSteeringMemberCandidates) Less(i, j int) bool {
	return bytes.Compare(sms[i].Address, sms[j].Address) == -1
}

func (sms SortedSteeringMemberCandidates) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}
