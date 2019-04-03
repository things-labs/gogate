package models

import (
	_ "github.com/mattn/go-sqlite3"
)

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

// 是否有通用对应的设备
func HasGeneralDevice(pid int, sn string) bool {
	_, err := LookupProduct(pid)
	if err != nil {
		return false
	}

	if db.Where(&GeneralDeviceInfo{ProductId: pid, Sn: sn}).
		First(&GeneralDeviceInfo{}).RecordNotFound() {
		return false
	}
	return true
}

// 创建通用设备
func CreateGeneralDevice(pid int, sn string) error {
	_, err := LookupProduct(pid)
	if err != nil {
		return ErrProductNotExist
	}
	return (&GeneralDeviceInfo{ProductId: pid, Sn: sn}).CreateGeneralDevice()
}

// 创建通用设备
func (this *GeneralDeviceInfo) CreateGeneralDevice() error {
	return db.Create(this).Error
}

// 删除通用设备
func DeleteGeneralDevice(pid int, sn string) error {
	_, err := LookupProduct(pid)
	if err != nil {
		return ErrProductNotExist
	}
	return (&GeneralDeviceInfo{ProductId: pid, Sn: sn}).DeleteGeneralDevice()
}

// 删除通用设备
func (this *GeneralDeviceInfo) DeleteGeneralDevice() error {
	return db.Where(this).Unscoped().Delete(this).Error
}

// 查找通用设备
func FindGeneralDevice(pid int) []GeneralDeviceInfo {
	_, err := LookupProduct(pid)
	if err != nil {
		return []GeneralDeviceInfo{}
	}
	devs := []GeneralDeviceInfo{}
	db.Where(&GeneralDeviceInfo{ProductId: pid}).Find(&devs)
	return devs
}
