package ltl

//import (
//	. "github.com/smartystreets/goconvey/convey"
//)

/*
//测试0
var HdrOrg = &FrameHdr{
	TrunkID:     0x1234,
	NodeNo:      0x55,
	TransSeqNum: 0x01,
	CommandID:   LTL_CMD_CONFIGURE_REPORTING, // 0x06
	FrameHdrCtl_t: FrameHdrCtl{
		Type_l:            LTL_FRAMECTL_TYPE_PROFILE,
		DisableDefaultRsp: LTL_FRAMECTL_DIS_DEFAULT_RSP_OFF, // false
	},
}
var HdrOrgSlice = []byte{0x34, 0x12, 0x55, 0x01, 0x00, 0x06}

//测试1
var HdrOrg1 = &FrameHdr{
	TrunkID:     0x5678,
	NodeNo:      0xaa,
	TransSeqNum: 0x79,
	CommandID:   LTL_CMD_READ_ATTRIBUTES, //0x00
	FrameHdrCtl_t: FrameHdrCtl{
		Type_l:            LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC,
		DisableDefaultRsp: LTL_FRAMECTL_DIS_DEFAULT_RSP_ON, // true
	},
}
var HdrOrgSlice1 = []byte{0x78, 0x56, 0xaa, 0x79, 0x05, 0x00}

func TestFrameIsProfileCmd(t *testing.T) {
	Convey("frame control type field - profile", t, func() {
		So(IsProfileCmd(LTL_FRAMECTL_TYPE_PROFILE), ShouldBeTrue)
		So(IsProfileCmd(LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC), ShouldBeFalse)
	})
}

func TestIsTrunkCmd(t *testing.T) {
	Convey("frame control type field - trunk specific", t, func() {
		So(IsTrunkCmd(LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC), ShouldBeTrue)
		So(IsTrunkCmd(LTL_FRAMECTL_TYPE_PROFILE), ShouldBeFalse)
	})
}

func TestHdrSize(t *testing.T) {
	Convey("frame head size", t, func() {
		So(hdrSize(), ShouldEqual, 6)
	})
}
func TestEncodeHdr(t *testing.T) {
	Convey("frame head encode", t, func() {
		hdrExp := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, false, HdrOrg.CommandID)
		So(reflect.DeepEqual(hdrExp, HdrOrg), ShouldBeTrue)

		hdrExp1 := encodeHdr(HdrOrg1.TrunkID, HdrOrg1.NodeNo, HdrOrg1.TransSeqNum, LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC, true, HdrOrg1.CommandID)
		So(reflect.DeepEqual(hdrExp1, HdrOrg1), ShouldBeTrue)
	})

}

func TestBuildHdr(t *testing.T) {
	Convey("frame head build", t, func() {
		So(bytes.Equal(HdrOrg.buildHdr(), HdrOrgSlice), ShouldBeTrue)
		So(bytes.Equal(HdrOrg1.buildHdr(), HdrOrgSlice1), ShouldBeTrue)
	})
}
func TestParseHdr(t *testing.T) {
	Convey("frame head parse", t, func() {
		hdr, reamin := parseHdr(HdrOrgSlice)
		So(reflect.DeepEqual(hdr, HdrOrg), ShouldBeTrue)
		So(len(reamin), ShouldBeZeroValue)

		hdr1, reamin1 := parseHdr(HdrOrgSlice1)
		So(reflect.DeepEqual(hdr1, HdrOrg1), ShouldBeTrue)
		So(len(reamin1), ShouldBeZeroValue)
	})
}
func TestIsValidDataType(t *testing.T) {
	Convey("是否有效数据类型", t, func() {
		Convey("是否有效数据类型 - 不是", func() {
			So(isAnalogDataType(LTL_DATATYPE_UNKNOWN), ShouldBeFalse)
		})

		Convey("是否有效数据类型 - 是", func() {
			So(isValidDataType(LTL_DATATYPE_CHAR_STR), ShouldBeTrue)
			So(isValidDataType(LTL_DATATYPE_UINT16_ARRAY), ShouldBeTrue)
		})
	})
}

func TestIsBaseDataType(t *testing.T) {
	Convey("是否基本数据类型", t, func() {
		Convey("是否基本数据类型 - 不是", func() {
			So(isBaseDataType(LTL_DATATYPE_CHAR_STR), ShouldBeFalse)
			So(isBaseDataType(LTL_DATATYPE_UINT16_ARRAY), ShouldBeFalse)
		})

		Convey("是否基本数据类型 - 是", func() {
			So(isBaseDataType(LTL_DATATYPE_UINT16), ShouldBeTrue)
			So(isBaseDataType(LTL_DATATYPE_BOOLEAN), ShouldBeTrue)
		})
	})
}

func TestIsComplexDataType(t *testing.T) {
	Convey("是否复杂数据类型", t, func() {
		Convey("是否复杂数据类型 - 不是", func() {
			So(isComplexDataType(LTL_DATATYPE_BOOLEAN), ShouldBeFalse)
			So(isComplexDataType(LTL_DATATYPE_UINT16), ShouldBeFalse)
		})

		Convey("是否复杂数据类型 - 是", func() {
			So(isComplexDataType(LTL_DATATYPE_UINT16_ARRAY), ShouldBeTrue)
		})
	})
}

func TestIsAnalogDataType(t *testing.T) {
	Convey("判否模拟数据类型", t, func() {
		Convey("判否模拟数据类型 - 不是", func() {
			So(isAnalogDataType(LTL_DATATYPE_BOOLEAN), ShouldBeFalse)
			So(isAnalogDataType(LTL_DATATYPE_UINT16_ARRAY), ShouldBeFalse)
		})

		Convey("判否模拟数据类型 - 是", func() {
			So(isAnalogDataType(LTL_DATATYPE_UINT16), ShouldBeTrue)
		})
	})
}

func TestGetBaseDataTypeLength(t *testing.T) {
	Convey("获得基本数据类型长度", t, func() {
		So(getBaseDataTypeLength(LTL_DATATYPE_BOOLEAN), ShouldEqual, 1)
		So(getBaseDataTypeLength(LTL_DATATYPE_UINT16), ShouldEqual, 2)
		So(getBaseDataTypeLength(LTL_DATATYPE_UINT32), ShouldEqual, 4)
		So(getBaseDataTypeLength(LTL_DATATYPE_UINT64), ShouldEqual, 8)
		So(getBaseDataTypeLength(LTL_DATATYPE_UNKNOWN), ShouldEqual, 0)
	})
}

func TestMarshal(t *testing.T) {
	Convey("序列数据为bytes", t, func() {
		var actual []byte
		var dataType byte
		var err error

		Convey("序列数据为bytes - *int16", func() {
			var ti16 int16 = 103
			var pti16 *int16 = &ti16

			actual, dataType, err = marshal(pti16)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT16)
			So(bytes.Equal(actual, utils.Little_Endian.Putuint16(uint16(*pti16))), ShouldBeTrue)
		})
		Convey("序列数据为bytes - bool", func() {
			actual, dataType, err = marshal(true)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_BOOLEAN)
			So(bytes.Equal(actual, []byte{1}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - int8", func() {
			actual, dataType, err = marshal(int8(8))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT8)
			So(bytes.Equal(actual, []byte{8}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - uint8", func() {
			actual, dataType, err = marshal(uint8(8))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT8)
			So(bytes.Equal(actual, []byte{8}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - int16", func() {
			actual, dataType, err = marshal(int16(1234))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT16)
			So(bytes.Equal(actual, []byte{0xd2, 0x04}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - uint16", func() {
			actual, dataType, err = marshal(uint16(1234))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT16)
			So(bytes.Equal(actual, []byte{0xd2, 0x04}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - int32", func() {
			actual, dataType, err = marshal(int32(12345678))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT32)
			So(bytes.Equal(actual, []byte{0x4e, 0x61, 0xbc, 0x00}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - uint32", func() {
			actual, dataType, err = marshal(uint32(12345678))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT32)
			So(bytes.Equal(actual, []byte{0x4e, 0x61, 0xbc, 0x00}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - int64", func() {
			actual, dataType, err = marshal(int64(1234567812345678))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT64)
			So(bytes.Equal(actual, []byte{0x4e, 0xef, 0xe7, 0x37, 0xd5, 0x62, 0x04, 0x00}), ShouldBeTrue)
		})
		Convey("序列数据为bytes - uint64", func() {
			actual, dataType, err = marshal(uint64(1234567812345678))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT64)
			So(bytes.Equal(actual, []byte{0x4e, 0xef, 0xe7, 0x37, 0xd5, 0x62, 0x04, 0x00}), ShouldBeTrue)
		})

		Convey("序列数据为bytes - float32", func() {
			actual, dataType, err = marshal(float32(3.3702805504e+12))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_SINGLE_PREC)
			So(bytes.Equal(actual, []byte{0x18, 0x2d, 0x44, 0x54}), ShouldBeTrue)
		})

		Convey("序列数据为bytes - float64", func() {
			actual, dataType, err = marshal(float64(3.141592653589793))
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_DOUBLE_PREC)
			So(bytes.Equal(actual, []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}), ShouldBeTrue)
		})

		Convey("序列数据为bytes - string, *string", func() {
			var str string = "hello world"
			var pstr *string = &str
			expect := []byte{}
			expect = append(expect, byte(len(*pstr)))
			expect = append(expect, ([]byte(*pstr))...)
			actual, dataType, err = marshal(pstr)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_CHAR_STR)
			So(bytes.Equal(actual, expect), ShouldBeTrue)

			expect = []byte{}
			expect = append(expect, byte(len(str)))
			expect = append(expect, ([]byte(str))...)
			actual, dataType, err = marshal(str)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_CHAR_STR)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []int8", func() {
			array := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
			expect := []byte{}

			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)

			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT8_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []uint8", func() {
			array := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
			expect := []byte{}

			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)

			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT8_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []int16", func() {
			array := []int16{1, 2, 3, 4, 5, 6}
			expect := []byte{}
			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)
			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT16_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []uint16", func() {
			array := []uint16{1, 2, 3, 4, 5, 6}
			expect := []byte{}
			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)
			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT16_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []int32", func() {
			array := []int32{1, 2, 3, 4, 5, 6}
			expect := []byte{}
			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)
			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT32_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})

		Convey("序列数据为bytes - []uint32", func() {
			array := []uint32{1, 2, 3, 4, 5, 6}
			expect := []byte{}
			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)
			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT32_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []int64", func() {
			array := []int64{1, 2, 3, 4, 5, 6}
			expect := []byte{}
			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)
			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_INT64_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})
		Convey("序列数据为bytes - []uint64", func() {
			array := []uint64{1, 2, 3, 4, 5, 6}
			expect := []byte{}
			expect = append(expect, byte(len(array)))
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, array)
			expect = append(expect, buf.Bytes()...)
			actual, dataType, err = marshal(array)
			So(err, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UINT64_ARRAY)
			So(bytes.Equal(actual, expect), ShouldBeTrue)
		})

		Convey("序列数据为bytes - 未知类型", func() {
			actual, dataType, err = marshal(int(1))
			So(actual, ShouldBeNil)
			So(dataType, ShouldEqual, LTL_DATATYPE_UNKNOWN)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUnMarshal(t *testing.T) {
	Convey("bytes反序列化为值", t, func() {
		Convey("bytes反序列化为值 - bool", func() {
			v, err := unMarshal(LTL_DATATYPE_BOOLEAN, []byte{0})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Bool(), ShouldBeFalse)

			v, err = unMarshal(LTL_DATATYPE_BOOLEAN, []byte{1})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Bool(), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - int8", func() {
			v, err := unMarshal(LTL_DATATYPE_INT8, []byte{0x80})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Int(), ShouldEqual, math.MinInt8)
		})

		Convey("bytes反序列化为值 - int16", func() {
			v, err := unMarshal(LTL_DATATYPE_INT16, []byte{0x00, 0x80})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Int(), ShouldEqual, math.MinInt16)
		})

		Convey("bytes反序列化为值 - int32", func() {
			v, err := unMarshal(LTL_DATATYPE_INT32, []byte{0x00, 0x00, 0x00, 0x80})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Int(), ShouldEqual, math.MinInt32)
		})

		Convey("bytes反序列化为值 - int64", func() {
			v, err := unMarshal(LTL_DATATYPE_INT64, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Int(), ShouldEqual, math.MinInt64)
		})

		Convey("bytes反序列化为值 - uint8", func() {
			v, err := unMarshal(LTL_DATATYPE_UINT8, []byte{55})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Uint(), ShouldEqual, 55)
		})

		Convey("bytes反序列化为值 - uint16", func() {
			v, err := unMarshal(LTL_DATATYPE_UINT16, []byte{0x34, 0x12})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Uint(), ShouldEqual, 0x1234)
		})

		Convey("bytes反序列化为值 - uint32", func() {
			v, err := unMarshal(LTL_DATATYPE_UINT32, []byte{0x78, 0x56, 0x34, 0x12})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Uint(), ShouldEqual, 0x12345678)
		})

		Convey("bytes反序列化为值 - uint64", func() {
			v, err := unMarshal(LTL_DATATYPE_UINT64, []byte{0x78, 0x56, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Uint(), ShouldEqual, 0x1234567812345678)
		})

		Convey("bytes反序列化为值 - float32", func() {
			v, err := unMarshal(LTL_DATATYPE_SINGLE_PREC, []byte{0x18, 0x2d, 0x44, 0x54})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Float(), ShouldEqual, 3.3702805504e+12)
		})

		Convey("bytes反序列化为值 - float64", func() {
			v, err := unMarshal(LTL_DATATYPE_DOUBLE_PREC, []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40})
			So(err, ShouldBeNil)
			So(reflect.ValueOf(v).Float(), ShouldAlmostEqual, 3.141592653589793)
		})

		Convey("bytes反序列化为值 - string", func() {
			b := "hello world"
			v, err := unMarshal(LTL_DATATYPE_CHAR_STR, []byte(b))
			So(err, ShouldBeNil)
			So(strings.EqualFold(reflect.ValueOf(v).String(), b), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []int8", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x77}
			expect := []int8{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x77}
			v, err := unMarshal(LTL_DATATYPE_INT8_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []int16", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40}
			expect := []int16{0x2d18, 0x5444, 0x217b, 0x4009, 0x2d18, 0x5444, 0x217b, 0x4009}
			v, err := unMarshal(LTL_DATATYPE_INT16_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []int32", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40}
			expect := []int32{0x54442d18, 0x4009217b, 0x54442d18, 0x4009217b}
			v, err := unMarshal(LTL_DATATYPE_INT32_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []int64", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40}
			expect := []int64{0x4009217b54442d18, 0x4009217b54442d18}
			v, err := unMarshal(LTL_DATATYPE_INT64_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []uint8", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x77}
			expect := []uint8{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x77}
			v, err := unMarshal(LTL_DATATYPE_UINT8_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []uint16", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40}
			expect := []uint16{0x2d18, 0x5444, 0x217b, 0x4009, 0x2d18, 0x5444, 0x217b, 0x4009}
			v, err := unMarshal(LTL_DATATYPE_UINT16_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []uint32", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40}
			expect := []uint32{0x54442d18, 0x4009217b, 0x54442d18, 0x4009217b}
			v, err := unMarshal(LTL_DATATYPE_UINT32_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - []uint64", func() {
			b := []byte{0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40, 0x18, 0x2d, 0x44, 0x54, 0x7b, 0x21, 0x09, 0x40}
			expect := []uint64{0x4009217b54442d18, 0x4009217b54442d18}
			v, err := unMarshal(LTL_DATATYPE_UINT64_ARRAY, b)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(v, expect), ShouldBeTrue)
		})

		Convey("bytes反序列化为值 - 未知类型", func() {
			b := []byte{0x18, 0x2d}
			v, err := unMarshal(LTL_DATATYPE_UNKNOWN, b)
			So(err, ShouldNotBeNil)
			So(v, ShouldBeNil)
		})
	})
}

func TestStr2AppString(t *testing.T) {
	Convey("字符串转LTL应用字符串", t, func() {
		org := "lchtime"
		expect := []byte{7, 'l', 'c', 'h', 't', 'i', 'm', 'e'}

		Convey("字符串转LTL应用字符串 -- 空间足够", func() {
			So(bytes.Equal(Str2AppString(org, 10), expect), ShouldBeTrue)

		})
		Convey("字符串转LTL应用字符串 -- 空间限制", func() {
			expect[0] = 6 - OCTET_CHAR_HEADROOM_LEN // 只能存5个字符
			So(bytes.Equal(Str2AppString(org, 6), expect[0:6]), ShouldBeTrue)
		})
	})
}

type lout struct {
	cc            chan struct{}
	inCommingMsg  chan *IncomingMsgPkt
	actual        []byte
	actualAddress uint16
}

func (this *lout) WriteMsg(DstAddr uint16, p []byte) error {
	this.actual = p              // 回传
	this.actualAddress = DstAddr // 回传测试

	return nil
}

func (this *lout) CloseChan() <-chan struct{} {
	return this.cc
}

func (this *lout) IncommingMsg() <-chan *IncomingMsgPkt {
	return this.inCommingMsg
}

var dstAddress uint16 = 0x1122

func TestSendCommand(t *testing.T) {
	Convey("sendcommand 发送命令", t, func() {
		cmdFormat := []byte{0x55, 0x66, 0x77}
		inslout := &lout{}

		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}
		instance.Start()

		hdr := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, false, HdrOrg.CommandID)
		expect := hdr.buildHdr()

		Convey("sendcommand 无数据域", func() {
			err := instance.SendCommand(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, false, HdrOrg.CommandID, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		Convey("sendcommand 有数据域", func() {
			expect = append(expect, cmdFormat...) // 追加数据域
			err := instance.SendCommand(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, false, HdrOrg.CommandID, cmdFormat)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})
	})
}

func TestSendSpecificCmd(t *testing.T) {
	Convey("SendSpecificCmd 集下特殊命令请求", t, func() {
		cmdFormat := []byte{0x55, 0x66, 0x77}
		inslout := &lout{}
		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}

		hdr := encodeHdr(HdrOrg1.TrunkID, HdrOrg1.NodeNo, HdrOrg1.TransSeqNum, LTL_FRAMECTL_TYPE_TRUNK_SPECIFIC, true, HdrOrg1.CommandID)
		expect := hdr.buildHdr()

		Convey("SendSpecificCmd 无数据域", func() {
			err := instance.SendSpecificCmd(dstAddress, HdrOrg1.TrunkID, HdrOrg1.NodeNo, HdrOrg1.TransSeqNum, true, HdrOrg1.CommandID, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		Convey("SendSpecificCmd 有数据域", func() {
			expect = append(expect, cmdFormat...)
			err := instance.SendSpecificCmd(dstAddress, HdrOrg1.TrunkID, HdrOrg1.NodeNo, HdrOrg1.TransSeqNum, true, HdrOrg1.CommandID, cmdFormat)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})
	})
}

func TestSendReadReq(t *testing.T) {
	Convey("SendReadReq 读属性命令请求", t, func() {
		arrtid := []uint16{0x0001, 0x0102, 0xaabb}
		inslout := &lout{}
		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}

		HdrOrg.CommandID = LTL_CMD_READ_ATTRIBUTES // 0x00
		HdrOrg.Type_l = LTL_FRAMECTL_TYPE_PROFILE

		hdr := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, HdrOrg.CommandID)
		expect := hdr.buildHdr()

		Convey("SendReadReq 无数据域", func() {
			err := instance.SendReadReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		Convey("SendReadReq 有数据域", func() {
			buf := &bytes.Buffer{}
			binary.Write(buf, binary.LittleEndian, arrtid)
			expect = append(expect, buf.Bytes()...)
			err := instance.SendReadReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, arrtid)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

	})
}

func TestSendWriteReq(t *testing.T) {
	Convey("写属性命令请求", t, func() {

		inslout := &lout{}
		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}

		expect0 := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, LTL_CMD_WRITE_ATTRIBUTES).buildHdr()

		Convey("SendWriteReq 无数据域", func() {
			err := instance.SendWriteReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect0, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		expect1 := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, LTL_CMD_WRITE_ATTRIBUTES_UNDIVIDED).buildHdr()
		Convey("SendWriteReqUndivided 无数据域", func() {
			err := instance.SendWriteReqUndivided(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect1, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		expect2 := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, LTL_CMD_WRITE_ATTRIBUTES_NORSP).buildHdr()
		Convey("SendWriteReqNoRsp 无数据域", func() {
			err := instance.SendWriteReqNoRsp(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect2, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		wRec := []WriteRec{
			{
				AttrID:    1,
				AttriData: uint32(0x12345678),
			},
			{
				AttrID:    2,
				AttriData: uint8(0xaa),
			},
		}
		for _, ws := range wRec {
			expect0 = append(expect0, utils.Little_Endian.Putuint16(ws.AttrID)...)
			expect1 = append(expect1, utils.Little_Endian.Putuint16(ws.AttrID)...)
			expect2 = append(expect2, utils.Little_Endian.Putuint16(ws.AttrID)...)
			bs, dataType, err := marshal(ws.AttriData)
			if err != nil {
				panic(err)
			}
			expect0 = append(expect0, dataType)
			expect1 = append(expect1, dataType)
			expect2 = append(expect2, dataType)
			expect0 = append(expect0, bs...)
			expect1 = append(expect1, bs...)
			expect2 = append(expect2, bs...)
		}

		Convey("SendWriteReq 有数据域", func() {
			err := instance.SendWriteReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, wRec)

			So(err, ShouldBeNil)
			So(bytes.Equal(expect0, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		Convey("SendWriteReqUndivided 有数据域", func() {
			err := instance.SendWriteReqUndivided(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, wRec)

			So(err, ShouldBeNil)
			So(bytes.Equal(expect1, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		Convey("SendWriteReqNoRsp 有数据域", func() {
			err := instance.SendWriteReqNoRsp(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, wRec)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect2, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})
	})
}

func TestSendConfigReportReq(t *testing.T) {
	Convey("配置报告属性", t, func() {

		inslout := &lout{}
		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}

		expect := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, LTL_CMD_CONFIGURE_REPORTING).buildHdr()

		Convey("SendConfigReportReq 无数据域", func() {
			err := instance.SendWriteReportCfgReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		wCfg := []WriteReportCfgRec{
			{
				AttrID:           1,
				MinReportInt:     1000,
				ReportableChange: uint32(0x12345678),
			},
			{
				AttrID:           2,
				MinReportInt:     2000,
				ReportableChange: true,
			},
		}

		for _, ws := range wCfg {
			expect = append(expect, utils.Little_Endian.Putuint16(ws.AttrID)...)
			bs, dataType, err := marshal(ws.ReportableChange)
			if err != nil || !isBaseDataType(dataType) {
				panic(err)
			}
			expect = append(expect, dataType)
			expect = append(expect, utils.Little_Endian.Putuint16(ws.MinReportInt)...)
			if isAnalogDataType(dataType) {
				expect = append(expect, bs...)
			}
		}

		Convey("SendWriteReq 有数据域", func() {
			err := instance.SendWriteReportCfgReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, wCfg)

			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

	})
}

func TestSendReadReportCfgReq(t *testing.T) {
	Convey("读报告属性命令请求", t, func() {
		inslout := &lout{}
		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}

		expect := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, LTL_CMD_READ_CONFIGURE_REPORTING).buildHdr()
		Convey("SendReadReportCfgReq 无数据域", func() {
			err := instance.SendReadReportCfgReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, nil)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})

		Convey("SendReadReportCfgReq 有数据域", func() {
			arrtid := []uint16{0x0001, 0x0102, 0xaabb}
			buf := &bytes.Buffer{}
			binary.Write(buf, binary.LittleEndian, arrtid)
			expect = append(expect, buf.Bytes()...)

			err := instance.SendReadReportCfgReq(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, arrtid)
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})
	})
}

func TestSendDefaultRspCmd(t *testing.T) {
	Convey("默认应答命令", t, func() {
		inslout := &lout{}
		instance := &Ltl_t{
			WriteCloseMsgComming: inslout,
		}

		expect := encodeHdr(HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, LTL_FRAMECTL_TYPE_PROFILE, true, LTL_CMD_DEFAULT_RSP).buildHdr()
		expect = append(expect, 9, 10)
		Convey("SendDefaultRspCmd 有数据", func() {
			err := instance.SendDefaultRspCmd(dstAddress, HdrOrg.TrunkID, HdrOrg.NodeNo, HdrOrg.TransSeqNum, DefaultRsp{9, 10})
			So(err, ShouldBeNil)
			So(bytes.Equal(expect, inslout.actual), ShouldBeTrue)
			So(inslout.actualAddress == dstAddress, ShouldBeTrue)
		})
	})
}

// func TestParseInWriteRspCmd(t *testing.T) {
// 	dataOk := []byte{0}
// 	data := []byte{0x84, 0x01, 0x00, 0x86, 0x05, 0x00}
// 	expectOk := []WriteRspStatus_t{{Status: 0}}
// 	expect := []WriteRspStatus_t{{0x84, 0x0001}, {0x86, 0x0005}}

// 	inslout := &lout{}
// 	instance := NewInstance(inslout)

// 	Convey("解析写属性应答命令", t, func() {
// 		Convey("解析写属性应答命令 成功的回复", func() {
// 			s, err := instance.ParseInWriteRspCmd(dataOk)
// 			So(err, ShouldBeNil)
// 			So(reflect.DeepEqual(s, expectOk), ShouldBeTrue)
// 		})

// 		Convey("解析写属性应答命令 有失败的回复", func() {
// 			s, err := instance.ParseInWriteRspCmd(data)
// 			So(err, ShouldBeNil)
// 			So(reflect.DeepEqual(s, expect), ShouldBeTrue)
// 		})
// 	})
// }

// func TestParseInConfigReportRspCmd(t *testing.T) {
// 	dataOk := []byte{0}
// 	data := []byte{0x84, 0x01, 0x00, 0x86, 0x05, 0x00}
// 	expectOk := []CfgReportRspStatus_t{{Status: 0}}
// 	expect := []CfgReportRspStatus_t{{0x84, 0x0001}, {0x86, 0x0005}}

// 	inslout := &lout{}
// 	instance := NewInstance(inslout)

// 	Convey("解析配置报告应答命令", t, func() {
// 		Convey("解析配置报告应答命令 成功的回复", func() {
// 			s, err := instance.ParseInConfigReportRspCmd(dataOk)
// 			So(err, ShouldBeNil)
// 			So(reflect.DeepEqual(s, expectOk), ShouldBeTrue)
// 		})

// 		Convey("解析配置报告应答命令 有失败的回复", func() {
// 			s, err := instance.ParseInConfigReportRspCmd(data)
// 			So(err, ShouldBeNil)
// 			So(reflect.DeepEqual(s, expect), ShouldBeTrue)
// 		})
// 	})
// }

// func TestParseInDefaultRspCmd(t *testing.T) {
// 	Convey("解析默认回复应答命令", t, func() {
// 		Convey("解析默认回复应答命令 数据长度正确", func() {
// 			instance := NewInstance(&lout{})
// 			orgRsp := &DefaultRsp_t{
// 				CommandID:  1,
// 				StatusCode: 5,
// 			}
// 			testS := []byte{orgRsp.CommandID, orgRsp.StatusCode}

// 			atual, err := instance.ParseInDefaultRspCmd(testS)
// 			So(err, ShouldBeNil)
// 			So(reflect.DeepEqual(atual, orgRsp), ShouldBeTrue)
// 		})

// 		Convey("解析默认回复应答命令 数据长度不正确", func() {
// 			instance := NewInstance(&lout{})

// 			testS := []byte{1}

// 			atual, err := instance.ParseInDefaultRspCmd(testS)
// 			So(err, ShouldNotBeNil)
// 			So(atual, ShouldBeNil)
// 		})
// 	})
// }
*/
