package elink

import (
	"reflect"
	"testing"
)

func TestControllerRegister_AddToRouter(t *testing.T) {
	type args struct {
		resource   string
		c          ControllerInterface
		mapMethods []string
	}
	tests := []struct {
		name    string
		this    *ControllerRegister
		args    args
		wantErr bool
	}{
		{"固定路由,unknown方法", NewControllerRegister(), args{"abcd", &Controller{}, []string{"unknown:get"}}, true},
		{"固定路由,无效map方法", NewControllerRegister(), args{"abcd", &Controller{}, []string{"aa"}}, true},
		{"固定路由,控制器未有对应方法", NewControllerRegister(), args{"abcd", &Controller{}, []string{"aa:bb"}}, true},
		{"固定路由,常规方法", NewControllerRegister(), args{resource: "abcd", c: &Controller{}}, false},
		{"固定路由,自定义方法", NewControllerRegister(), args{"abcd", &Controller{}, []string{"aa:StopRun"}}, false},
		{"匹配路由,常规方法", NewControllerRegister(), args{"abc.@", &Controller{}, []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.AddToRouter(tt.args.resource, tt.args.c, tt.args.mapMethods...); (err != nil) != tt.wantErr {
				t.Errorf("ControllerRegister.AddToRouter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func newControllerRegisterWithRouter() *ControllerRegister {
	cr := NewControllerRegister()
	_ = cr.AddToRouter("j.k.h", &Controller{})
	_ = cr.AddToRouter("a.b.c", &Controller{})
	_ = cr.AddToRouter("a.b.c.@", &Controller{})
	_ = cr.AddToRouter("a.b.@", &Controller{})
	return cr
}
func TestControllerRegister_MatchRouter(t *testing.T) {

	type args struct {
		resource string
	}
	tests := []struct {
		name      string
		this      *ControllerRegister
		args      args
		wantQuery QueryValues
		wantErr   bool
	}{
		{"不匹配", newControllerRegisterWithRouter(), args{"a.d"}, nil, true},
		{"全匹配", newControllerRegisterWithRouter(), args{"a.b.c"}, QueryValues{}, false},
		{"长匹配", newControllerRegisterWithRouter(), args{"a.b.c.d"}, QueryValues{"splat": []string{"d"}}, false},
		{"短匹配", newControllerRegisterWithRouter(), args{"a.b.d.e"}, QueryValues{"splat": []string{"d", "e"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotQuery, err := tt.this.MatchRouter(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ControllerRegister.MatchRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Errorf("ControllerRegister.MatchRouter() gotQuery = %v, wantQuery %v", gotQuery, tt.wantQuery)
				return
			}
		})
	}
}

func TestControllerRegister_DeleteRouter(t *testing.T) {
	cr := newControllerRegisterWithRouter()
	cr.DeleteRouter("a.b.@")
	_, _, err := cr.MatchRouter("a.b.d")
	if err == nil {
		t.Errorf("ControllerRegister.MatchRouter() error = %v, wantErr %v", err, true)
		return
	}
	_, _, err = cr.MatchRouter("j.k.h")
	if err != nil {
		t.Errorf("ControllerRegister.MatchRouter() error = %v, wantErr %v", err, false)
		return
	}
	cr.DeleteRouter("j.k.h")
	_, _, err = cr.MatchRouter("j.k.h")
	if err == nil {
		t.Errorf("ControllerRegister.MatchRouter() error = %v, wantErr %v", err, true)
		return
	}
}

func TestControllerRegister_addToRouter(t *testing.T) {

	type args struct {
		resource string
		ci       *ControllerInfo
	}
	tests := []struct {
		name    string
		this    *ControllerRegister
		args    args
		wantErr bool
	}{
		{"空资源", NewControllerRegister(), args{"", &ControllerInfo{}}, true},
		{"前置空 如.@", NewControllerRegister(), args{".@", &ControllerInfo{}}, true},
		{"固定资源", NewControllerRegister(), args{"a.b.c", &ControllerInfo{}}, false},
		{"匹配资源", NewControllerRegister(), args{"a.b.@", &ControllerInfo{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.addToRouter(tt.args.resource, tt.args.ci); (err != nil) != tt.wantErr {
				t.Errorf("ControllerRegister.addToRouter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
