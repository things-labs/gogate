package models

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tsDev1 = DeviceInfo{IeeeAddr: 33445566, NwkAddr: 1234, Capacity: 1, ProductId: 100}
	tsDev2 = DeviceInfo{IeeeAddr: 11223344, NwkAddr: 5678, Capacity: 2, ProductId: 100}
)

func TestDevll(t *testing.T) {
	Convey("设备表", t, func() {

		// Convey("增加更新设备信息", func() {
		// 	Isneed, err := AddDevll(tsDev1.IeeeAddr, tsDev1.NwkAddr, tsDev1.Capacity, tsDev1.ProductId-1)
		// 	So(Isneed, ShouldBeTrue)
		// 	So(err, ShouldBeNil)

		// 	Isneed, err = AddDevll(tsDev1.IeeeAddr, tsDev1.NwkAddr, tsDev1.Capacity, tsDev1.ProductId)
		// 	So(Isneed, ShouldBeTrue)
		// 	So(err, ShouldBeNil)

		// 	Isneed, err = AddDevll(tsDev1.IeeeAddr, tsDev1.NwkAddr, tsDev1.Capacity, tsDev1.ProductId)
		// 	So(Isneed, ShouldBeFalse)
		// 	So(err, ShouldBeNil)

		// 	Isneed, err = AddDevll(tsDev2.IeeeAddr, tsDev2.NwkAddr, tsDev2.Capacity, tsDev2.ProductId)
		// 	So(Isneed, ShouldBeFalse)
		// 	So(err, ShouldBeNil)
		// })

		// Convey("用nwkAddr查找设备", func() {
		// 	dev, err := LookupDevllByNwkAddr(tsDev1.NwkAddr)
		// 	So(err, ShouldBeNil)
		// 	So(dev.GetIeeeAddr(), ShouldEqual, tsDev1.IeeeAddr)
		// 	So(dev.GetNwkAddr(), ShouldEqual, tsDev1.NwkAddr)
		// 	So(dev.GetCapacity(), ShouldEqual, tsDev1.Capacity)
		// 	So(dev.GetProductID(), ShouldEqual, tsDev1.ProductId)
		// 	So(dev.GetID(), ShouldBeGreaterThan, 0)

		// 	_, err = LookupDevllByNwkAddr(0xffff)
		// 	So(err, ShouldNotBeNil)
		// })

		// Convey("用ieeeAddr查找设备", func() {
		// 	dev, err := LookupDevllByIeeeAddr(tsDev1.IeeeAddr)
		// 	So(err, ShouldBeNil)
		// 	So(dev.GetIeeeAddr(), ShouldEqual, tsDev1.IeeeAddr)
		// 	So(dev.GetNwkAddr(), ShouldEqual, tsDev1.NwkAddr)
		// 	So(dev.GetCapacity(), ShouldEqual, tsDev1.Capacity)
		// 	So(dev.GetProductID(), ShouldEqual, tsDev1.ProductId)

		// 	_, err = LookupDevllByIeeeAddr(0xffff)
		// 	So(err, ShouldNotBeNil)
		// })

		// Convey("删除设备", func() {
		// 	Isneed, err := DeleteDevll(tsDev1.IeeeAddr)
		// 	So(Isneed, ShouldBeTrue)
		// 	So(err, ShouldBeNil)

		// 	Isneed, err = DeleteDevll(0xffff)
		// 	So(Isneed, ShouldBeFalse)
		// 	So(err, ShouldBeNil)
		// })

	})
}

var inTrunkStr = `{"trunkID":[6,7,8,9]}`
var outTrunkStr = `{"trunkID":[3,4,5,6]}`
var bindListStr = `{"id":[1,2,3,4]}`
var inTrunkSlice = []uint16{6, 7, 8, 9}
var outTrunkSlice = []uint16{3, 4, 5, 6}
var bindListSlice = []uint{1, 2, 3, 4}

var nodeInfo0 = &DeviceNodeInfo{
	NwkAddr:      1234,
	NodeNo:       1,
	IeeeAddr:     1122334455667788,
	InTrunkList:  inTrunkStr,
	OutTrunkList: outTrunkStr,
	DstBindList:  bindListStr,
}

func TestNbiNode(t *testing.T) {
	Convey("节点表", t, func() {
		Convey("获取绑定id列表,json字符串转换成Id列表", func() {
			l, err := nodeInfo0.GetDstBindList()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(l, bindListSlice), ShouldBeTrue)
		})
		Convey("设置绑定id列表,id列表转换成json字符串", func() {
			actual := &DeviceNodeInfo{}
			err := actual.setDstBindList(bindListSlice)
			So(err, ShouldBeNil)
			So(strings.Compare(actual.DstBindList, nodeInfo0.DstBindList), ShouldBeZeroValue)
		})

		Convey("添加新的绑定到id列表", func() {

		})
		Convey("从绑定id列表删除一个绑定", func() {

		})
		Convey("获得节点输入输出集列表", func() {
			lin, lout, err := nodeInfo0.GetTrunkIDList()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(lin, inTrunkSlice), ShouldBeTrue)
			So(reflect.DeepEqual(lout, outTrunkSlice), ShouldBeTrue)
		})

		Convey("设置节点输入输出集列表", func() {
			actual := &DeviceNodeInfo{}
			err := actual.SetTrunkIDlist(inTrunkSlice, outTrunkSlice)
			So(err, ShouldBeNil)
			So(strings.Compare(actual.InTrunkList, nodeInfo0.InTrunkList), ShouldBeZeroValue)
			So(strings.Compare(actual.OutTrunkList, nodeInfo0.OutTrunkList), ShouldBeZeroValue)
		})
	})
}

func TestNBI(t *testing.T) {
	Convey("节点增删改查", t, func() {

	})
}
