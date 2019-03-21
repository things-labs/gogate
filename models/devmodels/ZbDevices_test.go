package devmodels

import (
	"reflect"
	"strings"
	"testing"

	"github.com/slzm40/common"

	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	tsDevNode = ZbDeviceNodeInfo{
		ID:           6655,
		NwkAddr:      1234,
		NodeNo:       5,
		IeeeAddr:     33445566,
		InTrunkList:  "",
		OutTrunkList: "",
		SrcBindList:  "",
		DstBindList:  "",
	}
)

func TestSplitInternalString(t *testing.T) {
	Convey(`splite internal string with ","`, t, func() {
		s := splitInternalString("1")
		So(reflect.DeepEqual(s, []string{"1"}), ShouldBeTrue)
		s = splitInternalString("")
		So(len(s), ShouldBeZeroValue)
	})
}

func TestJoinInternalString(t *testing.T) {
	Convey(`join internal string with ","`, t, func() {
		s := joinInternalString([]string{"1", "2", "3", "4", "5"})
		So(strings.EqualFold(s, "1,2,3,4,5"), ShouldBeTrue)
		s = joinInternalString([]string{})
		So(strings.EqualFold(s, ""), ShouldBeTrue)
	})

}

func TestparseInternalString(t *testing.T) {
	Convey("解析内部的String 逗号分隔", t, func() {
		Convey("解析内部的String 逗号分隔 - 列表均为空", func() {
			tsDevNode.parseInternalString()
			So(len(tsDevNode.inTrunk), ShouldBeZeroValue)
			So(len(tsDevNode.outTrunk), ShouldBeZeroValue)
			So(len(tsDevNode.srcBind), ShouldBeZeroValue)
			So(len(tsDevNode.dstBind), ShouldBeZeroValue)
		})

		Convey("解析内部的String 逗号分隔 - 列表均有值", func() {
			expTK := []string{"1", "2", "3", "4"}
			expTK1 := []string{"5", "6", "7", "8"}
			exp := []string{"1", "3", "5", "7"}
			exp1 := []string{"2", "4", "6", "8"}
			tsDevNode.InTrunkList = `1,2,3,4`
			tsDevNode.OutTrunkList = `5,6,7,8`
			tsDevNode.SrcBindList = `1,3,5,7`
			tsDevNode.DstBindList = `2,4,6,8`

			tsDevNode.parseInternalString()
			So(reflect.DeepEqual(tsDevNode.inTrunk, expTK), ShouldBeTrue)
			So(reflect.DeepEqual(tsDevNode.outTrunk, expTK1), ShouldBeTrue)
			So(reflect.DeepEqual(tsDevNode.srcBind, exp), ShouldBeTrue)
			So(reflect.DeepEqual(tsDevNode.dstBind, exp1), ShouldBeTrue)
		})
	})
}

func TestGetDeviceNodeInfo(t *testing.T) {
	Convey("获取设备节点属性", t, func() {

		tsDevNode.InTrunkList = `1,2,3,4`
		tsDevNode.OutTrunkList = `1,2,3,4`
		tsDevNode.SrcBindList = `1,2,3,4`
		tsDevNode.DstBindList = `1,2,3,4`
		expTK := []string{"1", "2", "3", "4"}
		exp := []string{"1", "2", "3", "4"}
		tsDevNode.parseInternalString()

		Convey("获取id", func() {
			So(tsDevNode.GetID(), ShouldEqual, 6655)
		})
		Convey("获取网络地址", func() {
			So(tsDevNode.GetNwkAddr(), ShouldEqual, 1234)
		})
		Convey("获取ieee地址", func() {
			So(tsDevNode.GetIeeeAddr(), ShouldEqual, 33445566)
		})
		Convey("获取节点号", func() {
			So(tsDevNode.GetNodeNum(), ShouldEqual, 5)
		})
		Convey("获取输入输出集", func() {
			inTk1, outTk1 := tsDevNode.GetTrunkIDList()
			So(reflect.DeepEqual(inTk1, expTK), ShouldBeTrue)
			So(reflect.DeepEqual(outTk1, expTK), ShouldBeTrue)
		})
		Convey("获取源和目标绑定表", func() {
			srcBd, dstBd := tsDevNode.GetBindList()
			So(reflect.DeepEqual(srcBd, exp), ShouldBeTrue)
			So(reflect.DeepEqual(dstBd, exp), ShouldBeTrue)
		})
	})
}

func TestGetDeviceInfo(t *testing.T) {
	var dev = &ZbDeviceInfo{
		Model: gorm.Model{
			ID: 1,
		},
		IeeeAddr:  11223344,
		NwkAddr:   5566,
		Capacity:  2,
		ProductId: 80000,
	}

	Convey("获取设备信息", t, func() {
		Convey("获取设备ieee地址", func() {
			So(dev.GetIeeeAddr(), ShouldEqual, 11223344)
		})
		Convey("获取设备网络地址", func() {
			So(dev.GetNwkAddr(), ShouldEqual, 5566)
		})
		Convey("获取设备能力", func() {
			So(dev.GetCapacity(), ShouldEqual, 2)
		})
		Convey("获取设备产品id", func() {
			So(dev.GetProductID(), ShouldEqual, 80000)
		})

		Convey("获取设备数据据ID值", func() {
			So(dev.GetID(), ShouldEqual, 1)
		})
	})
}

var tsDev = &ZbDeviceInfo{
	IeeeAddr:  11223344,
	NwkAddr:   5566,
	Capacity:  2,
	ProductId: 80000,
}

var tsDev0 = &ZbDeviceInfo{
	IeeeAddr:  55667788,
	NwkAddr:   1122,
	Capacity:  1,
	ProductId: 80000,
}

func TestDevice(t *testing.T) {
	Convey("设备表和设备节点表", t, func() {
		Convey("创建设备和所有的节点", func() {
			err := UpdateZbDeviceAndNode(11223344, 5566, 2, 80000)
			So(err, ShouldBeNil)

			err = UpdateZbDeviceAndNode(55667788, 1122, 1, 80000)
			So(err, ShouldBeNil)

			err = UpdateZbDeviceAndNode(11111111, 1111, 2, 80000)
			So(err, ShouldBeNil)

			err = UpdateZbDeviceAndNode(22222222, 2222, 1, 80000)
			So(err, ShouldBeNil)
		})

		Convey("通过网络地址查询设备", func() {
			oDev, err := LookupZbDeviceByNwkAddr(5566)
			So(err, ShouldBeNil)
			So(oDev.IeeeAddr, ShouldEqual, tsDev.IeeeAddr)
			So(oDev.Capacity, ShouldEqual, tsDev.Capacity)
			So(oDev.ProductId, ShouldEqual, tsDev.ProductId)
		})

		Convey("通过ieee地址查询设备", func() {
			oDev, err := LookupZbDeviceByIeeeAddr(11223344)
			So(err, ShouldBeNil)
			So(oDev.NwkAddr, ShouldEqual, tsDev.NwkAddr)
			So(oDev.Capacity, ShouldEqual, tsDev.Capacity)
			So(oDev.ProductId, ShouldEqual, tsDev.ProductId)
		})

		Convey("通过网络地址,节点查询设备节点", func() {
			oDevNode, err := LookupZbDeviceNodeByNN(5566, 0)
			So(err, ShouldBeNil)
			So(oDevNode.ID, ShouldBeGreaterThan, 0)
		})

		Convey("通过Ieee地址,节点查询设备节点", func() {
			oDevNode, err := LookupZbDeviceNodeByIN(11223344, 1)
			So(err, ShouldBeNil)
			So(oDevNode.ID, ShouldBeGreaterThan, 0)
		})

		Convey("通过ID查询设备节点", func() {
			o1, _ := LookupZbDeviceNodeByNN(5566, 0)

			o2, err := LookupZbDeviceNodeByID(common.FormatBaseTypes(o1.ID))

			So(err, ShouldBeNil)
			So(reflect.DeepEqual(o1, o2), ShouldBeTrue)
		})

		Convey("更新设备能力属性", func() {
			testDev, err := LookupZbDeviceByNwkAddr(5566)
			So(err, ShouldBeNil)

			err = testDev.updateCapacity(1)
			So(err, ShouldBeNil)
			err = testDev.updateCapacity(tsDev.Capacity)
		})

		Convey("更新设备和设备所有设备节点的网络地址", func() {
			testDev, err := LookupZbDeviceByNwkAddr(5566)
			So(err, ShouldBeNil)

			err = testDev.updateZbDeviceAndNodeNwkAddr(7788)
			So(err, ShouldBeNil)
		})

		Convey("绑定两个互补的设备节点", func() {
			err := BindZbDeviceNode(11223344, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(11223344, 2, 11223344, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(55667788, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(11111111, 2, 22222222, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(22222222, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(11111111, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
		})

		Convey("根据网络址,节点号查找绑定表的所有设备节点", func() {
			devnodes, err := BindFindZbDeviceNodeByNN(7788, 2, 3)
			So(err, ShouldBeNil)
			So(len(devnodes), ShouldEqual, 2)
		})

		Convey("根据ieee地址,节点号查找绑定表的所有设备节点", func() {
			devnodes, err := BindFindZbDeviceNodeByIN(11223344, 2, 3)
			So(err, ShouldBeNil)
			So(len(devnodes), ShouldEqual, 2)
		})

		Convey("解除两个设备节点的绑定", func() {
			err := UnZbBindDeviceNode(11223344, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
			err = UnZbBindDeviceNode(11223344, 2, 11223344, 3, 3)
			So(err, ShouldBeNil)
		})

		Convey("删除设备", func() {
			err := DeleteZbDeveiceAndNode(11111111)
			So(err, ShouldBeNil)
			err = DeleteZbDeveiceAndNode(22222222)
			So(err, ShouldBeNil)
		})

	})
}
