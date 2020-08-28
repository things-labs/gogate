package npi

import (
	"bytes"
	"encoding/binary"
)

/*********************************MT_AF**************************************************/
/********** Commands *************/
//3.2.1.2
func (this *Monitor) Af_DataReq(dstAddr, ClusterID uint16, dstEndpoint, srcEndPoint, TransID, Options, Radius byte, data []byte) (bool, error) {
	var s = []interface{}{
		dstAddr,
		dstEndpoint,
		srcEndPoint,
		ClusterID,
		TransID,
		Options,
		Radius,
		byte(len(data)),
	}

	buf := &bytes.Buffer{}
	for _, v := range s {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			return false, err
		}
	}
	odat := make([]byte, 0, buf.Len()+len(data))
	odat = append(odat, buf.Bytes()...)
	odat = append(odat, data...)

	rspPdu, err := this.SendSynchData(MT_AF_DATA_REQUEST, odat)
	if err != nil {
		return false, err
	}

	return Common_ResponseParse(rspPdu)
}

/********** Callbacks *************/
type AfDataConfirm_t struct {
	Status, Endpoint, TransId byte
}

// 3.2.1.1
func Af_DataConfirmParse(pdu *Npdu) (*AfDataConfirm_t, error) {
	var err error

	if err := pdu.checkDataValid(3); err != nil {
		return nil, err
	}

	cfm := &AfDataConfirm_t{}

	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, cfm); err != nil {
		return nil, err
	}

	return cfm, nil
}

type AfIncomingMsgHead_t struct {
	Groupid, ClusterID, SrcAddr                               uint16
	SrcEndPoint, DstEndpoint, WasBroadcast, Rssi, SecurityUse byte
	Timestamp                                                 uint32
	TransID, Length                                           byte
}
type AfIncomingMsg_t struct {
	AfIncomingMsgHead_t
	Data []byte
}

// 3.2.1.1
func Af_IncomingMsgParse(pdu *Npdu) (*AfIncomingMsg_t, error) {
	var err error
	var head AfIncomingMsgHead_t

	if err = pdu.checkDataValid(17); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(pdu.Data)
	if err = binary.Read(buf, binary.LittleEndian, &head); err != nil {
		return nil, err
	}

	return &AfIncomingMsg_t{
		AfIncomingMsgHead_t: head,
		Data:                buf.Bytes()[:head.Length], // 只要需求内的数据,其它数据丢弃
	}, nil
}
