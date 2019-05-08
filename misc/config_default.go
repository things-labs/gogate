package misc

const (
	SmartAppDefaultCfg = `#enable orm DB log 
ormDbLog=false

[logs]
# 调试输出引擎(conn/console)
adapter = console
# 调试等级
level = 7
# 使能文件或行号输出
isEFCD = true
# 使能异步
isAsync = false
#if 0,use default(2), other user. seeSetLogFuncCallDepth
logFCD = 0
#使用conn时配置
net = udp
addr = 192.168.1.199:9000
`
	UsartDefaultCfg = `[COM0]
Name=/dev/ttyS1
BaudRate=115200
DataBit=8
Parity=N
StopBit=1
FlowType=0
`
)
