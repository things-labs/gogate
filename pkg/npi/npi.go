package npi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/x/numeric"

	"github.com/tarm/serial"
)

//	frame 帧结构 -- | 帧头 | 命令0 | 命令1 | 数据域Adu | xor |
const (
	npi_async_size_max       = 256             // 异步npdu缓冲chan大小
	npi_sync_reponse_timeout = 6 * time.Second // 同步请求超时时间
)

// PDU帧结构
type Npdu struct {
	CmdId uint16 // 命令0 命令1
	Data  []byte // 数据域
}

type Monitor struct {
	sport  *serial.Port
	lctx   context.Context
	cancel context.CancelFunc

	asyncIn             chan *Npdu // 底层mac PDU
	syncResponseLock    sync.Mutex // 同步应答锁
	syncResponseWaiting *uint32    // 原子操作,同步等待应答
	syncResponse        chan *Npdu // 用于通知同步应答的通道
	cb                  sync.Map
}

// npi打开
func Open(cfg *serial.Config) (*Monitor, error) {
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	m := &Monitor{
		sport:               sp,
		lctx:                ctx,
		cancel:              cancel,
		asyncIn:             make(chan *Npdu, npi_async_size_max),
		syncResponseWaiting: new(uint32),
		syncResponse:        make(chan *Npdu, 1),
	}
	go m.asyncframeRoutine() // 启动处理异步NPDU协程
	go m.usartReadRoutine()  // 启动读协程

	return m, nil
}

// 获取上下文关系
func (this *Monitor) Context() context.Context {
	if this.lctx != nil {
		return this.lctx
	}
	return context.Background()
}

// 是否已经启动
func (this *Monitor) HasStarted() bool {
	return this.lctx.Err() == nil
}

// 关闭
func (this *Monitor) Close() {
	this.sport.Close()
	this.cancel()
}

// npi同步请求,得到应答PDU
func (this *Monitor) SendSynchData(commandId uint16, data []byte) (*Npdu, error) {
	this.syncResponseLock.Lock()
	defer this.syncResponseLock.Unlock()
	if err := this.WriteMessage(commandId|MT_RPC_CMD_SREQ, data); err != nil {
		return nil, err
	}

	atomic.StoreUint32(this.syncResponseWaiting, 1) // 开始等待一个同步应答
	tm := time.NewTimer(npi_sync_reponse_timeout)
	select {
	case pdu := <-this.syncResponse:
		tm.Stop()
		return pdu, nil

	case <-tm.C:
		atomic.StoreUint32(this.syncResponseWaiting, 0) // 超时,则不等待
	}

	return nil, errors.New("NPI: syncResponse time out")
}

// npi异步请求
func (this *Monitor) SendAsynchData(commandId uint16, data []byte) error {
	return this.WriteMessage(commandId|MT_RPC_CMD_AREQ, data)
}

// 注册一个异步回调
func (this *Monitor) AddAsyncCb(commandID uint16, cb func(*Npdu)) error {
	if cb == nil {
		return errors.New("NPI: AddAsyncCb provide is nil")
	}
	this.cb.Store(commandID, cb)
	return nil
}

// 批量注册异步回调, 通过map
func (this *Monitor) AddAsyncCbs(cbs map[uint16]func(*Npdu)) error {
	if cbs == nil {
		return errors.New("NPI: AddAsyncCbs provide is nil")
	}
	for k, v := range cbs {
		if _, ok := this.cb.Load(k); !ok {
			this.cb.Store(k, v)
		}
	}
	return nil
}

// 根据id寻找回调
func (this *Monitor) matchAsyncCb(commandID uint16) (func(*Npdu), bool) {
	cb, ok := this.cb.Load(commandID)
	if !ok {
		return nil, false
	}
	return cb.(func(*Npdu)), true
}

//根据id删除回调
func (this *Monitor) DeleteAsyncCb(commandID uint16) {
	this.cb.Delete(commandID)
}

/* 向低层发送一个帧 */
func (this *Monitor) WriteMessage(cmd uint16, data []byte) error {
	wr := make([]byte, 0, len(data)+MT_RPC_FRAME_OVHD+MT_RPC_PDU_HDR_SZ)
	wr = append(wr, MT_RPC_UART_SOF, byte(len(data)), byte(cmd>>8), byte(cmd)) // 0xfe + 长度 + cmd0 + cmd1
	wr = append(wr, data...)                                                   // 数据域
	wr = append(wr, npi_calfcs(wr[MT_RPC_FRAME_HEAD_SZ:]))                     // fcs
	_, err := this.sport.Write(wr)

	return err
}

// 读协程,处理同步应答,异步应答交给其它协程
func (this *Monitor) usartReadRoutine() {
	//logs.Info("npi: Read routine Started!")

	readbuf := make([]byte, 1024)
	tmpbuf := make([]byte, 0)
	for {
		cnt, err := this.sport.Read(readbuf)
		if err != nil { // 处理读错误
			if err != io.EOF {
				//	logs.Error("npi: Read failed, %s", err)
				break
			} else if cnt == 0 { // 文件尾无数据
				continue
			} else { // cnt > 0
				// do nothing
			}
		} else if cnt == 0 { // 长度为0,重读
			continue
		}
		tmpbuf = append(tmpbuf, readbuf[:cnt]...) //和前一次处理的包拼接

		//logs.Debug("npi: Raw count - %d,Raw data: %#v", len(tmpbuf), tmpbuf[0:len(tmpbuf)]) // 要处理的长度和数据
		for {
			//寻找帧头,并返回索引值
			if idx := bytes.IndexByte(tmpbuf, MT_RPC_UART_SOF); idx < 0 { // 找不到.丢弃整条帧
				//logs.Debug("npi: frame sof: not find!")
				tmpbuf = make([]byte, 0)
				break
			} else { // 已经找到frame头
				tmpbuf = tmpbuf[idx:] // 切到frame头后面的所有数据
			}

			if len(tmpbuf) < (MT_RPC_FRAME_HEAD_SZ + MT_RPC_PDU_LEN_SZ) {
				break
			}

			adulen := int(tmpbuf[MT_RPC_FRAME_HEAD_SZ+MT_RPC_POS_LEN]) //获取adu数据长度
			framelen := adulen + MT_RPC_FRAME_OVHD + MT_RPC_PDU_HDR_SZ // 计算总frame帧长
			if len(tmpbuf) < framelen {                                // 帧小于应读帧长,跳出
				break
			}

			//logs.Debug("npi: frame count - %d,frame data: %#v", framelen, tmpbuf[0:framelen]) // 要处理的长度和数据
			// 完整npi帧,foc校验过,成功则处理,失败则丢弃
			if npi_calfcs(tmpbuf[MT_RPC_FRAME_HEAD_SZ:framelen]) == 0 {
				msg := tmpbuf[(MT_RPC_FRAME_HEAD_SZ + MT_RPC_PDU_LEN_SZ):(framelen - MT_RPC_FRAME_FCS_SZ)]

				// 处理pdu帧
				this.Procframe(numeric.BuildUint16(msg[1], msg[0]), msg[2:])
			}

			// foc 校验失败还是成功,跳过或丢弃这个npi帧,获得余下的帧数据
			tmpbuf = tmpbuf[framelen:]
			if len(tmpbuf) < MT_RPC_FRAME_HEAD_SZ+MT_RPC_PDU_LEN_SZ { // 帧小于最小帧长,跳出
				break
			}
		}
	}

	//logs.Error("npi: device closed! and read routine closed") // 设备关闭?
	this.Close()
}

// 处理pdu帧
func (this *Monitor) Procframe(cmd uint16, data []byte) {
	cmdType := cmd & MT_RPC_CMD_TYPE_MASK
	cmdID := cmd & MT_RPC_SUBSYSTEM_MASK
	//logs.Debug("npi: Type is sync: %t, commandID: 0x%04x", cmdType == MT_RPC_CMD_SRSP, cmdID)

	if cmdType == MT_RPC_CMD_SRSP { // synchronous response
		if atomic.CompareAndSwapUint32(this.syncResponseWaiting, 1, 0) { // 相同,则将应答等待置0, 不同不改变
			this.syncResponse <- &Npdu{cmdID, data} // 投递同步Pdu
		}
	} else { // must be an asynchronous message
		this.asyncIn <- &Npdu{cmdID, data} // 投递Mpdu帧
	}
}

// 异步帧处理goroutine
func (this *Monitor) asyncframeRoutine() {
	//logs.Info("npi: asyncframe routine started")

	for {
		select {
		case <-this.lctx.Done():
			//logs.Warn("npi: asyncframe routine closed!")
			return
		case pdu := <-this.asyncIn:
			cb, ok := this.matchAsyncCb(pdu.CmdId)
			if !ok {
				//	logs.Warn("NPI: commandID: 0x%04x not implement", pdu.CmdId)
				break
			}
			cb(pdu)
		}
	}
}

func npi_calfcs(buf []byte) byte {
	xorresult := byte(0)
	for _, dat := range buf {
		xorresult ^= dat
	}
	return xorresult
}
