package yorm

import "strings"

func camel2underscore(name string) string {
	if name == "" {
		return name
	}
	var bs []rune
	for _, s := range name {
		if 'A' <= s && s <= 'Z' {
			s += 32
			bs = append(bs, '_')
		}
		bs = append(bs, s)
	}
	if bs[0] == '_' {
		return string(bs[1:])
	}
	return string(bs)
}
func underscore2camel(name string) string {
	lengh := len(name)
	if lengh <= 2 {
		return name
	}
	ss := strings.Split(name, "_")
	ns := ""
	for _, s := range ss {
		if s != "" {
			rs := make([]rune, len([]rune(s)))
			copy(rs, []rune(s))
			if s[0] > 'a' && s[0] < 'z' {
				rs[0] -= 32
			}
			ns += string(rs)
		}
	}
	return ns
}
