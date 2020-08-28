package ltl

func (this *Ltl_t) SendReadReqBasic(DstAddr uint16, val interface{}) error {
	return this.SendReadReq(DstAddr, TrunkID_GeneralBasic, NodeNumReserved,
		[]uint16{0, 1, 2, 3, 4, 5, 6, 7}, val)
}

func (this *Ltl_t) SendSpecificCmdBasic(DstAddr uint16, cmd byte) error {
	return this.SendSpecificCmd(DstAddr, TrunkID_GeneralBasic, NodeNumReserved,
		RESPONSETYPE_NO, LTL_FRAMECTL_CLIENT_SERVER_DIR, cmd, nil, nil)
}
