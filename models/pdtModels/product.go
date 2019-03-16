package pdtModels

// 产品类型
const (
	ProductTypes_General = iota // 通用产品
	ProductTypes_Zigbee         // zigbee产品
	ProductTypes_Modbus
)

type ProductInfo struct {
	Number           int    // 编号  用于识别相同的类型走的通道
	Types            int    // 类型
	TypesName        string // 类型名称
	ModelSpec        string // 型号规格
	ModelName        string // 型号名称
	Description      string // 描述
	ManufacturerName string // 制造商名字
}

var DeviceProductInfos map[int]*ProductInfo = map[int]*ProductInfo{
	20000: &ProductInfo{0, ProductTypes_Zigbee, "smart zigbee", "sws01", "创1开关", "一开智能开关", "lchtime"},
	20001: &ProductInfo{0, ProductTypes_Zigbee, "smart zigbee", "sws02", "创2开关", "二开智能开关", "lchtime"},
	20002: &ProductInfo{0, ProductTypes_Zigbee, "smart zigbee", "sws03", "创3开关", "三开智能开关", "lchtime"},
	20003: &ProductInfo{0, ProductTypes_Zigbee, "smart zigbee", "wsd00", "创2温湿度", "简易温湿度传感器", "lchtime"},
	20100: &ProductInfo{0, ProductTypes_Modbus, "modbus", "测试型号", "sensor测试", "mb测试", "lchtime"},
}

// 查找产品
func LookupProduct(pid int) (*ProductInfo, bool) {
	v, isexist := DeviceProductInfos[pid]

	return v, isexist
}

// 此产品id是否存在
func HasProduct(pid int) bool {
	_, isexist := DeviceProductInfos[pid]

	return isexist
}
