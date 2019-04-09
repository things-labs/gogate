package models

// 设备产品类型
const (
	PTypes_General = iota // 通用产品
	PTypes_Zigbee         // zigbee产品
	PTypes_Modbus
)

// 所有的产品id列表, 必需注册到DeviceProductInfos, zigbee的产品需要另外注册到zigbee的设备产品里
const (
	PID_DZMS01       = 20000 + iota // LC_DZMS01型号温湿度传感器
	PID_DZSW01                      // LC_DZSW01型号一位智能开关
	PID_DZSW02                      // LC_DZSW02型号二位智能开关
	PID_DZSW03                      // LC_DZSW03型号三位智能开关
	PID_RESERVE      = 0
	PID_ZIGBEE_TEST  = 80000
	PID_GENERAL_TEST = 80001
)

type ProductInfo struct {
	Number           int    // 编号  用于识别相同的类型走的通道 0为默认通道
	Types            int    // 类型
	TypesName        string // 类型名称
	ModelSpec        string // 型号规格
	ModelName        string // 型号名称
	Description      string // 描述
	ManufacturerName string // 制造商名字
}

var DeviceProductInfos map[int]*ProductInfo = map[int]*ProductInfo{
	PID_GENERAL_TEST: &ProductInfo{0, PTypes_General, "general", "LC_GTEST", "普通设备测试", "普通设备测试", "lchtime"},
	PID_DZMS01:       &ProductInfo{0, PTypes_Zigbee, "smart zigbee", "LC_DZMS01", "温湿度传感器", "温湿度传感器", "lchtime"},
	PID_DZSW01:       &ProductInfo{0, PTypes_Zigbee, "smart zigbee", "LC_DZSW01", "一位智能开关", "一开智能开关", "lchtime"},
	PID_DZSW02:       &ProductInfo{0, PTypes_Zigbee, "smart zigbee", "LC_DZSW02", "二位智能开关", "二位智能开关", "lchtime"},
	PID_DZSW03:       &ProductInfo{0, PTypes_Zigbee, "smart zigbee", "LC_DZSW03", "三位智能开关", "三位智能开关", "lchtime"},
}

// 查找产品
func LookupProduct(pid int) (*ProductInfo, error) {
	if pid == PID_RESERVE {
		return nil, ErrProductNotExist
	}
	if v, exist := DeviceProductInfos[pid]; exist {
		return v, nil
	}
	return nil, ErrProductNotExist
}

// 此产品id是否存在
func HasProduct(pid int) bool {
	if pid == PID_RESERVE {
		return false
	}
	_, isexist := DeviceProductInfos[pid]

	return isexist
}
