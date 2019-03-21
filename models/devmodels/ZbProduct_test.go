package devmodels

/*
import (
	//"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var inTk1 = []uint16{1, 2}
var outTk1 = []uint16{3, 4}
var inTk2 = []uint16{3}
var outTk2 = []uint16{5, 6, 7, 8}
var spareTk = []uint16{}
var node0 = &NodeDsc{InTrunk: []uint16{}, OutTrunk: []uint16{}}
var node1 = &NodeDsc{InTrunk: inTk1, OutTrunk: outTk1}
var node2 = &NodeDsc{InTrunk: inTk2, OutTrunk: outTk2}
var js = `{"NodeDscList":[{"InTrunk":[],"OutTrunk":[]},{"InTrunk":[1,2],"OutTrunk":[3,4]},{"InTrunk":[3],"OutTrunk":[5,6,7,8]}]}`

func TestNodeDsc(t *testing.T) {
	Convey("节点描述输入输出集表", t, func() {
		expect := &NodeDsc{InTrunk: inTk1, OutTrunk: outTk1}
		actual := &NodeDsc{}
		Convey("设置节点描述-表有值时", func() {
			actual.SetTrunkID(inTk1, outTk1)
			So(reflect.DeepEqual(actual, expect), ShouldBeTrue)
		})

		Convey("设置节点描述-表无值时", func() {
			actual.SetTrunkID(nil, nil)
			So(reflect.DeepEqual(actual, &NodeDsc{
				InTrunk:  spareTk,
				OutTrunk: spareTk,
			}), ShouldBeTrue)
		})

		Convey("获得节点输入输出集表-表有值时", func() {
			iTK, oTK := expect.GetTrunkID()
			So(reflect.DeepEqual(iTK, inTk1), ShouldBeTrue)
			So(reflect.DeepEqual(oTK, outTk1), ShouldBeTrue)
		})
		Convey("获得节点输入输出集表-表无值时", func() {
			iTK, oTK := (&NodeDsc{InTrunk: spareTk, OutTrunk: spareTk}).GetTrunkID()
			So(reflect.DeepEqual(iTK, spareTk), ShouldBeTrue)
			So(reflect.DeepEqual(oTK, spareTk), ShouldBeTrue)
		})
	})
}

func TestProductNodeDscList(t *testing.T) {
	Convey("产品节点描述列表", t, func() {
		nodeDsc := make([]*NodeDsc, 0)
		nodeDsc = append(nodeDsc, node0, node1, node2)
		pdt := &Product{}

		Convey("设置节点描列表-有节点时", func() {
			err := pdt.SetNodeDscList(nodeDsc)
			So(err, ShouldBeNil)
			So(strings.Compare(js, pdt.NodeList), ShouldBeZeroValue)
		})

		Convey("设置节点描列表-无节点时", func() {
			err := pdt.SetNodeDscList(nil)
			So(err, ShouldBeNil)
			So(strings.Compare(pdt.NodeList, _default_product_node_list), ShouldBeZeroValue)
		})

		Convey("获取节点描列表-有节点描述", func() {
			pdt.NodeList = js
			actual, err := pdt.GetDeviceNodeDscList()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(actual, nodeDsc), ShouldBeTrue)
		})

		Convey("获取节点描列表-无节点描述", func() {
			pdt.NodeList = _default_product_node_list
			actual, err := pdt.GetDeviceNodeDscList()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(actual, []*NodeDsc{}), ShouldBeTrue)
		})

	})
}

func TestProduct(t *testing.T) {
	Convey("产品数据库增删改查", t, func() {
		pid := uint32(1000)
		nodeDsc := make([]*NodeDsc, 0)
		nodeDsc = append(nodeDsc, node0, node1, node2)

		Convey("增加产品", func() {
			err := AddProduct(pid, nodeDsc, "switch")
			So(err, ShouldBeNil)

			UpdateProductDescritption(2000, "")
			err = AddProduct(2000, nil, "新产品1")
			So(err, ShouldBeNil)

			pdt := &Product{
				ProductId: 3000,
				NodeList:  js,
			}
			So(pdt.AddProduct(), ShouldBeNil)
		})

		Convey("更新获取产品描述", func() {
			err := UpdateProductDescritption(2000, "新产品3")
			So(err, ShouldBeNil)

			ds, err := GetProductDescription(2000)
			So(err, ShouldBeNil)
			So(strings.Compare(ds, "新产品3"), ShouldBeZeroValue)
		})

		Convey("查找产品", func() {
			pdt, err := LookupProduct(pid)
			So(err, ShouldBeNil)
			So(pdt.ProductId, ShouldEqual, pid)

			pdt, err = LookupProduct(4000)
			So(err, ShouldNotBeNil)
			So(pdt, ShouldBeNil)
		})

		Convey("获得产品节点列表", func() {
			actual, err := LookupProductDeviceNodeDscList(pid)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(actual, nodeDsc), ShouldBeTrue)
		})

		Convey("删除产品", func() {
			err := DeleteProduct(pid)
			So(err, ShouldBeNil)
		})

	})
}
*/
