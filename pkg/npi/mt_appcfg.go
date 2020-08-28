package npi

import (
	"bytes"
	"encoding/binary"
)

const (
	Cms_mode_Init             = 0x00
	Cms_mode_TouchLink        = 0x01
	Cms_mode_NetworkSteer     = 0x02
	Cms_mode_NetworkFormation = 0x04
	Cms_mode_FindBind         = 0x08
)

/*********************************MT_APP_CONFIG**************************************************/
// 3.13.1.5
func (this *Monitor) Appcfg_BdbStartCommissioningReq(cms_mode byte) (bool, error) {
	rspPdu, err := this.SendSynchData(MT_APP_CNF_BDB_START_COMMISSIONING, []byte{cms_mode})
	if err != nil {
		return false, err
	}

	return Common_ResponseParse(rspPdu)
}

type AppcfgCommissioningNotice_t struct {
	Status             byte
	Commissioning_mode byte
}

/********** Callback *************/
// 3.13.2.1
func Appcfg_BdbCommissioningNoticeParse(pdu *Npdu) (*AppcfgCommissioningNotice_t, error) {
	var err error

	if err = pdu.checkDataValid(2); err != nil {
		return nil, err
	}

	notice := &AppcfgCommissioningNotice_t{}
	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, notice); err != nil {
		return nil, err
	}

	return notice, nil
}
