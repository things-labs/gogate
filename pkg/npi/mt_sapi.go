package npi

/********** Commands *************/
// 3.7.1.1
func (this *Monitor) Zb_System_Reset() error {
	return this.SendAsynchData(MT_SAPI_SYS_RESET, nil)
}

// 3.7.1.3
func (this *Monitor) Zb_PermitJoingReq(DestAddr uint16, Timeout byte) (bool, error) {
	rspPdu, err := this.SendSynchData(MT_SAPI_PMT_JOIN_REQ, []byte{byte(DestAddr), byte(DestAddr >> 8), Timeout})
	if err != nil {
		return false, err
	}
	return Common_ResponseParse(rspPdu)
}
