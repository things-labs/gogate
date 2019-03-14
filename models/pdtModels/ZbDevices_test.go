package pdtModels

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	tsDevNode = ZbDeviceNodeInfo{
		ID:           6655,
		NwkAddr:      1234,
		NodeNo:       5,
		IeeeAddr:     33445566,
		InTrunkList:  _default_trunkid_list,
		OutTrunkList: _default_trunkid_list,
		SrcBindList:  _default_bind_list,
		DstBindList:  _default_bind_list,
	}
)

func TestparseInternalJsonString(t *testing.T) {
	Convey("解析内部的jsonString", t, func() {
		Convey("解析内部的jsonString - 列表均为空", func() {
			err := tsDevNode.parseInternalJsonString()
			So(err, ShouldBeNil)
			So(len(tsDevNode.inTrunk), ShouldBeZeroValue)
			So(len(tsDevNode.outTrunk), ShouldBeZeroValue)
			So(len(tsDevNode.srcBind), ShouldBeZeroValue)
			So(len(tsDevNode.dstBind), ShouldBeZeroValue)
		})

		Convey("解析内部的jsonString - 列表均有值", func() {
			expTK := []uint16{1, 2, 3, 4}
			expTK1 := []uint16{5, 6, 7, 8}
			exp := []uint{1, 3, 5, 7}
			exp1 := []uint{2, 4, 6, 8}
			tsDevNode.InTrunkList = `{"trunkID":[1,2,3,4]}`
			tsDevNode.OutTrunkList = `{"trunkID":[5,6,7,8]}`
			tsDevNode.SrcBindList = `{"id":[1,3,5,7]}`
			tsDevNode.DstBindList = `{"id":[2,4,6,8]}`

			err := tsDevNode.parseInternalJsonString()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(tsDevNode.inTrunk, expTK), ShouldBeTrue)
			So(reflect.DeepEqual(tsDevNode.outTrunk, expTK1), ShouldBeTrue)
			So(reflect.DeepEqual(tsDevNode.srcBind, exp), ShouldBeTrue)
			So(reflect.DeepEqual(tsDevNode.dstBind, exp1), ShouldBeTrue)
		})

		Convey("获取设备节点Id列表 - 列表值空字符串或错误Json格式", func() {
			tsDevNode.InTrunkList = ""
			tsDevNode.OutTrunkList = ""
			tsDevNode.SrcBindList = ""
			tsDevNode.DstBindList = ""

			err := tsDevNode.parseInternalJsonString()
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetDeviceNodeInfo(t *testing.T) {
	Convey("获取设备节点属性", t, func() {
		tsDevNode.InTrunkList = `{"trunkID":[1,2,3,4]}`
		tsDevNode.OutTrunkList = `{"trunkID":[1,2,3,4]}`
		tsDevNode.SrcBindList = `{"id":[1,2,3,4]}`
		tsDevNode.DstBindList = `{"id":[1,2,3,4]}`
		expTK := []uint16{1, 2, 3, 4}
		exp := []uint{1, 2, 3, 4}
		err := tsDevNode.parseInternalJsonString()
		So(err, ShouldBeNil)

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
		Convey("获取源绑定表", func() {
			bd := tsDevNode.GetSrcBindList()
			So(reflect.DeepEqual(bd, exp), ShouldBeTrue)
		})

		Convey("获取目标绑定表", func() {
			bd := tsDevNode.GetDstBindList()
			So(reflect.DeepEqual(bd, exp), ShouldBeTrue)
		})
	})
}

func TestSetTrunkIDList(t *testing.T) {
	Convey("设置设备节点集id列表", t, func() {
		spareTk := []uint16{}

		Convey("设置设备节点Id列表 - 列表均为空", func() {
			err := tsDevNode.SetTrunkIDlist(spareTk, spareTk)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.InTrunkList, _default_trunkid_list), ShouldBeZeroValue)
			So(strings.Compare(tsDevNode.InTrunkList, _default_trunkid_list), ShouldBeZeroValue)
		})

		Convey("设置设备节点Id列表 - 列表均有值", func() {
			expTK := []uint16{1, 2, 3, 4}
			expStr := `{"trunkID":[1,2,3,4]}`

			err := tsDevNode.SetTrunkIDlist(expTK, expTK)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.InTrunkList, expStr), ShouldBeZeroValue)
			So(strings.Compare(tsDevNode.OutTrunkList, expStr), ShouldBeZeroValue)
		})

		Convey("设置设备节点Id列表 - 列表为nil", func() {
			err := tsDevNode.SetTrunkIDlist(nil, nil)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.InTrunkList, _default_trunkid_list), ShouldBeZeroValue)
			So(strings.Compare(tsDevNode.OutTrunkList, _default_trunkid_list), ShouldBeZeroValue)
		})
	})
}

func TestSetBindList(t *testing.T) {
	Convey("设置设备节点绑定列表", t, func() {
		spare := []uint{}
		exp := []uint{1, 2, 3, 4}
		expStr := `{"id":[1,2,3,4]}`

		Convey("设置设备节点: 目的绑定列表 - 列表均为空", func() {
			err := tsDevNode.setDstBindList(spare)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.DstBindList, _default_bind_list), ShouldBeZeroValue)
		})

		Convey("设置设备节点: 目的绑定列表 - 列表均有值", func() {
			err := tsDevNode.setDstBindList(exp)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.DstBindList, expStr), ShouldBeZeroValue)
		})

		Convey("设置设备节点: 目的绑定列表 - 列表为nil", func() {
			err := tsDevNode.setDstBindList(nil)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.DstBindList, _default_bind_list), ShouldBeZeroValue)
		})

		Convey("设置设备节点: 源绑定列表 - 列表均为空", func() {
			err := tsDevNode.setSrcBindList(spare)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.SrcBindList, _default_bind_list), ShouldBeZeroValue)
		})

		Convey("设置设备节点: 源绑定列表 - 列表均有值", func() {
			err := tsDevNode.setSrcBindList(exp)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.SrcBindList, expStr), ShouldBeZeroValue)
		})

		Convey("设置设备节点: 绑定源列表 - 列表为nil", func() {
			err := tsDevNode.setSrcBindList(nil)
			So(err, ShouldBeNil)
			So(strings.Compare(tsDevNode.SrcBindList, _default_bind_list), ShouldBeZeroValue)
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
		ProductId: 3000,
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
			So(dev.GetProductID(), ShouldEqual, 3000)
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
	ProductId: 3000,
}

var tsDev0 = &ZbDeviceInfo{
	IeeeAddr:  55667788,
	NwkAddr:   1122,
	Capacity:  1,
	ProductId: 3000,
}

func TestDevice(t *testing.T) {
	Convey("设备表和设备节点表", t, func() {
		Convey("创建设备和所有的节点", func() {
			err := UpdateZbDeviceAndANode(11223344, 5566, 2, 3000)
			So(err, ShouldBeNil)

			err = UpdateZbDeviceAndANode(55667788, 1122, 1, 3000)
			So(err, ShouldBeNil)

			err = UpdateZbDeviceAndANode(11111111, 1111, 2, 3000)
			So(err, ShouldBeNil)

			err = UpdateZbDeviceAndANode(22222222, 2222, 1, 3000)
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

			o2, err := LookupZbDeviceNodeByID(o1.ID)

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
			testDev.updateZbDeviceAndNodeNwkAddr(5566)
		})

		Convey("绑定两个互补的设备节点", func() {
			err := BindZbDeviceNode(11223344, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(11223344, 2, 11223344, 3, 3)
			So(err, ShouldBeNil)
			err = BindZbDeviceNode(55667788, 2, 55667788, 3, 3)
			So(err, ShouldBeNil)
		})

		Convey("根据网络址,节点号查找绑定表的所有设备节点", func() {
			devnodes, err := BindFindZbDeviceNodeByNN(5566, 2, 3)
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
