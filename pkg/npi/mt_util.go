package npi

import (
	"bytes"
	"encoding/binary"
)

/********** Command *************/
type UtilDeviceInfoHead_t struct {
	Status                                   byte
	IeeeAddr                                 uint64
	ShortAddr                                uint16
	DeviceType, DeviceState, NumAssocDevices byte
}
type UtilDeviceInfo_t struct {
	UtilDeviceInfoHead_t
	AssocDeviceList []uint16 // Array of 16-bits of network addresses of Reduce Function Devices associated to the local device
}

// 3.10.1.1
func (this *Monitor) Util_GetDeviceInfoReq() (*UtilDeviceInfo_t, error) {
	var err error
	var head UtilDeviceInfoHead_t
	var AssocDeviceList []uint16
	var rspPdu *Npdu

	rspPdu, err = this.SendSynchData(MT_UTIL_GET_DEVICE_INFO, nil)
	if err != nil {
		return nil, err
	}

	if err = rspPdu.checkDataValid(14); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(rspPdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, &head); err != nil {
		return nil, err
	}

	if head.NumAssocDevices > 0 {
		AssocDeviceList = make([]uint16, head.NumAssocDevices)
		if err = binary.Read(buf, binary.LittleEndian, AssocDeviceList); err != nil {
			return nil, err
		}
	}

	return &UtilDeviceInfo_t{UtilDeviceInfoHead_t: head, AssocDeviceList: AssocDeviceList}, nil
}

/********** Callback *************/
