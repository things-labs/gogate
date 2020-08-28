package npi

import (
	"bytes"
	"encoding/binary"
)

/*********************************MT_ZDO**************************************************/
// 3.12.1.1
func (this *Monitor) Zdo_NwkAddrReq(IeeeAddr uint64, reqType, StatIndex byte) (bool, error) {
	data := make([]byte, 6)

	binary.LittleEndian.PutUint64(data, IeeeAddr)
	data[4] = reqType
	data[5] = StatIndex
	rspPdu, err := this.SendSynchData(MT_ZDO_NWK_ADDR_REQ, data)
	if err != nil {
		return false, nil
	}

	return Common_ResponseParse(rspPdu)
}

// 3.12.1.2
func (this *Monitor) Zdo_IeeeAddrReq(ShortAddr uint16, reqType, StatIndex byte) (bool, error) {
	rspPdu, err := this.SendSynchData(MT_ZDO_NWK_ADDR_REQ,
		[]byte{byte(ShortAddr), byte(ShortAddr >> 8), reqType, StatIndex})
	if err != nil {
		return false, nil
	}

	return Common_ResponseParse(rspPdu)
}

// 3.12.1.20
func (this *Monitor) Zdo_MgmtLeaveReq(DstAddr uint16, DeviceAddr uint64, RemoveChild_Rejoin byte) (bool, error) {
	data := make([]byte, 11)
	binary.LittleEndian.PutUint16(data, DstAddr)
	binary.LittleEndian.PutUint64(data[2:], DeviceAddr)
	data[10] = RemoveChild_Rejoin

	rspPdu, err := this.SendSynchData(MT_ZDO_MGMT_LEAVE_REQ, data)
	if err != nil {
		return false, err
	}

	return Common_ResponseParse(rspPdu)
}

type ZdoMgmtPermitJionRsp_t struct {
	SrcAddr uint16
	Status  byte
}

/********** Callback *************/
// 3.12.2.21
func Zdo_MgmtPermitJoinRspParse(pdu *Npdu) (*ZdoMgmtPermitJionRsp_t, error) {
	var err error

	if err = pdu.checkDataValid(3); err != nil {
		return nil, err
	}

	rsp := &ZdoMgmtPermitJionRsp_t{}
	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// 3.12.2.22
func Zdo_StateChangeIndParse(pdu *Npdu) (byte, error) {
	var err error

	if err = pdu.checkDataValid(1); err != nil {
		return 0, err
	}

	return pdu.Data[0], nil
}

type ZdoEnddeviceAnnceInd_t struct {
	SrcAddr, NwkAddr uint16
	IeeeAddr         uint64
	Capabilities     byte
}

// 3.12.2.23
func Zdo_EnddeviceAnnceIndParse(pdu *Npdu) (*ZdoEnddeviceAnnceInd_t, error) {
	var err error

	if err = pdu.checkDataValid(13); err != nil {
		return nil, err
	}
	annce := &ZdoEnddeviceAnnceInd_t{}
	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, annce); err != nil {
		return nil, err
	}

	return annce, nil
}

type ZdoLeaveInd_t struct {
	SrcAddr uint16
	ExtAddr uint64
	Request bool //Boolean, TRUE = request, FALSE = indication.
	Remove  bool //Boolean, TRUE = remove children
	Rejoin  bool //Boolean, TRUE = rejoin
}

// 3.12.2.30
func Zdo_LeaveIndParse(pdu *Npdu) (*ZdoLeaveInd_t, error) {
	var err error

	if err = pdu.checkDataValid(13); err != nil {
		return nil, err
	}

	leave := &ZdoLeaveInd_t{}
	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, leave); err != nil {
		return nil, err
	}

	return leave, nil
}
