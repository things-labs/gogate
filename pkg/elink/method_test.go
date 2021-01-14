package elink

import "testing"

func TestHasMethod(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"小写的方法", args{"get"}, true},
		{"大写的方法", args{"GET"}, true},
		{"混合大小写的方法", args{"Get"}, true},
		{"不包含的方法", args{"option"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasMethod(tt.args.method); got != tt.want {
				t.Errorf("HasMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsUnknownMethod(t *testing.T) {
	type args struct {
		str []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"包含小写unknown方法", args{[]string{"unknown", "opt"}}, true},
		{"包含大写unknown方法", args{[]string{"UNKNOWN", "opt"}}, true},
		{"包含混合大小写unknown方法", args{[]string{"UNknown", "opt"}}, true},
		{"不包含unknown方法", args{[]string{"te", "opt"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsUnknownMethod(tt.args.str); got != tt.want {
				t.Errorf("ContainsUnknownMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasMessageType(t *testing.T) {
	type args struct {
		msgType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"包含小写的消息类型", args{"time"}, true},
		{"包含大写的消息类型", args{"TIME"}, true},
		{"包含混合大小写的消息类型", args{"TIme"}, true},
		{"不包含的消息类型", args{"message"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasMessageType(tt.args.msgType); got != tt.want {
				t.Errorf("HasMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}
