package npi

import (
	"bytes"
	"encoding/binary"
)

/*********************************MT_SYS**************************************************/
/********** command *************/
//3.8.1.1   Type:  0: soft reset , 1 : hard reset
func (this *Monitor) Sys_ResetReq(Type byte) error {
	return this.SendAsynchData(MT_SYS_RESET_REQ, []byte{Type})
}

// 3.8.1.2 return mt command Capabilities
func (this *Monitor) Sys_Ping() (uint16, error) {
	rspPdu, err := this.SendSynchData(MT_SYS_PING, nil)
	if err != nil {
		return 0, err
	}

	if err = rspPdu.checkDataValid(2); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(rspPdu.Data), nil
}

type SysVersion_t struct {
	TransportRev, ProductId, MajorRel, MinorRel, MaintRel byte
}

func (this *Monitor) Sys_Version() (*SysVersion_t, error) {
	rspPdu, err := this.SendSynchData(MT_SYS_VERSION, nil)
	if err != nil {
		return nil, err
	}

	if err = rspPdu.checkDataValid(5); err != nil {
		return nil, err
	}

	return &SysVersion_t{
		TransportRev: rspPdu.Data[0],
		ProductId:    rspPdu.Data[1],
		MajorRel:     rspPdu.Data[2],
		MinorRel:     rspPdu.Data[3],
		MaintRel:     rspPdu.Data[4],
	}, nil
}

type SysResetInd_t struct {
	Reason, TransportRev, ProductId, MajorRel, MinorRel, HwRev byte
}

/********** Callback *************/
// 3.8.2.1
func Sys_ResetIndParse(pdu *Npdu) (*SysResetInd_t, error) {
	var err error

	if err = pdu.checkDataValid(6); err != nil {
		return nil, err
	}

	srt := &SysResetInd_t{}
	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, srt); err != nil {
		return nil, err
	}

	return srt, nil
}
