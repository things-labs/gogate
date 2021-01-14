package elink

import (
	"reflect"
	"testing"
)

func Test_CodeErrorMessage(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want CodeErrorMessageInfo
	}{
		{"not defined", args{1000}, codeErrorMessageList[CodeErrSysNotSupport]},
		{"defined", args{CodeErrSysInProcess}, codeErrorMessageList[CodeErrSysInProcess]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CodeErrorMessage(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CodeErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
