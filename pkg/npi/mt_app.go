package npi

import (
	"errors"
)

const error_pdu_length_not_enough = "packet not enough"

// 确认数据域是否有效,长度是否合法
func (this *Npdu) checkDataValid(expectLength uint16) error {
	if this.Data == nil || len(this.Data) < int(expectLength) {
		return errors.New(error_pdu_length_not_enough)
	}

	return nil
}

// return true: success  false : failed 一般回复包解析,
func Common_ResponseParse(pdu *Npdu) (bool, error) {
	if err := pdu.checkDataValid(1); err != nil {
		return false, err
	}

	if pdu.Data[0] == 0 {
		return true, nil
	} else {
		return false, nil
	}
}
