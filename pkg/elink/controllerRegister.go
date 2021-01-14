package elink

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

// 错误信息
var ErrRouterNotMatch = errors.New("router not match")
var ErrChannelNotMatch = errors.New("channel not match")

// ControllerInfo 通道控制器信息
type ControllerInfo struct {
	controllerType reflect.Type
	methods        map[string]string
	initialize     func() ControllerInterface
}

// RunMode 运行模式
func (this *ControllerInfo) RunMode(rawMethod string) string {
	rawMethod = strings.ToLower(rawMethod)
	if m, ok := this.methods[rawMethod]; ok { // 有配置自定方法,直接使用对应的方法函数名
		return m
	} else if m, ok := this.methods["*"]; ok && HasMethod(rawMethod) { // 有配置 * ,method必须为elink method
		return m
	} else if HasMethod(rawMethod) { // 否则使用elink method
		return rawMethod
	}
	return MethodUnknown
}

// ControllerRegister 控制器注册器资源路由
type ControllerRegister struct {
	routers sync.Map
}

// NewControllerRegister 新建一个控制注册器
func NewControllerRegister() *ControllerRegister {
	return &ControllerRegister{}
}

// AddToRouter 添加路由到资源
//第三个参数就是用来设置对应 method 到函数名，即自定义method到功能函数,定义如下
//-  *表示任意的elink method 都执行该功能函数,不包含其它的自定义的method
//-  使用 method:funcname 格式来展示
//-  多个不同的, 格式使用 ; 分割
//-  多个 method 对应同一个funcname,method之间通过 , 来分割
// 如果同时存在 * 和自定义elnik Method对应的功能函数，那么优先执行自定义功能函数
func (sf *ControllerRegister) AddToRouter(resource string, c ControllerInterface, mapMethods ...string) error {
	reflectVal := reflect.ValueOf(c)         // 这其实是个指针的值
	t := reflect.Indirect(reflectVal).Type() //通过Indirect获得指针指向的值,并获得类型
	methods := make(map[string]string)
	if len(mapMethods) > 0 {
		semi := strings.Split(mapMethods[0], ";")
		for _, v := range semi {
			colon := strings.Split(v, ":")
			if len(colon) != 2 {
				return errors.New("method mapping format is invalid")
			}
			comma := strings.Split(colon[0], ",")
			if ContainsUnknownMethod(comma) {
				return errors.New("forbid used 'unknown' method")
			}
			for _, m := range comma { // 方法列表
				if val := reflectVal.MethodByName(colon[1]); val.IsValid() { // 找到reflectVal是否持有方法,通过判断是否返回零值
					methods[strings.ToLower(m)] = colon[1] // 将method 和 调用函数名映射
				} else {
					return errors.New("'" + colon[1] + "' method doesn't exist in the controller " + t.Name())
				}
			}
		}
	}

	cInfo := &ControllerInfo{
		controllerType: t,
		methods:        methods,
	}
	cInfo.initialize = func() ControllerInterface {
		Valc := reflect.New(cInfo.controllerType) // 创建指向controllerType类型的零值指针

		execController, ok := Valc.Interface().(ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface") // 不可能发生
		}

		elemVal := reflect.ValueOf(c).Elem()
		elemType := reflect.TypeOf(c).Elem()
		execElem := reflect.ValueOf(execController).Elem()

		numOfFields := elemVal.NumField()
		for i := 0; i < numOfFields; i++ {
			fieldType := elemType.Field(i)
			elemField := execElem.FieldByName(fieldType.Name)
			if elemField.CanSet() {
				fieldVal := elemVal.Field(i)
				elemField.Set(fieldVal)
			}
		}
		return execController
	}
	return sf.addToRouter(resource, cInfo)
}

// MatchRouter 找到资源,匹配路由,返回资源匹配参数值 splat
func (sf *ControllerRegister) MatchRouter(resource string) (*ControllerInfo, QueryValues, error) {
	rs := SplitResource(resource) // 主要去除无效"."
	// 先匹配固定路由
	if prefix := strings.Join(rs, "."); len(prefix) > 0 { // 固定resource 匹配
		if v, ok := sf.routers.Load(prefix); ok && v.(*Root).fixInfo != nil {
			return v.(*Root).fixInfo, QueryValues{}, nil
		}
	}

	// 匹配前置prefix a.b.c.d  先匹配长的,往短匹配que
	for length := len(rs); length > 0; {
		s := rs[:length-1]
		if prefix := strings.Join(s, "."); len(prefix) > 0 { // 正则前置
			if v, ok := sf.routers.Load(prefix); ok && v.(*Root).regInfo != nil {
				return v.(*Root).regInfo, QueryValues{"splat": rs[length-1:]}, nil
			}
		}
		length = len(s)
	}

	return nil, nil, ErrRouterNotMatch
}

// DeleteRouter 删除资源,路由相关
func (sf *ControllerRegister) DeleteRouter(resource string) {
	var prefix string
	var isFix bool

	rs := SplitResource(resource)
	length := len(rs)
	if length == 0 {
		return
	}

	//固定路由还是正则路由
	if strings.EqualFold(rs[length-1], "@") {
		prefix = strings.Join(rs[:length-1], ".")
		isFix = false
	} else {
		prefix = strings.Join(rs, ".")
		isFix = true
	}

	if len(prefix) == 0 {
		return
	}

	if v, ok := sf.routers.Load(prefix); ok {
		root := v.(*Root)
		if isFix {
			root.fixInfo = nil
		} else {
			root.regInfo = nil
		}

		if root.fixInfo == nil && root.regInfo == nil {
			sf.routers.Delete(prefix)
		}
	}
}

// Root 资源根
type Root struct {
	fixInfo *ControllerInfo
	regInfo *ControllerInfo
}

// addToRouter 增加到路由
func (sf *ControllerRegister) addToRouter(resource string, ci *ControllerInfo) error {
	var prefix string
	var isFix bool

	rs := SplitResource(resource)
	length := len(rs)
	if length == 0 {
		return errors.New("resource is empty string")
	}

	//固定路由还是正则路由
	if strings.EqualFold(rs[length-1], "@") {
		prefix = strings.Join(rs[:length-1], ".")
		isFix = false
	} else {
		prefix = strings.Join(rs, ".")
		isFix = true
	}

	if len(prefix) == 0 {
		return errors.New("resource prefix is empty")
	}

	// 增加
	var root *Root
	if v, ok := sf.routers.Load(prefix); ok {
		root = v.(*Root)
	} else {
		root = &Root{}
	}

	if isFix {
		root.fixInfo = ci
	} else {
		root.regInfo = ci
	}
	sf.routers.Store(prefix, root)

	return nil
}
