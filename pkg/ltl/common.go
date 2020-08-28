package ltl

import (
	"bytes"
	"encoding/binary"
)

type FrameHdrCtl struct {
	Type_l, DisableDefaultRsp, Dir byte
}

type FrameHdr struct {
	TrunkID                        uint16
	NodeNo, TransSeqNum, CommandID byte
	FrameHdrCtl
}

// frame control direction 帧控制域 请求来自服务器
func IsFromServer(direction byte) bool {
	return direction == LTL_FRAMECTL_SERVER_CLIENT_DIR
}

// frame control type Profile 帧控制域 框架下命令
func IsProfileCmd(type_l byte) bool {
	return type_l == LTL_FRAMECTL_TYPE_PROFILE
}

// frame control type Trunk Specific 帧控制域 集下特殊命令
func IsTrunkCmd(type_l byte) bool {
	return type_l == LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC
}

// 帧头大小
func hdrSize() int {
	//trunkID + nodeNO + seqnum + frame control + cmdID
	return (2 + 1 + 1 + 1 + 1)
}

// 编码帧头到结构体
func encodeHdr(trunkId uint16, nodeNO, seqNum, frameCtlType, dir byte, disableDefaultRsp bool, cmd byte) *FrameHdr {
	hdrctl := FrameHdrCtl{
		Type_l:            (frameCtlType & LTL_FRAMECTL_TYPE_MASK),
		DisableDefaultRsp: LTL_FRAMECTL_DIS_DEFAULT_RSP_OFF,
		Dir:               dir,
	}

	if disableDefaultRsp {
		hdrctl.DisableDefaultRsp = LTL_FRAMECTL_DIS_DEFAULT_RSP_ON
	}

	return &FrameHdr{
		TrunkID:     trunkId,
		TransSeqNum: seqNum,
		NodeNo:      nodeNO,
		CommandID:   cmd,
		FrameHdrCtl: hdrctl,
	}
}

// Build the Frame Control byte 帧头序列化成字节流
func (this *FrameHdr) buildHdr() []byte {
	fc := this.Type_l
	fc |= this.DisableDefaultRsp << 2
	fc |= this.Dir << 3
	s := make([]byte, 0, hdrSize())
	s = append(s, byte(this.TrunkID), byte(this.TrunkID>>8))

	return append(s, this.NodeNo, this.TransSeqNum, fc, this.CommandID)
}

// 解析出帧头,获取帧头后的数据切片 get hdr and pointer pass header
func parseHdr(data []byte) (*FrameHdr, []byte) {
	hdr := &FrameHdr{}
	inData := bytes.NewBuffer(data)

	hdr.TrunkID = binary.LittleEndian.Uint16(inData.Next(2))
	hdr.NodeNo, _ = inData.ReadByte()
	hdr.TransSeqNum, _ = inData.ReadByte()

	fc, _ := inData.ReadByte()
	hdr.CommandID, _ = inData.ReadByte()

	if (fc & LTL_FRAMECTL_TYPE_MASK) > 0 {
		hdr.Type_l = LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC
	} else {
		hdr.Type_l = LTL_FRAMECTL_TYPE_PROFILE
	}

	if (fc & LTL_FRAMECTL_DISALBE_DEFAULT_RSP_MASK) > 0 {
		hdr.DisableDefaultRsp = LTL_FRAMECTL_DIS_DEFAULT_RSP_ON
	} else {
		hdr.DisableDefaultRsp = LTL_FRAMECTL_DIS_DEFAULT_RSP_OFF
	}

	if (fc & LTL_FRAMECTL_DIRECTION_MASK) > 0 {
		hdr.Dir = LTL_FRAMECTL_SERVER_CLIENT_DIR
	} else {
		hdr.Dir = LTL_FRAMECTL_CLIENT_SERVER_DIR
	}

	return hdr, inData.Bytes()
}

// 底层,命令发送
func (this *Ltl_t) SendCommand(DstAddr, trunkId uint16, nodeNO, seqNum, hdrctl_type, dir byte,
	disableDefaultRsp bool, cmd byte, cmdFormart []byte) error {
	hdr := encodeHdr(trunkId, nodeNO, seqNum, hdrctl_type, dir, disableDefaultRsp, cmd).buildHdr()
	data := make([]byte, 0, len(hdr)+len(cmdFormart))
	data = append(data, hdr...)
	data = append(data, cmdFormart...)

	return this.WriteMsg(DstAddr, data)
}

// 发送profile级命令
func (this *Ltl_t) SendProfileCmd(DstAddr, trunkId uint16, nodeNO, seqNum,
	cmd byte, cmdFormart []byte) error {
	return this.SendCommand(DstAddr, trunkId, nodeNO, seqNum,
		LTL_FRAMECTL_TYPE_PROFILE, LTL_FRAMECTL_CLIENT_SERVER_DIR, true, cmd, cmdFormart)
}
