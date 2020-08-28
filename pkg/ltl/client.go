package ltl

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type IncomingMsgPkt struct {
	IsBroadCast bool
	SrcAddr     uint16
	ApduData    []byte
}

/*********************帧格式****************************/
// |trunkID | nodeNO | seqnum | frame control | command |

var (
	ErrInvalidDataLength = errors.New("Invalid data length")
	ErrInvalidDataType   = errors.New("Invalid data type")
)

type ProcessIn interface {
	ProInSpecificCmd(srcAddr uint16, hdr *FrameHdr, cmdFormart []byte, val interface{}) byte
	ProInReadRspCmd(srcAddr uint16, hdr *FrameHdr, rdRspStatus []ReadRspStatus, val interface{}) error
	ProInWriteRspCmd(srcAddr uint16, hdr *FrameHdr, wrStatus []WriteRspStatus, val interface{}) error
	ProInWriteRpCfgRspCmd(srcAddr uint16, hdr *FrameHdr, crStatus []WriteRpCfgRspStatus, val interface{}) error
	ProInReadRpCfgRspCmd(srcAddr uint16, hdr *FrameHdr, rcStatus []ReadRpCfgRspStatus, val interface{}) error
	ProInDefaultRsp(srcAddr uint16, hdr *FrameHdr, dfStatus *DefaultRsp, val interface{}) error
	ProInReportCmd(srcAddr uint16, hdr *FrameHdr, rRec []ReportRec) error
}

type WriteCloseMsgComming interface {
	WriteMsg(DstAddr uint16, Data []byte) error
	IncommingMsg() <-chan *IncomingMsgPkt
}

// 内嵌接口,需求一个指定实例,外层结构体中，可以调用内层接口定义的函数
type Ltl_t struct {
	WriteCloseMsgComming
	tsmb sync.Map // 地址对id使用表的映射
	c    *cache.Cache
}

//, defaultExpiration, cleanupInterval time.Duration
func NewClient(wcm WriteCloseMsgComming) *Ltl_t {
	l := &Ltl_t{
		WriteCloseMsgComming: wcm,
		c:                    cache.New(5*time.Second, 2*time.Minute),
	}
	l.c.OnEvicted(expireCb)
	return l
}

/****************************请求*******************************/
const (
	RESPONSETYPE_NO      = 0 // 无需应答
	RESPONSETYPE_DEFAULT = 1 // 默认应答
	RESPONSETYPE_SELF    = 2 // 命令应答
)

// 发送集下特殊命令
func (this *Ltl_t) SendSpecificCmd(DstAddr, trunkId uint16, nodeNO, dir,
	rspType byte, cmd byte, cmdFormart []byte, val interface{}) error {
	var err error
	var seqNum byte

	//无需应答
	if rspType == RESPONSETYPE_NO {
		return this.SendCommand(DstAddr, trunkId, nodeNO, TsmIDReseverd,
			LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC, dir, true, cmd, cmdFormart)
	}

	disableDefaultRsp := false
	if rspType > RESPONSETYPE_DEFAULT {
		disableDefaultRsp = true
	}

	if dir == LTL_FRAMECTL_CLIENT_SERVER_DIR {
		seqNum, err = this.acquireID(DstAddr)
		if err != nil {
			return err
		}
	}

	err = this.SendCommand(DstAddr, trunkId, nodeNO, seqNum,
		LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC, dir, disableDefaultRsp, cmd, cmdFormart)
	if err != nil {
		if dir == LTL_FRAMECTL_CLIENT_SERVER_DIR {
			this.releaseID(DstAddr, seqNum)
		}
		return err
	}
	if dir == LTL_FRAMECTL_CLIENT_SERVER_DIR {
		if !disableDefaultRsp {
			cmd = LTL_CMD_DEFAULT_RSP
		}
		this.hang(DstAddr, seqNum, cmd, val)
	}
	return nil
}

// 读属性请求
func (this *Ltl_t) SendReadReq(DstAddr, trunkId uint16, nodeNO byte, AttrID []uint16, val interface{}) error {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.LittleEndian, AttrID); err != nil {
		return err
	}

	seqNum, err := this.acquireID(DstAddr)
	if err != nil {
		return err
	}

	err = this.SendProfileCmd(DstAddr, trunkId, nodeNO, seqNum, LTL_CMD_READ_ATTRIBUTES, buf.Bytes())
	if err != nil {
		this.releaseID(DstAddr, seqNum)
		return err
	}
	this.hang(DstAddr, seqNum, LTL_CMD_READ_ATTRIBUTES_RSP, val)
	return nil
}

// 写属性记录
type WriteRec struct {
	AttrID    uint16
	AttriData interface{}
}

// 底层,写属性请求
// 属性状态记录不包括写属性成功的应答,如果所有属性都是写成功的,那只记录一个写成功, 否则只包括写失败的
func (this *Ltl_t) sendWriteRequest(DstAddr, trunkId uint16, nodeNO,
	cmd byte, WritewrRec []WriteRec, val interface{}) error {
	var buf, bs []byte
	var dataType byte
	var err error

	for _, v := range WritewrRec {
		bs, dataType, err = marshal(v.AttriData)
		if err != nil {
			return ErrInvalidDataType
		}
		buf = append(buf, byte(v.AttrID), byte(v.AttrID>>8))
		buf = append(buf, dataType)
		buf = append(buf, bs...)
	}

	if cmd == LTL_CMD_WRITE_ATTRIBUTES_NORSP {
		return this.SendProfileCmd(DstAddr, trunkId, nodeNO, TsmIDReseverd, cmd, buf)
	}

	seqNum, err := this.acquireID(DstAddr)
	if err != nil {
		return err
	}
	err = this.SendProfileCmd(DstAddr, trunkId, nodeNO, seqNum, cmd, buf)
	if err != nil {
		this.releaseID(DstAddr, seqNum)
		return err
	}
	this.hang(DstAddr, seqNum, LTL_CMD_WRITE_ATTRIBUTES_RSP, val)
	return nil
}

// 写属性请求
func (this *Ltl_t) SendWriteReq(DstAddr, trunkId uint16, nodeNO byte, WritewrRec []WriteRec, val interface{}) error {
	return this.sendWriteRequest(DstAddr, trunkId, nodeNO, LTL_CMD_WRITE_ATTRIBUTES, WritewrRec, val)
}

// 写属性完整请求,失败将发生回滚
func (this *Ltl_t) SendWriteReqUndivided(DstAddr, trunkId uint16, nodeNO byte, WritewrRec []WriteRec, val interface{}) error {
	return this.sendWriteRequest(DstAddr, trunkId, nodeNO, LTL_CMD_WRITE_ATTRIBUTES_UNDIVIDED, WritewrRec, val)
}

// 写属性不必应答命令请求
func (this *Ltl_t) SendWriteReqNoRsp(DstAddr, trunkId uint16, nodeNO byte, WritewrRec []WriteRec) error {
	return this.sendWriteRequest(DstAddr, trunkId, nodeNO, LTL_CMD_WRITE_ATTRIBUTES_NORSP, WritewrRec, nil)
}

// 写配置报告记录
type WriteReportCfgRec struct {
	AttrID           uint16
	MinReportInt     uint16
	ReportableChange interface{}
}

// 写报告配置请求
/*最小报告间隔:  以秒为单位
  数字量: 值为 0x0000 将忽略报告间隔, 只要值发生变化就进行报告
       值为 0xffff,那么属性将不报告,净被移除设备报告列表之外
        其它值,按间隔时间报告,当发生值变化报告时,重新计时
  模拟量: 值为 0x0000 将忽略报告间隔, 报告由灵敏值决定
         值为0xffff,那么 此属性将不报告,将被移除设备报告列表之外.*/
func (this *Ltl_t) SendWriteReportCfgReq(DstAddr, trunkId uint16, nodeNO byte, wrReportCfgRec []WriteReportCfgRec, val interface{}) error {
	var buf, bs []byte
	var dataType byte
	var err error

	for _, cfgwrRec := range wrReportCfgRec {
		bs, dataType, err = marshal(cfgwrRec.ReportableChange)
		if err != nil || !isBaseDataType(dataType) {
			return ErrInvalidDataType
		}

		buf = append(buf, byte(cfgwrRec.AttrID), byte(cfgwrRec.AttrID>>8))
		buf = append(buf, dataType)
		buf = append(buf, byte(cfgwrRec.MinReportInt), byte(cfgwrRec.MinReportInt>>8))
		if isAnalogDataType(dataType) {
			buf = append(buf, bs...)
		}
	}

	seqNum, err := this.acquireID(DstAddr)
	if err != nil {
		return err
	}

	err = this.SendProfileCmd(DstAddr, trunkId, nodeNO, seqNum, LTL_CMD_CONFIGURE_REPORTING, buf)
	if err != nil {
		this.releaseID(DstAddr, seqNum)
		return err
	}
	this.hang(DstAddr, seqNum, LTL_CMD_CONFIGURE_REPORTING_RSP, val)
	return nil
}

// 读报告配置请求
func (this *Ltl_t) SendReadReportCfgReq(DstAddr, trunkId uint16, nodeNO byte, attrid []uint16, val interface{}) error {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.LittleEndian, attrid); err != nil {
		return err
	}

	seqNum, err := this.acquireID(DstAddr)
	if err != nil {
		return err
	}

	err = this.SendProfileCmd(DstAddr, trunkId, nodeNO, seqNum, LTL_CMD_READ_CONFIGURE_REPORTING, buf.Bytes())
	if err != nil {
		this.releaseID(DstAddr, seqNum)
		return err
	}
	this.hang(DstAddr, seqNum, LTL_CMD_READ_CONFIGURE_REPORTING_RSP, val)
	return nil
}

// 默认回复
type DefaultRsp struct {
	CommandID  byte
	StatusCode byte
}

// 默认应答命令
func (this *Ltl_t) SendDefaultRspCmd(DstAddr, trunkId uint16, nodeNO, seqNum byte, defaultRsp DefaultRsp) error {
	return this.SendProfileCmd(DstAddr, trunkId, nodeNO, seqNum, LTL_CMD_DEFAULT_RSP, []byte{defaultRsp.CommandID, defaultRsp.StatusCode})
}

/****************************解析*******************************/

// 属性值
type AttrValues struct {
	DataType byte
	Data     []byte
}

// 读属性回复状态
type ReadRspStatus struct {
	AttrID uint16
	Status byte
	AttrValues
}

// 解析读属性回复
func parseInReadRspCmd(data []byte) ([]ReadRspStatus, error) {
	var err error
	var datalen int
	var roomLen byte
	var dataType byte
	var tmp ReadRspStatus

	rdRspStatus := []ReadRspStatus{}
	inData := bytes.NewBuffer(data)
	for inData.Len() >= 3 {
		tmp.AttrID = binary.LittleEndian.Uint16(inData.Next(2))
		tmp.Status, _ = inData.ReadByte()

		if tmp.Status != LTL_STATUS_SUCCESS {
			rdRspStatus = append(rdRspStatus, tmp)
			continue
		}

		dataType, err = inData.ReadByte()
		if err != nil {
			return nil, ErrInvalidDataLength
		}

		if !isValidDataType(dataType) {
			return nil, ErrInvalidDataType
		}

		if isComplexDataType(dataType) {
			if roomLen, err = inData.ReadByte(); err != nil {
				return nil, ErrInvalidDataLength
			}
			datalen = int(roomLen)
		} else {
			datalen = getBaseDataTypeLength(dataType)
		}

		if inData.Len() < datalen {
			return nil, ErrInvalidDataLength
		}
		tmp.Data = inData.Next(datalen)
		tmp.DataType = dataType
		rdRspStatus = append(rdRspStatus, tmp)
	}

	return rdRspStatus, nil
}

// 写属性回复状态
type WriteRspStatus struct {
	Status byte   // should be LTL_STATUS_SUCCESS or error
	AttrID uint16 // attribute ID
}

// 解析写回复命令
func parseInWriteRspCmd(data []byte) ([]WriteRspStatus, error) {
	var err error
	var tmp WriteRspStatus

	if len(data) == 1 && data[0] == LTL_SUCCESS {
		return []WriteRspStatus{{Status: data[0]}}, nil
	}

	wrStatus := []WriteRspStatus{}
	inData := bytes.NewReader(data)
	for inData.Len() >= 3 {
		if err = binary.Read(inData, binary.LittleEndian, &tmp); err != nil {
			return nil, err
		}

		wrStatus = append(wrStatus, tmp)
	}

	return wrStatus, nil
}

// 写报告配置回复状态
type WriteRpCfgRspStatus struct {
	Status byte
	AttrID uint16
}

// 解析报告配置回复命令
func parseInWriteRpCfgRspCmd(data []byte) ([]WriteRpCfgRspStatus, error) {
	var err error
	var tmp WriteRpCfgRspStatus

	if len(data) == 1 && data[0] == LTL_SUCCESS {
		return []WriteRpCfgRspStatus{{Status: data[0]}}, nil
	}

	crStatus := []WriteRpCfgRspStatus{}
	inData := bytes.NewBuffer(data)
	for inData.Len() >= 3 {
		if err = binary.Read(inData, binary.LittleEndian, &tmp); err != nil {
			return nil, err
		}

		crStatus = append(crStatus, tmp)
	}

	return crStatus, nil
}

// 读报告配置回复状态
type ReadRpCfgRspStatus struct {
	Status       byte
	AttrID       uint16
	MinReportInt uint16
	AttrValues
}

// 解析读报告配置回复命令
func parseInReadRpCfgRspCmd(data []byte) ([]ReadRpCfgRspStatus, error) {
	var err error
	var dataLen int
	var dataType byte
	var bs []byte
	var tmp ReadRpCfgRspStatus

	inData := bytes.NewBuffer(data)
	rcStatus := []ReadRpCfgRspStatus{}
	for inData.Len() >= 3 {
		tmp.Status, _ = inData.ReadByte()
		tmp.AttrID = binary.LittleEndian.Uint16(inData.Next(2))
		if tmp.Status != LTL_STATUS_SUCCESS {
			rcStatus = append(rcStatus, tmp)
			continue
		}

		if dataType, err = inData.ReadByte(); err != nil {
			return nil, ErrInvalidDataLength
		}

		if !isBaseDataType(dataType) {
			return nil, ErrInvalidDataType
		}

		if bs, err = inData.ReadBytes(2); err != nil {
			return nil, ErrInvalidDataLength
		}
		tmp.MinReportInt = binary.LittleEndian.Uint16(bs)

		if isAnalogDataType(dataType) {
			if dataLen = getBaseDataTypeLength(dataType); inData.Len() < dataLen {
				return nil, ErrInvalidDataLength
			}

			tmp.Data = inData.Next(dataLen)
		}
		tmp.DataType = dataType
		rcStatus = append(rcStatus, tmp)
	}

	return rcStatus, nil
}

// 报告记录
type ReportRec struct {
	AttrID uint16
	AttrValues
}

func parseInReportCmd(data []byte) ([]ReportRec, error) {
	var dataLen int
	var roomLen byte
	var dataType byte
	var tmp ReportRec
	var err error

	inData := bytes.NewBuffer(data)
	rRec := []ReportRec{}
	for inData.Len() >= 3 {
		tmp.AttrID = binary.LittleEndian.Uint16(inData.Next(2))
		dataType, _ = inData.ReadByte()
		if !isValidDataType(dataType) {
			return nil, ErrInvalidDataType
		}

		if isComplexDataType(dataType) {
			if roomLen, err = inData.ReadByte(); err != nil {
				return nil, ErrInvalidDataLength
			}
			dataLen = int(roomLen)
		} else {
			dataLen = getBaseDataTypeLength(dataType)
		}

		if inData.Len() < dataLen {
			return nil, ErrInvalidDataLength
		}

		tmp.Data = inData.Next(dataLen)
		tmp.DataType = dataType
		rRec = append(rRec, tmp)
	}

	return rRec, nil
}

// 默认回复
func parseInDefaultRspCmd(data []byte) (*DefaultRsp, error) {
	if len(data) < 2 {
		return nil, ErrInvalidDataLength
	}

	return &DefaultRsp{data[0], data[1]}, nil
}

func (this *Ltl_t) ServerInApdu(ctx context.Context, pi ProcessIn) {
	var pkt *IncomingMsgPkt
	var status byte
	var errs error

	//logs.Info("ltl: serverInApdu routine started")
	for {
		select {
		case <-ctx.Done():
			//logs.Warn("ltl: serverInApdu routine closed!")
			return
		case pkt = <-this.IncommingMsg():
		}
		if len(pkt.ApduData) < hdrSize() {
			continue
		}
		hdr, remainData := parseHdr(pkt.ApduData)
		//logs.Debug("ServerInApdu: hdr- %#v, data- %#v", hdr, remainData)

		// 是否为标准命令
		if IsProfileCmd(hdr.Type_l) {
			if hdr.CommandID < LTL_CMD_PROFILE_MAX &&
				(hdr.CommandID == LTL_CMD_READ_ATTRIBUTES_RSP ||
					hdr.CommandID == LTL_CMD_WRITE_ATTRIBUTES_RSP ||
					hdr.CommandID == LTL_CMD_CONFIGURE_REPORTING_RSP ||
					hdr.CommandID == LTL_CMD_READ_CONFIGURE_REPORTING_RSP ||
					hdr.CommandID == LTL_CMD_DEFAULT_RSP) {
				if hdr.TransSeqNum == TsmIDReseverd {
					continue
				}
				cmdId, val, err := this.FindItem(pkt.SrcAddr, hdr.TransSeqNum)
				if err != nil || cmdId != hdr.CommandID {
					continue
				}
				//logs.Debug("find the item:ID is %d", hdr.TransSeqNum)

				switch hdr.CommandID {
				case LTL_CMD_READ_ATTRIBUTES_RSP:
					rdRspStatus, err := parseInReadRspCmd(remainData)
					if err != nil {
						errs = err
						break
					}
					errs = pi.ProInReadRspCmd(pkt.SrcAddr, hdr, rdRspStatus, val)
				case LTL_CMD_WRITE_ATTRIBUTES_RSP:
					wrStatus, err := parseInWriteRspCmd(remainData)
					if err != nil {
						errs = err
						break
					}
					errs = pi.ProInWriteRspCmd(pkt.SrcAddr, hdr, wrStatus, val)
				case LTL_CMD_CONFIGURE_REPORTING_RSP:
					crStatus, err := parseInWriteRpCfgRspCmd(remainData)
					if err != nil {
						errs = err
						break
					}
					errs = pi.ProInWriteRpCfgRspCmd(pkt.SrcAddr, hdr, crStatus, val)
				case LTL_CMD_READ_CONFIGURE_REPORTING_RSP:
					rcStatus, err := parseInReadRpCfgRspCmd(remainData)
					if err != nil {
						errs = err
						break
					}
					errs = pi.ProInReadRpCfgRspCmd(pkt.SrcAddr, hdr, rcStatus, val)

				case LTL_CMD_DEFAULT_RSP:
					dfStatus, err := parseInDefaultRspCmd(remainData)
					if err != nil {
						errs = err
						break
					}
					errs = pi.ProInDefaultRsp(pkt.SrcAddr, hdr, dfStatus, val)
				default:
					errs = errors.New("no support command")
				}
				if errs != nil {
					//logs.Error("ServerInApdu: type - profile ,commandID - %d, %s", hdr.CommandID, errs)
				}
				continue
			} else if hdr.CommandID == LTL_CMD_REPORT_ATTRIBUTES {
				rwrRec, err := parseInReportCmd(remainData)
				if err != nil {
					//logs.Error("ServerInApdu: type - profile ,commandID - %d, %s", hdr.CommandID, err)
					continue
				}
				err = pi.ProInReportCmd(pkt.SrcAddr, hdr, rwrRec)
				if err != nil {
					//logs.Error("ServerInApdu: type - profile ,commandID - %d, %s", hdr.CommandID, err)
					continue
				}
			} else {
				status = LTL_STATUS_UNSUP_GENERAL_COMMAND
			}
		} else if IsTrunkCmd(hdr.Type_l) {
			if IsFromServer(hdr.Dir) {
				if hdr.TransSeqNum == TsmIDReseverd {
					continue
				}
				cmdId, val, err := this.FindItem(pkt.SrcAddr, hdr.TransSeqNum)
				if err != nil || cmdId != hdr.CommandID {
					continue
				}
				// The return value of the plugin function will be
				//  LTL_STATUS_SUCCESS - Supported and need default response
				//  LTL_STATUS_FAILURE - Unsupported
				//  LTL_STATUS_CMD_HAS_RSP - Supported and do not need default rsp
				status = pi.ProInSpecificCmd(pkt.SrcAddr, hdr, remainData, val)
			}
		}

		//非广播消息可以进行默认应答
		if (!pkt.IsBroadCast) &&
			(hdr.DisableDefaultRsp == LTL_FRAMECTL_DIS_DEFAULT_RSP_OFF) &&
			!IsFromServer(hdr.FrameHdrCtl.Dir) &&
			hdr.TransSeqNum > 0 {
			this.SendDefaultRspCmd(pkt.SrcAddr, hdr.TrunkID, hdr.NodeNo, hdr.TransSeqNum,
				DefaultRsp{CommandID: hdr.CommandID, StatusCode: status})
		}
	}
}
