package elink

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	// MWC 多层通配符
	MWC = "#"
	// SWC 单层通配符
	SWC = "+"
	// SEP 主题分隔符
	SEP = "/"
	// RSEP 资源分隔符
	RSEP = "."
	// QSEP 参数分隔符
	QSEP = "&"
	// KVSEP 键值分隔符
	KVSEP = "="
)

var ErrInvalidTopicLength = errors.New("invalid topic length")

// Request 请求信息
type Request struct {
	Topic   *TopicLayer // 主题
	Param   QueryValues // resource上提取的参数
	Values  QueryValues // 查询参数 key=val1&key1=val2
	Payload []byte
}

// TopicLayer 主题层级
type TopicLayer struct {
	Channel, ProductKey, Mac string
	Resource, Method         string
	UserKey                  string
	RawQuery                 string
}

// DecodeTopic 解析主题层级,分隔符“/” channel/productKey/mac/resource/method/query1=xx&query2=xxx
func DecodeTopic(topic string) (*TopicLayer, error) {
	s := strings.Split(topic, SEP)
	if len(s) < 6 {
		return nil, ErrInvalidTopicLength
	}

	tp := &TopicLayer{
		Channel:    s[0],
		ProductKey: s[1],
		Mac:        s[2],
		Resource:   s[3],
		Method:     s[4],
		UserKey:    s[5],
	}

	if len(s) >= 7 {
		tp.RawQuery = s[6]
	}
	return tp, nil
}

// Query 解析主题层级查询参数
func (this *TopicLayer) Query() QueryValues {
	return DecodeQuery(this.RawQuery)
}

// EncodeTopic 编码主题格式,分隔符"/" 例topic/sub1/sub2/sub3
func EncodeTopic(s []string) string {
	return strings.Join(s, SEP)
}

// EncodeReplyTopic 回复通道主题, 例 channel/productKey/mac/resource/method/query
func EncodeReplyTopic(tp *TopicLayer, query ...string) string {
	s := []string{fmt.Sprintf("%s_reply", tp.Channel), tp.ProductKey, tp.Mac,
		tp.Resource, tp.Method, tp.UserKey}
	if len(query) > 0 {
		s = append(s, query[0])
	}
	return EncodeTopic(s)
}

// EncodePushTopic 格式化主题
func EncodePushTopic(channel, productKey, mac, resource, method, messageType string) string {
	return EncodeTopic([]string{channel, productKey, mac, resource,
		strings.ToLower(method), strings.ToLower(messageType)})
}

// FormatResource 组成资源,使用"."分隔, 例device.sub1.sub2.sub3
// prefix一般是固定的, attach主要是动态变化,根据情况组合
func FormatResource(prefix string, attach ...string) string {
	if len(attach) == 0 {
		return prefix
	}

	s := make([]string, 0, len(attach)+1)
	s = append(s, prefix)
	s = append(s, attach...)
	return strings.Join(s, RSEP)
}

// SplitResource 分割资源
// "." - > []
// "resource" -> ["resource"]
// "resource." -> ["resource"]
// "resource.@" -> ["resource","@"]
// "resource.device" -> ["resource", "device"]
//分割资源,返回一个切片,不包含'.'
func SplitResource(s string) []string {
	s = strings.Trim(s, ".")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ".")
}

// QueryValues 查询参数
type QueryValues map[string][]string

// NewQueryValues 新建查询参数对象
func NewQueryValues() QueryValues {
	return make(QueryValues)
}

// DecodeQuery 解析主题层次的查询参数
func DecodeQuery(query string) QueryValues {
	m := make(QueryValues)
	for query != "" {
		key := query
		if i := strings.IndexAny(key, QSEP); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, KVSEP); i >= 0 {
			key, value = key[:i], key[i+1:]
		}

		m[key] = append(m[key], value)
	}
	return m
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (sf QueryValues) Get(key string) string {
	if sf == nil {
		return ""
	}
	vs := sf[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// GetValue gets the value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string slice. To access multiple values, use the map
// directly.
func (sf QueryValues) GetValues(key string) []string {
	if sf == nil {
		return []string{}
	}
	return sf[key]
}

// Set sets the key to value. It replaces any existing values.
func (sf QueryValues) Set(key, value string) QueryValues {
	sf[key] = []string{value}
	return sf
}

// Add adds the value to key. It appends to any existing values associated with key.
func (sf QueryValues) Add(key string, value ...string) QueryValues {
	sf[key] = append(sf[key], value...)
	return sf
}

// Del deletes the values associated with key.
func (sf QueryValues) Del(key string) {
	delete(sf, key)
}

// EncodeQuery encodes the QueryValues into form ("bar=baz&foo=qux") sorted by key.
func (sf QueryValues) EncodeQuery() string {
	if sf == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(sf))
	for k := range sf {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := sf[k]
		for _, this := range vs {
			if buf.Len() > 0 {
				buf.WriteString(QSEP)
			}
			buf.WriteString(k)
			buf.WriteString(KVSEP)
			buf.WriteString(this)
		}
	}
	return buf.String()
}
