package elink

import (
	"reflect"
	"testing"
)

func TestDecodeTopic(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		want    *TopicLayer
		wantErr bool
	}{
		{"主题格式非法", args{"a/b/c/d/e"},
			nil, true},
		{"不含查询参数", args{"a/b/c/d/e/1"},
			&TopicLayer{"a", "b", "c", "d", "e", "1", ""}, false},
		{"含查询参数", args{"a/b/c/d/e/1/key1=a&key2=b"},
			&TopicLayer{"a", "b", "c", "d", "e", "1", "key1=a&key2=b"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeTopic(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopicLayer_Query(t *testing.T) {
	tests := []struct {
		name string
		this *TopicLayer
		want QueryValues
	}{
		{"", &TopicLayer{RawQuery: "key1=a&key2=b"}, map[string][]string{"key1": []string{"a"}, "key2": []string{"b"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopicLayer.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeTopic(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{[]string{"a", "b", "c"}}, "a/b/c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeTopic(tt.args.s); got != tt.want {
				t.Errorf("EncodeTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeReplyTopic(t *testing.T) {
	type args struct {
		tp    *TopicLayer
		query []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"不包含query", args{tp: &TopicLayer{
			Channel: "channel", ProductKey: "key", Mac: "mac", Resource: "res", Method: "method", UserKey: "1"}},
			"channel_reply/key/mac/res/method/1"},
		{"包含query", args{tp: &TopicLayer{
			Channel: "channel", ProductKey: "key", Mac: "mac", Resource: "res", Method: "method", UserKey: "1"},
			query: []string{"query=10"}},
			"channel_reply/key/mac/res/method/1/query=10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeReplyTopic(tt.args.tp, tt.args.query...); got != tt.want {
				t.Errorf("EncodeReplyTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodePushTopic(t *testing.T) {
	type args struct {
		channel     string
		productKey  string
		mac         string
		resource    string
		method      string
		messageType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"channel", "key", "mac", "res", "method", "type"},
			"channel/key/mac/res/method/type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodePushTopic(tt.args.channel, tt.args.productKey, tt.args.mac, tt.args.resource, tt.args.method, tt.args.messageType); got != tt.want {
				t.Errorf("EncodePushTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatResource(t *testing.T) {
	type args struct {
		prefix string
		attach []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"只有prefix", args{prefix: "pre"}, "pre"},
		{"带1个attach", args{"pre", []string{"att"}}, "pre.att"},
		{"带2个attach", args{"pre", []string{"att", "att1"}}, "pre.att.att1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatResource(tt.args.prefix, tt.args.attach...); got != tt.want {
				t.Errorf("FormatResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitResource(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{` -> []`, args{""}, []string{}},
		{`. -> []`, args{"."}, []string{}},
		{`resource -> [resource]`, args{"resource"}, []string{"resource"}},
		{"resource. -> [resource]", args{"resource."}, []string{"resource"}},
		{"resource.@ -> [resource @]", args{"."}, []string{}},
		{"resource.device -> [resource,device]", args{"resource.device"}, []string{"resource", "device"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitResource(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryValues(t *testing.T) {
	tests := []struct {
		name string
		want QueryValues
	}{
		{"", QueryValues{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDecodeQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want QueryValues
	}{
		{"单个值",
			args{query: "a=1&b=2"},
			QueryValues{"a": []string{"1"}, "b": []string{"2"}},
		},
		{
			"多个值",
			args{query: "a=1&a=2&a=banana"},
			QueryValues{"a": []string{"1", "2", "banana"}},
		},
		{
			"带有空的键",
			args{query: "=1&a=2&a=banana"},
			QueryValues{"a": []string{"2", "banana"}, "": []string{"1"}},
		},
		{
			"&起头",
			args{query: "&a=2&a=banana"},
			QueryValues{"a": []string{"2", "banana"}},
		},
	}
	for _, tt := range tests {
		form := DecodeQuery(tt.args.query)
		if len(form) != len(tt.want) {
			t.Errorf("test %s: len(form) = %d, want %d", tt.name, len(form), len(tt.want))
		}
		for k, evs := range tt.want {
			vs, ok := form[k]
			if !ok {
				t.Errorf("test %s: Missing key %q", tt.name, k)
				continue
			}
			if len(vs) != len(evs) {
				t.Errorf("test %s: len(form[%q]) = %d, want %d", tt.name, k, len(vs), len(evs))
				continue
			}
			for j, ev := range evs {
				if v := vs[j]; v != ev {
					t.Errorf("test %s: form[%q][%d] = %q, want %q", tt.name, k, j, v, ev)
				}
			}
		}
	}
}

func TestQueryValues_QueryValue(t *testing.T) {
	var vNil QueryValues
	if g, e := vNil.Get("foo"), ""; g != e {
		t.Errorf("Get(foo) = %q, want %q", g, e)
	}
	// Case sensitive:
	if g, e := vNil.GetValues("bar"), []string{}; !reflect.DeepEqual(g, e) {
		t.Errorf("GetValues(bar) = %v, want %v", g, e)
	}

	v := DecodeQuery("foo=bar&bar=1&bar=2")
	if g, e := v.Get("foo"), "bar"; g != e {
		t.Errorf("Get(foo) = %q, want %q", g, e)
	}
	// Case sensitive:
	if g, e := v.GetValues("bar"), []string{"1", "2"}; !reflect.DeepEqual(g, e) {
		t.Errorf("GetValues(bar) = %v, want %v", g, e)
	}

	// Case sensitive:
	if g, e := v.Get("Foo"), ""; g != e {
		t.Errorf("Get(Foo) = %q, want %q", g, e)
	}
	if g, e := v.Get("bar"), "1"; g != e {
		t.Errorf("Get(bar) = %q, want %q", g, e)
	}
	if g, e := v.Get("baz"), ""; g != e {
		t.Errorf("Get(baz) = %q, want %q", g, e)
	}
	v.Del("bar")
	if g, e := v.Get("bar"), ""; g != e {
		t.Errorf("second Get(bar) = %q, want %q", g, e)
	}
	v.Set("foo", "xxx")
	if g, e := v.Get("foo"), "xxx"; g != e {
		t.Errorf("second Get(foo) = %q, want %q", g, e)
	}
	v.Add("add", "haha")
	if g, e := v.Get("add"), "haha"; g != e {
		t.Errorf("second Get(add) = %q, want %q", g, e)
	}
}

func TestQueryValues_EncodeQuery(t *testing.T) {
	tests := []struct {
		name string
		this QueryValues
		want string
	}{
		{"", nil, ""},
		{"", QueryValues{"q": {"puppies"}, "oe": {"utf8"}}, "oe=utf8&q=puppies"},
		{"", QueryValues{"q": {"dogs", "b", "7"}}, "q=dogs&q=b&q=7"},
		{"", QueryValues{
			"a": {"a1", "a2", "a3"},
			"b": {"b1", "b2", "b3"},
			"c": {"c1", "c2", "c3"},
		}, "a=a1&a=a2&a=a3&b=b1&b=b2&b=b3&c=c1&c=c2&c=c3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.EncodeQuery(); got != tt.want {
				t.Errorf("QueryValues.EncodeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
