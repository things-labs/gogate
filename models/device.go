package models

type GeneralDeviceInfo struct {
	ID        uint
	ProductId int
	Sn        string
}

const generaldeviceInfo_sql = `CREATE TABLE "general_device_infos" (
			"id" integer primary key autoincrement,
			"product_id" integer NOT NULL,
			"sn" bigint NOT NULL,
			UNIQUE(product_id,sn) ON CONFLICT FAIL)`

func init() {
	RegisterDbTableInitFunction(GeneralDeviceDbTableInit)
}

// 通用设备数据表初始化
func GeneralDeviceDbTableInit() error {
	if !db.HasTable("general_device_infos") {
		return db.Raw(generaldeviceInfo_sql).Scan(&GeneralDeviceInfo{}).Error
	}
	return nil
}

// 是否有通用对应的设备(pid sn)
func HasGeneralDevice(pid int, sn string) bool {
	if !HasProduct(pid) || len(sn) == 0 {
		return false
	}
	return db.Where(&GeneralDeviceInfo{ProductId: pid, Sn: sn}).
		First(&GeneralDeviceInfo{}).Error == nil
}

// 查找通用设备对应的信息
func LookupGeneralDevice(pid int, sn string) (*GeneralDeviceInfo, error) {
	if !HasProduct(pid) || len(sn) == 0 {
		return nil, ErrDeviceNotExist
	}

	dev := &GeneralDeviceInfo{}
	err := db.Where(&GeneralDeviceInfo{ProductId: pid, Sn: sn}).First(dev).Error
	return dev, err
}

// 创建通用设备
func CreateGeneralDevice(pid int, sn string) error {
	return (&GeneralDeviceInfo{ProductId: pid, Sn: sn}).CreateGeneralDevice()
}

// 删除通用设备
func DeleteGeneralDevice(pid int, sn string) error {
	return (&GeneralDeviceInfo{ProductId: pid, Sn: sn}).DeleteGeneralDevice()
}

// 创建通用设备
func (this *GeneralDeviceInfo) CreateGeneralDevice() error {
	if HasGeneralDevice(this.ProductId, this.Sn) {
		return nil
	}
	return db.Create(this).Error
}

// 删除通用设备
func (this *GeneralDeviceInfo) DeleteGeneralDevice() error {
	if HasGeneralDevice(this.ProductId, this.Sn) {
		return db.Where(this).Unscoped().Delete(this).Error
	}
	return nil
}

// 查找通用设备列表
func FindGeneralDevice(pid int) []GeneralDeviceInfo {
	devs := []GeneralDeviceInfo{}
	if !HasProduct(pid) {
		return devs
	}
	db.Where(&GeneralDeviceInfo{ProductId: pid}).Find(&devs)
	return devs
}
