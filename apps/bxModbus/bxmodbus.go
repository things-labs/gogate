package bxModbus

import (
	"container/list"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/goburrow/modbus"
)

const (
	scan_rate_time    = 20 * time.Millisecond // ms  请求列表扫描速率
	ready_request_cnt = 50                    // 就绪请求最大数量
)

type Config struct {
	// Device path (/dev/ttyS0)
	Name string
	// Baud rate (default 19200)
	BaudRate int
	// Data bits: 5, 6, 7 or 8 (default 8)
	DataBits int
	// Stop bits: 1 or 2 (default 1)
	StopBits int
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	Parity string
}

// node info
type NodeReg struct {
	Holding_addr_start  uint16
	Input_addr_start    uint16
	Coils_addr_start    uint16
	Discrete_addr_start uint16
	HoldingReg          []uint16
	InputReg            []uint16
	CoilsReg            []uint8
	DiscreteInputsReg   []uint8
}

// object
type Mbj struct {
	handler  *modbus.RTUClientHandler
	client   modbus.Client     // 客户端
	node     sync.Map          // 节点信息存储表
	reqlist  *list.List        // 请求列表
	locklist sync.Mutex        // 请求列表锁
	readyReq chan *pollRequest // 就绪表
}

// read poll request
type pollRequest struct {
	slaveId  byte
	funcCode byte
	regAddr  uint16
	quantity uint16
	scanrate time.Duration // scan rate
	scancnt  time.Duration // scan cnt
	txcnt    uint64        // tx count
	errcnt   uint64        // error count
}

type parseRspHandler struct {
}

func BxRtuInit(cfg *Config, timeout time.Duration) (*Mbj, error) {

	handler := modbus.NewRTUClientHandler(cfg.Name)
	handler.BaudRate = cfg.BaudRate
	handler.DataBits = cfg.DataBits
	handler.Parity = cfg.Parity
	handler.StopBits = cfg.StopBits
	handler.Timeout = timeout
	handler.SlaveId = 0xff

	if err := handler.Connect(); err != nil {
		return nil, err
	} else {
		client := modbus.NewClient(handler)
		bmi := &Mbj{
			handler:  handler,
			client:   client,
			reqlist:  list.New(),
			readyReq: make(chan *pollRequest, ready_request_cnt),
		}

		go bmi.readPoll()
		go bmi.scanRequestList()
		logs.Debug("modbus rtu started")

		return bmi, nil
	}
}

// 客户端新建一个节点
func (bmi *Mbj) BxNewNode(slaveId byte, holding_addr_start, holding_num, input_addr_start, input_num,
	coils_addr_start, coils_num, discrete_addr_start, discrete_num uint16) {
	var coilsByte, DiscsByte uint16

	if (byte(coils_num) & 0x07) > 0 {
		coilsByte = (coils_num >> 3) + 1
	} else {
		coilsByte = coils_num >> 3
	}
	if (byte(discrete_num) & 0x07) > 0 {
		DiscsByte = (discrete_num >> 3) + 1
	} else {
		DiscsByte = discrete_num >> 3
	}

	newnode := &NodeReg{
		Holding_addr_start:  holding_addr_start,
		Input_addr_start:    input_addr_start,
		Coils_addr_start:    coils_addr_start,
		Discrete_addr_start: discrete_addr_start,
		HoldingReg:          make([]uint16, holding_num),
		InputReg:            make([]uint16, input_num),
		CoilsReg:            make([]uint8, coilsByte),
		DiscreteInputsReg:   make([]uint8, DiscsByte),
	}

	bmi.node.Store(slaveId, newnode)
}

// 从客户端删除一个节点
func (bmi *Mbj) BxDeleteNode(slaveId byte) {
	bmi.node.Delete(slaveId)
}

// 客户端获得节点寄存器列表
func (bmi *Mbj) BxGetNode(slaveId byte) (*NodeReg, bool) {
	nd, ok := bmi.node.Load(slaveId)
	if ok {
		nds, ok := nd.(*NodeReg)
		return nds, ok
	} else {
		return nil, false
	}
}

// 新建一个读轮询功能码
func (bmi *Mbj) BxNewReadPoll(slaveId, funcCode byte, address, quantity uint16, scanrate time.Duration) *list.Element {

	if scanrate < scan_rate_time {
		scanrate = scan_rate_time
	}

	req := &pollRequest{
		slaveId:  slaveId,
		funcCode: funcCode,
		regAddr:  address,
		quantity: quantity,
		scanrate: scanrate,
	}

	bmi.locklist.Lock()
	defer bmi.locklist.Unlock()

	return bmi.reqlist.PushBack(req)
}

// 删除一个读轮询功能码
func (bmi *Mbj) BxDeleteReadPoll(req *list.Element) {
	bmi.locklist.Lock()
	defer bmi.locklist.Unlock()

	bmi.reqlist.Remove(req)
}

//  写单个线圈
func (bmi *Mbj) WriteSingleCoil(slaveId byte, address, value uint16) ([]byte, error) {
	bmi.handler.SlaveId = slaveId
	return bmi.client.WriteSingleCoil(address, value)
}

// 写多个线圈
func (bmi *Mbj) WriteMultipleCoils(slaveId byte, address, quantity uint16, value []byte) ([]byte, error) {
	bmi.handler.SlaveId = slaveId
	return bmi.client.WriteMultipleCoils(address, quantity, value)
}

// 写单个寄存器
func (bmi *Mbj) WriteSingleRegister(slaveId byte, address, value uint16) ([]byte, error) {
	bmi.handler.SlaveId = slaveId
	return bmi.client.WriteSingleRegister(address, value)
}

// 写多个寄存器
func (bmi *Mbj) WriteMultipleRegisters(slaveId byte, address, quantity uint16, value []byte) (results []byte, err error) {
	bmi.handler.SlaveId = slaveId
	return bmi.client.WriteMultipleRegisters(address, quantity, value)
}

// 读写多个寄存器
func (bmi *Mbj) ReadWriteMultipleRegisters(slaveId byte, readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	bmi.handler.SlaveId = slaveId
	return bmi.client.ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity, value)
}

// 掩码写寄存器
func (bmi *Mbj) MaskWriteRegister(slaveId byte, address, andMask, orMask uint16) (results []byte, err error) {
	bmi.handler.SlaveId = slaveId
	return bmi.client.MaskWriteRegister(address, andMask, orMask)
}

// 扫描请求列表
func (bmi *Mbj) scanRequestList() {
	var req *pollRequest
	var tmp *list.Element

	for {
		<-time.After(scan_rate_time)

		bmi.locklist.Lock()
		for e := bmi.reqlist.Front(); e != nil; e = tmp {
			req = e.Value.(*pollRequest)
			req.scancnt += scan_rate_time
			if req.scancnt > req.scanrate {
				req.scancnt = 0
				tmp = e.Next()
				bmi.reqlist.Remove(e)
				bmi.locklist.Unlock()
				bmi.readyReq <- req
				bmi.locklist.Lock()
			} else {
				tmp = e.Next()
			}
		}
		bmi.locklist.Unlock()
	}

}

// 读协程
func (bmi *Mbj) readPoll() {
	var (
		results []byte
		cureq   *pollRequest
		err     error
	)

	for {
		cureq = <-bmi.readyReq // 查看是否有准备好的请求

		cureq.txcnt++
		bmi.handler.SlaveId = cureq.slaveId
		switch cureq.funcCode {
		// Bit access read
		case modbus.FuncCodeReadCoils:
			results, err = bmi.client.ReadCoils(cureq.regAddr, cureq.quantity)
		case modbus.FuncCodeReadDiscreteInputs:
			results, err = bmi.client.ReadDiscreteInputs(cureq.regAddr, cureq.quantity)
			// 16-bit access read
		case modbus.FuncCodeReadHoldingRegisters:
			results, err = bmi.client.ReadHoldingRegisters(cureq.regAddr, cureq.quantity)
		case modbus.FuncCodeReadInputRegisters:
			results, err = bmi.client.ReadInputRegisters(cureq.regAddr, cureq.quantity)
			// FIFO read
		case modbus.FuncCodeReadFIFOQueue:
			results, err = bmi.client.ReadFIFOQueue(cureq.regAddr)

			//		// Bit access write 采用直接操作,不直行Poll
			//		case modbus.FuncCodeWriteSingleCoil:
			//		case modbus.FuncCodeWriteMultipleCoils:
			//		case modbus.FuncCodeWriteSingleRegister:
			//		case modbus.FuncCodeWriteMultipleRegisters:
			//		case modbus.FuncCodeReadWriteMultipleRegisters:
			//		case modbus.FuncCodeMaskWriteRegister:
		}

		if err != nil {
			logs.Error("Read failed", err)
			cureq.errcnt++
		}

		logs.Debug("---------------------------------------------------------------")
		logs.Debug("BxReadPoll: Tx=%d,Err=%d,ID=%d,F=%d,Addr=%d,Cnt=%d,SR=%dms",
			cureq.txcnt, cureq.errcnt, cureq.slaveId, cureq.funcCode,
			cureq.regAddr, cureq.quantity, cureq.scanrate/time.Millisecond)
		logs.Debug("funcCode: %x, Value: %v", cureq.funcCode, results)

		if cureq.scanrate > 0 {
			bmi.reqlist.PushFront(cureq)
		}

	}
}
