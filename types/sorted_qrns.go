package types

type SortedQrns []*Qrn

func (sms SortedQrns) Len() int { return len(sms) }

func (sms SortedQrns) Less(i, j int) bool {
	if sms[j].Signature == nil {
		return true
	}

	if sms[i].Signature == nil {
		return false
	}
	return sms[i].Value > sms[j].Value
}

func (sms SortedQrns) Swap(i, j int) {
	sms[i], sms[j] = sms[j], sms[i]
}
