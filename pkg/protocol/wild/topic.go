package wild

import (
	"strings"
)

const (
	// MWC 多层通配符
	MWC = "#"
	// SWC 单层通配符
	SWC = "+"
	// SEP 主题分隔符
	SEP = "/"
	// SYS 系统级起始符
	SYS = "$"
	// Both wildcards
	_WC = "#+"
)

type Wild struct {
	wild []string
}

func NewWild(topic string) *Wild {
	return &Wild{strings.Split(topic, SEP)}
}

// 匹配实际主题, 进行层次匹配
func (this *Wild) Matches(parts []string) bool {
	i := 0
	for i < len(parts) {
		// 主题过长,不匹配
		if i >= len(this.wild) {
			return false
		}

		// 通配符"#",匹配所有
		if this.wild[i] == MWC {
			return true
		}

		// 对应层级字符不匹配,或不是层级通配符"+",表示不匹配
		if parts[i] != this.wild[i] && this.wild[i] != SWC {
			return false
		}
		i++
	}

	// 使 a/b/c/# 匹配 a/b/c
	if i == len(this.wild)-1 && this.wild[len(this.wild)-1] == MWC {
		return true
	}

	return i == len(this.wild)
}

// 判断主题有效性
// 单层中不允许出现 finance{{.MWC}}
// MWC通配符只允许出现在最后一层
//
func (this *Wild) Valid() bool {
	for i, part := range this.wild {
		// 单层中允许出现类似 finance#
		if HasWildcard(part) && len(part) != 1 {
			return false
		}
		// MWC 通配符只允许出现在出后一级
		if part == MWC && i != len(this.wild)-1 {
			return false
		}
	}
	return true
}

func HasWildcard(topic string) bool {
	return (strings.Contains(topic, MWC) || strings.Contains(topic, SWC))
}
