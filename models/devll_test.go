package models

import (
	"fmt"
	//"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tsDev1 = DeviceInfo{IeeeAddr: 33445566, NwkAddr: 1234, Capacity: 1, ProductId: 100}
	tsDev2 = DeviceInfo{IeeeAddr: 11223344, NwkAddr: 5678, Capacity: 2, ProductId: 100}
)

func TestDevll(t *testing.T) {
	Convey("设备表", t, func() {

		Convey("增加更新设备信息", func() {
			err := InserUpdateDevll(tsDev1.IeeeAddr, tsDev1.NwkAddr, tsDev1.Capacity, tsDev1.ProductId-1)
			So(err, ShouldBeNil)

			err = InserUpdateDevll(tsDev1.IeeeAddr, tsDev1.NwkAddr, tsDev1.Capacity, tsDev1.ProductId)
			So(err, ShouldBeNil)

			err = InserUpdateDevll(tsDev2.IeeeAddr, tsDev2.NwkAddr, tsDev2.Capacity, tsDev2.ProductId)
			So(err, ShouldBeNil)
		})

		Convey("用nwkAddr查找设备", func() {
			dev, err := FindDevllByNwk(tsDev1.NwkAddr)
			So(err, ShouldBeNil)
			So(dev.GetIeeeAddr(), ShouldEqual, tsDev1.IeeeAddr)
			So(dev.GetNwkAddr(), ShouldEqual, tsDev1.NwkAddr)
			So(dev.GetCapacity(), ShouldEqual, tsDev1.Capacity)
			So(dev.GetProductID(), ShouldEqual, tsDev1.ProductId)
			So(dev.GetID(), ShouldBeGreaterThan, 0)

			_, err = FindDevllByNwk(0xffff)
			So(err, ShouldNotBeNil)
		})

		Convey("用ieeeAddr查找设备", func() {
			dev, err := FindDevllByIeeeAddr(tsDev1.IeeeAddr)
			So(err, ShouldBeNil)
			So(dev.GetIeeeAddr(), ShouldEqual, tsDev1.IeeeAddr)
			So(dev.GetNwkAddr(), ShouldEqual, tsDev1.NwkAddr)
			So(dev.GetCapacity(), ShouldEqual, tsDev1.Capacity)
			So(dev.GetProductID(), ShouldEqual, tsDev1.ProductId)

			_, err = FindDevllByIeeeAddr(0xffff)
			So(err, ShouldNotBeNil)
		})

		Convey("删除设备", func() {
			err := DeleteDevll(tsDev1.IeeeAddr)
			So(err, ShouldBeNil)

			err = DeleteDevll(0xffff)
			So(err, ShouldBeNil)
		})

	})
}

func TestGetNbiBindList(t *testing.T) {
	Convey("节点表", t, func() {
		Convey("获取绑定id列表", func() {
			ifo := &NodeInfo{BindList: `{"id":[1,2,3,4]}`}
			l, err := ifo.GetNbiBindList()
			So(err, ShouldBeNil)

			fmt.Printf("%#v", l)
		})
	})
}

func TestNBI(t *testing.T) {
	Convey("节点增删改查", t, func() {

	})
}
