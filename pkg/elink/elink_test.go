package elink

import (
	"reflect"
	"testing"
)

var testCenters = NewCenter()

func TestCenter_Put(t *testing.T) {
	type args struct {
		channel string
		c       *ControllerRegister
	}
	tests := []struct {
		name    string
		this    *Center
		args    args
		wantErr bool
	}{
		{"第一次添加通道", testCenters, args{"raw", NewControllerRegister()}, false},
		{"重复添加通道", testCenters, args{"raw", NewControllerRegister()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.Put(tt.args.channel, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Center.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCenter_CenterOperate(t *testing.T) {
	testCenters.SetPanicHandler(func(interface{}) {})
	testCenters.SetErrorHandler(func(error) {})
	_ = testCenters.Put("raw", NewControllerRegister())

	list, want := testCenters.ChannelSelectorList(), []string{"raw"}
	if !reflect.DeepEqual(list, want) {
		t.Errorf("Center.ChannelSelectorList() = %v, want %v", list, want)
	}

	b := testCenters.Contains("raw")
	if !b {
		t.Errorf("Center.Contains(a) bool = %v, wantErr %v", b, true)
	}
	testCenters.Delete("raw")
	b = testCenters.Contains("raw")
	if b {
		t.Errorf("Center.Contains(a) bool = %v, wantErr %v", b, false)
	}
	_ = testCenters.Put("raw", NewControllerRegister())
	list, want = testCenters.ChannelSelectorList(), []string{"raw"}
	if !reflect.DeepEqual(list, want) {
		t.Errorf("Center.ChannelSelectorList() = %v, want %v", list, want)
	}
}

func TestCenter_Router(t *testing.T) {
	_ = testCenters.Put("raw", NewControllerRegister())
	type args struct {
		channel    string
		resource   string
		c          ControllerInterface
		mapMethods []string
	}
	tests := []struct {
		name    string
		this    *Center
		args    args
		wantErr bool
	}{
		{"无效通道", testCenters, args{c: &Controller{}}, true},
		{"无效资源", testCenters, args{channel: "raw", c: &Controller{}}, true},
		{"通道未注册", testCenters, args{channel: "ctrl", resource: "a", c: &Controller{}}, true},
		{"通道存在", testCenters, args{channel: "raw", resource: "a", c: &Controller{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.Router(tt.args.channel, tt.args.resource, tt.args.c, tt.args.mapMethods...); (err != nil) != tt.wantErr {
				t.Errorf("Center.Router() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type provide struct{}

func (this *provide) ErrorDefaultResponse(topic string) error {
	return nil
}

func (this *provide) WriteResponse(topic string, data interface{}) error {
	return nil
}
func TestCenter_Elink(t *testing.T) {

	_ = testCenters.Put("raw", NewControllerRegister())
	_ = testCenters.Router("raw", "ccc", &Controller{})
	_ = testCenters.Router("raw", "ccc", &Controller{}, "xxx:StopRun")

	pr := &provide{}

	err := testCenters.Server(pr, "raw/aaa/bbb/ccc/ddd", nil)
	if err == nil {
		t.Errorf("Center.Server error = %v,wantErr %v", err, ErrInvalidTopicLength)
	}

	err = testCenters.Server(pr, "raw/aaa/bbb/ccc/get/eee", nil)
	if err != nil {
		t.Errorf("Center.Server error = %v,wantErr %v", err, nil)
	}

	_ = testCenters.Server(pr, "xxx/aaa/bbb/ccc/get/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/cccc/get/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/ccc/get/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/ccc/ddd/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/ccc/post/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/ccc/put/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/ccc/delete/eee", nil)
	_ = testCenters.Server(pr, "raw/aaa/bbb/ccc/xxx/eee", nil)
}

func TestErrorDefaultResponse(t *testing.T) {
	type args struct {
		p    Provider
		tp   *TopicLayer
		code int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ErrorDefaultResponse(tt.args.p, tt.args.tp, tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("ErrorDefaultResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
