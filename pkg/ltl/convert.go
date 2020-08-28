package ltl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
)

const (
	//errInvalidType  = "invalid value type"
	errAssertFailed = "type assertion failed"
)

// A Number represents a ltl number literal.
type Number string

// String returns the literal text of the number.
func (n Number) String() string {
	return string(n)
}

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// 属性值序列化成切片,并得到类型
func marshal(attrData interface{}) ([]byte, byte, error) {
	var length int
	var tmpBytes []byte
	var dataType byte = LTL_DATATYPE_UNKNOWN

	switch attrData.(type) {
	case *bool, bool:
		dataType = LTL_DATATYPE_BOOLEAN
	case *int8, int8:
		dataType = LTL_DATATYPE_INT8
	case *uint8, uint8:
		dataType = LTL_DATATYPE_UINT8
	case *int16, int16:
		dataType = LTL_DATATYPE_INT16
	case *uint16, uint16:
		dataType = LTL_DATATYPE_UINT16
	case *int32, int32:
		dataType = LTL_DATATYPE_INT32
	case *uint32, uint32:
		dataType = LTL_DATATYPE_UINT32
	case *int64, int64:
		dataType = LTL_DATATYPE_INT64
	case *uint64, uint64:
		dataType = LTL_DATATYPE_UINT64
	case *float32, float32:
		dataType = LTL_DATATYPE_SINGLE_PREC
	case *float64, float64:
		dataType = LTL_DATATYPE_DOUBLE_PREC
	case []int8:
		dataType = LTL_DATATYPE_INT8_ARRAY
	case []int16:
		dataType = LTL_DATATYPE_INT16_ARRAY
	case []int32:
		dataType = LTL_DATATYPE_INT32_ARRAY
	case []int64:
		dataType = LTL_DATATYPE_INT64_ARRAY
	case []uint8:
		dataType = LTL_DATATYPE_UINT8_ARRAY
	case []uint16:
		dataType = LTL_DATATYPE_UINT16_ARRAY
	case []uint32:
		dataType = LTL_DATATYPE_UINT32_ARRAY
	case []uint64:
		dataType = LTL_DATATYPE_UINT64_ARRAY
	case *string, string: // NOTE:  理论上不应该出现*string,这里还是实现了
		dataType = LTL_DATATYPE_CHAR_STR
		str, ok := attrData.(string)
		if !ok {
			tmpstr := attrData.(*string)
			str = *(tmpstr)
		}
		length = len(str)

		p := make([]byte, 0, length+OCTET_CHAR_HEADROOM_LEN)
		p = append(p, byte(length))
		return append(p, ([]byte(str))...), dataType, nil
	}

	if dataType == LTL_DATATYPE_UNKNOWN {
		return nil, LTL_DATATYPE_UNKNOWN, ErrInvalidDataType
	}

	if isBaseDataType(dataType) {
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, attrData); err != nil {
			return nil, dataType, err
		}

		return buf.Bytes(), dataType, nil
	}

	out := new(bytes.Buffer)
	if err := binary.Write(out, binary.LittleEndian, attrData); err != nil {
		return nil, dataType, err
	}
	length = reflect.ValueOf(attrData).Len()
	tmpBytes = out.Bytes()

	p := make([]byte, 0, length+OCTET_CHAR_HEADROOM_LEN)
	p = append(p, byte(length))

	return append(p, tmpBytes...), dataType, nil
}

func isValidDataType(dataType byte) bool {
	return dataType < LTL_DATATYPE_UNKNOWN
}

// 是否为LTL基本数据类型
func isBaseDataType(dataType byte) bool {
	return dataType <= LTL_DATATYPE_DOUBLE_PREC
	// switch dataType {
	// case LTL_DATATYPE_BOOLEAN, LTL_DATATYPE_INT8, LTL_DATATYPE_INT16, LTL_DATATYPE_INT32, LTL_DATATYPE_INT64,
	// 	LTL_DATATYPE_UINT8, LTL_DATATYPE_UINT16, LTL_DATATYPE_UINT32, LTL_DATATYPE_UINT64,
	// 	LTL_DATATYPE_SINGLE_PwrRec, LTL_DATATYPE_DOUBLE_PwrRec:
	// 	return true
	// }
}

// 是否为LTL复杂数据类型,主要为ltl的字符串,字符数组,双字节数据, 最大不超过255
func isComplexDataType(dataType byte) bool {
	if dataType >= LTL_DATATYPE_CHAR_STR && dataType < LTL_DATATYPE_UNKNOWN {
		return true
	}
	// switch dataType {
	// case LTL_DATATYPE_CHAR_STR,
	// 	LTL_DATATYPE_INT8_ARRAY, LTL_DATATYPE_INT16_ARRAY, LTL_DATATYPE_INT32_ARRAY, LTL_DATATYPE_INT64_ARRAY,
	// 	LTL_DATATYPE_UINT8_ARRAY, LTL_DATATYPE_UINT16_ARRAY, LTL_DATATYPE_UINT32_ARRAY, LTL_DATATYPE_UINT64_ARRAY:
	// 	return true
	// }

	return false
}

// 是否LTL模拟数据类型
func isAnalogDataType(dataType byte) bool {
	if isBaseDataType(dataType) && dataType != LTL_DATATYPE_BOOLEAN {
		return true
	}

	// switch dataType {
	// case  LTL_DATATYPE_INT8, LTL_DATATYPE_INT16, LTL_DATATYPE_INT32, LTL_DATATYPE_INT64,
	// 	LTL_DATATYPE_UINT8, LTL_DATATYPE_UINT16, LTL_DATATYPE_UINT32, LTL_DATATYPE_UINT64,
	// 	LTL_DATATYPE_SINGLE_PwrRec, LTL_DATATYPE_DOUBLE_PwrRec:
	// 	return true
	// }

	return false
}

// 获取基本数据类型的长度
func getBaseDataTypeLength(dataType byte) int {
	var lens int = 0
	switch dataType {
	case LTL_DATATYPE_BOOLEAN, LTL_DATATYPE_UINT8, LTL_DATATYPE_INT8:
		lens = 1
	case LTL_DATATYPE_UINT16, LTL_DATATYPE_INT16:
		lens = 2
	case LTL_DATATYPE_UINT32, LTL_DATATYPE_INT32, LTL_DATATYPE_SINGLE_PREC:
		lens = 4
	case LTL_DATATYPE_UINT64, LTL_DATATYPE_INT64, LTL_DATATYPE_DOUBLE_PREC:
		lens = 8
	}
	// LTL_DATATYPE_UNKNOWN
	return lens
}

// 反序列化成值
func (this AttrValues) unMarshal() (interface{}, error) {
	switch this.DataType {
	case LTL_DATATYPE_BOOLEAN:
		return this.Bool()
	case LTL_DATATYPE_UINT8:
		return this.Uint8()
	case LTL_DATATYPE_INT8:
		v, err := this.Uint8()
		return int8(v), err

	case LTL_DATATYPE_UINT16:
		return this.Uint16()
	case LTL_DATATYPE_INT16:
		v, err := this.Uint8()
		return int16(v), err

	case LTL_DATATYPE_UINT32:
		return this.Uint32()
	case LTL_DATATYPE_INT32:
		v, err := this.Uint32()
		return int32(v), err

	case LTL_DATATYPE_UINT64:
		return this.Uint64()
	case LTL_DATATYPE_INT64:
		v, err := this.Uint64()
		return int64(v), err

	case LTL_DATATYPE_SINGLE_PREC:
		v, err := this.Float64()
		return float32(v), err

	case LTL_DATATYPE_DOUBLE_PREC:
		return this.Float64()

	case LTL_DATATYPE_CHAR_STR:
		return this.String()

	case LTL_DATATYPE_INT8_ARRAY:
		return this.ArrayInt8()
	case LTL_DATATYPE_UINT8_ARRAY:
		return this.ArrayUint8()

	case LTL_DATATYPE_INT16_ARRAY:
		return this.ArrayInt16()
	case LTL_DATATYPE_UINT16_ARRAY:
		return this.ArrayUint16()

	case LTL_DATATYPE_INT32_ARRAY:
		return this.ArrayInt32()
	case LTL_DATATYPE_UINT32_ARRAY:
		return this.ArrayUint32()

	case LTL_DATATYPE_INT64_ARRAY:
		return this.ArrayInt64()
	case LTL_DATATYPE_UINT64_ARRAY:
		return this.ArrayUint64()
	}
	//  LTL_DATATYPE_UNKNOWN
	return nil, ErrInvalidDataType
}

// 字符串转LTL应用字符串
func Str2AppString(RawStr string, ApplenMax byte) []byte {
	rawlen := len(RawStr)
	if rawlen > int(ApplenMax-OCTET_CHAR_HEADROOM_LEN) {
		rawlen = int(ApplenMax - OCTET_CHAR_HEADROOM_LEN)
	}
	appstrbyte := make([]byte, 0, rawlen+OCTET_CHAR_HEADROOM_LEN)
	appstrbyte = append(appstrbyte, byte(rawlen))
	rawbyte := []byte(RawStr)
	return append(appstrbyte, rawbyte[0:rawlen]...)
}

func (this AttrValues) Bool() (bool, error) {
	if this.DataType != LTL_DATATYPE_BOOLEAN &&
		len(this.Data) == 0 {
		return false, errors.New(errAssertFailed)
	}

	if this.Data[0] > 0 {
		return true, nil
	}

	return false, nil
}

func (this AttrValues) Uint8() (uint8, error) {
	if this.DataType != LTL_DATATYPE_UINT8 &&
		this.DataType != LTL_DATATYPE_INT8 &&
		len(this.Data) == 0 {
		return 0, errors.New(errAssertFailed)
	}

	return uint8(this.Data[0]), nil
}

func (this AttrValues) Uint16() (uint16, error) {
	if this.DataType != LTL_DATATYPE_UINT16 &&
		this.DataType != LTL_DATATYPE_INT16 &&
		len(this.Data) < 2 {
		return 0, errors.New(errAssertFailed)
	}

	return binary.LittleEndian.Uint16(this.Data), nil
}

func (this AttrValues) Uint32() (uint32, error) {
	if this.DataType != LTL_DATATYPE_UINT32 &&
		this.DataType != LTL_DATATYPE_INT32 &&
		len(this.Data) < 4 {
		return 0, errors.New(errAssertFailed)
	}

	return binary.LittleEndian.Uint32(this.Data), nil
}

func (this AttrValues) Uint64() (uint64, error) {
	if this.DataType != LTL_DATATYPE_UINT64 &&
		this.DataType != LTL_DATATYPE_INT64 &&
		len(this.Data) < 8 {
		return 0, errors.New(errAssertFailed)
	}

	return binary.LittleEndian.Uint64(this.Data), nil
}

func (this AttrValues) Float64() (float64, error) {
	switch this.DataType {
	case LTL_DATATYPE_SINGLE_PREC:
		if len(this.Data) < 4 {
			break
		}
		var v float32
		buf := bytes.NewReader(this.Data)
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return 0, err
		}
		return float64(v), nil
	case LTL_DATATYPE_DOUBLE_PREC:
		if len(this.Data) < 8 {
			break
		}
		var v float64
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, &v); err != nil {
			return 0, err
		}
		return v, nil
	}

	return 0, errors.New(errAssertFailed)
}

func (this AttrValues) String() (string, error) {
	if this.DataType == LTL_DATATYPE_CHAR_STR {
		s := make([]byte, len(this.Data))
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return "", err
		}
		return string(s), nil
	}
	return "", errors.New(errAssertFailed)
}

func (this AttrValues) ArrayInt8() ([]int8, error) {
	if this.DataType == LTL_DATATYPE_INT8_ARRAY {
		s := make([]int8, len(this.Data))
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}

func (this AttrValues) ArrayUint8() ([]uint8, error) {
	if this.DataType == LTL_DATATYPE_UINT8_ARRAY {
		s := make([]uint8, len(this.Data))
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}

func (this AttrValues) ArrayInt16() ([]int16, error) {
	if this.DataType == LTL_DATATYPE_INT16_ARRAY {
		s := make([]int16, len(this.Data)/2)
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}
func (this AttrValues) ArrayUint16() ([]uint16, error) {
	if this.DataType == LTL_DATATYPE_UINT16_ARRAY {
		s := make([]uint16, len(this.Data)/2)
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}
func (this AttrValues) ArrayInt32() ([]int32, error) {
	if this.DataType == LTL_DATATYPE_INT32_ARRAY {
		s := make([]int32, len(this.Data)/4)
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}
func (this AttrValues) ArrayUint32() ([]uint32, error) {
	if this.DataType == LTL_DATATYPE_UINT32_ARRAY {
		s := make([]uint32, len(this.Data)/4)
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}
func (this AttrValues) ArrayInt64() ([]int64, error) {
	if this.DataType == LTL_DATATYPE_INT64_ARRAY {
		s := make([]int64, len(this.Data)/8)
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}
func (this AttrValues) ArrayUint64() ([]uint64, error) {
	if this.DataType == LTL_DATATYPE_UINT64_ARRAY {
		s := make([]uint64, len(this.Data)/8)
		if err := binary.Read(bytes.NewReader(this.Data), binary.LittleEndian, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, errors.New(errAssertFailed)
}

func (this AttrValues) MustBool() bool {
	b, _ := this.Bool()
	return b
}

func (this AttrValues) MustUint8() uint8 {
	b, _ := this.Uint8()
	return b
}

func (this AttrValues) MustUint16() uint16 {
	b, _ := this.Uint16()
	return b
}

func (this AttrValues) MustUint32() uint32 {
	b, _ := this.Uint32()
	return b
}

func (this AttrValues) MustUint64() uint64 {
	b, _ := this.Uint64()
	return b
}

func (this AttrValues) MustFloat64() float64 {
	b, _ := this.Float64()
	return b
}

func (this AttrValues) MustString() string {
	b, _ := this.String()
	return b
}

func (this AttrValues) MustArrayInt8() []int8 {
	b, err := this.ArrayInt8()
	if err != nil {
		return []int8{}
	}
	return b
}

func (this AttrValues) MustArrayUint8() []uint8 {
	b, err := this.ArrayUint8()
	if err != nil {
		return []uint8{}
	}
	return b
}

func (this AttrValues) MustArrayInt16() []int16 {
	b, err := this.ArrayInt16()
	if err != nil {
		return []int16{}
	}
	return b
}
func (this AttrValues) MustArrayUint16() []uint16 {
	b, err := this.ArrayUint16()
	if err != nil {
		return []uint16{}
	}
	return b
}
func (this AttrValues) MustArrayInt32() []int32 {
	b, err := this.ArrayInt32()
	if err != nil {
		return []int32{}
	}
	return b
}
func (this AttrValues) MustArrayUint32() []uint32 {
	b, err := this.ArrayUint32()
	if err != nil {
		return []uint32{}
	}
	return b
}
func (this AttrValues) MustArrayInt64() []int64 {
	b, err := this.ArrayInt64()
	if err != nil {
		return []int64{}
	}
	return b
}
func (this AttrValues) MustArrayUint64() []uint64 {
	b, err := this.ArrayUint64()
	if err != nil {
		return []uint64{}
	}
	return b
}
