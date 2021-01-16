package util

import "strings"

var asta byte = '*'
var asta1 = "*"
var asta2 = "**"

type compType int8

const c_EQU compType = 0
const c_STR compType = 1
const c_STR_MID compType = 2
const c_STR_END compType = 3
const c_MID compType = 4
const c_MID_MID compType = 5
const c_MID_END compType = 6
const c_END compType = 7
const c_ANY compType = 8

type StrMatch struct {
	Pattern               string
	comp                  compType
	start, end, mid, mid2 string
}

func NewStrMatch(pattern string) StrMatch {
	m := StrMatch{Pattern: pattern}
	if pattern == asta1 || pattern == asta2 {
		m.comp = c_ANY
		return m
	}
	length := len(pattern)
	if length < 2 {
		m.comp = c_EQU
		m.mid = pattern
		return m
	}

	anyStart := pattern[0] == asta
	anyEnd := pattern[length-1] == asta
	apos := strings.IndexByte(pattern[1:], asta)
	if apos >= 0 {
		apos++
	}
	anyMid := apos > 0 && apos < length-1

	if anyMid {
		if anyStart && anyEnd {
			m.comp = c_MID_MID
			m.mid = pattern[1:apos]
			m.mid2 = pattern[apos+1 : length-1]
		} else if anyStart {
			m.comp = c_MID_END
			m.mid = pattern[1:apos]
			m.end = pattern[apos+1:]
		} else if anyEnd {
			m.comp = c_STR_MID
			m.start = pattern[0:apos]
			m.mid = pattern[apos+1 : length-1]
		} else {
			m.comp = c_STR_END
			m.start = pattern[0:apos]
			m.end = pattern[apos+1:]
		}
	} else {
		if anyStart && anyEnd {
			m.comp = c_MID
			m.mid = pattern[1 : length-1]
		} else if anyStart {
			m.comp = c_END
			m.end = pattern[1:]
		} else if anyEnd {
			m.comp = c_STR
			m.start = pattern[:length-1]
		} else {
			m.comp = c_EQU
			m.mid = pattern
		}
	}
	return m
}

func (m *StrMatch) include(target string) bool {
	if target == "" {
		return false;
	}
	switch m.comp {
	case c_ANY:
		return true
	case c_EQU:
		return m.mid == target
	case c_STR:
		return strings.HasPrefix(target, m.start)
	case c_STR_MID:
		return strings.HasPrefix(target, m.start) && strings.Contains(target, m.mid)
	case c_STR_END:
		return strings.HasPrefix(target, m.start) && strings.HasSuffix(target, m.end)
	case c_MID:
		return strings.Contains(target, m.mid)
	case c_MID_MID:
		x := strings.Index(target, m.mid)
		if x < 0 {
			return false
		}
		return strings.Contains(target, m.mid2) && strings.HasSuffix(target, m.end)
	case c_MID_END:
		return strings.Contains(target, m.mid) && strings.HasSuffix(target, m.end)
	case c_END:
		return strings.HasSuffix(target, m.end)
	default:
		return false
	}
}


type CommaSeparatedChainedStrMatcher struct {
	strMatches []StrMatch
}

func NewCommaSeparatedChainedStrMatcher(patterns string) CommaSeparatedChainedStrMatcher {
	matcher := CommaSeparatedChainedStrMatcher{}
	arrPatterns := strings.Split(patterns, ",")
	for _, pattern := range arrPatterns {
		matcher.strMatches = append(matcher.strMatches, StrMatch{Pattern: pattern})
	}
	return matcher
}

func (m CommaSeparatedChainedStrMatcher) IsMatch(target string) bool {
	for _, match := range m.strMatches {
		if match.include(target) {
			return true
		}
	}
	return false
}
