package yorm

func camel2underscore(name string) string {
	if name == "" {
		return name
	}
	bs := make([]rune, 0, 2*len(name))
	for _, s := range name {
		if 'A' <= s && s <= 'Z' {
			s += ('a' - 'A')
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
	ns := make([]rune, 0, len(name))
	isUnder := true
	for _, v := range name {
		r := v
		if isUnder {
			if v >= 'a' && v <= 'z' {
				r -= ('a' - 'A')
			}
		}
		isUnder = v == '_'
		if isUnder {
			continue
		}
		ns = append(ns, r)
	}
	return string(ns)
}
