package yorm

import "strings"

type columnTag struct {
	skip              bool
	columnName        string
	columnIsSet       bool
	defaultValue      string
	defaultValueIsSet bool
}

//parse value like column(name) ,return name
func parseBracketsValue(toParse, key string) (value string, isSet bool) {
	index := strings.Index(toParse, key)
	if index != 0 {
		return
	}
	keyLen := len(key)
	lengh := len(toParse)
	if keyLen+2 > lengh {
		return
	}
	if toParse[keyLen:keyLen+1] == "(" && toParse[lengh-1:] == ")" {
		return toParse[keyLen+1 : lengh-1], true
	}
	return
}

func parseTag(tagStr string) (t columnTag) {
	defer func() {
		if t.skip {
			//ignore other filed value
			t = columnTag{skip: true}
		}
	}()
	if tagStr == "" {
		return
	}
	tagStrs := strings.Split(tagStr, ";")
	for _, ts := range tagStrs {
		if ts == "-" {
			t.skip = true
			return
		}
		if !t.columnIsSet {
			t.columnName, t.columnIsSet = parseBracketsValue(ts, "column")
			if t.columnIsSet {
				continue
			}
		}
		if !t.defaultValueIsSet {
			t.defaultValue, t.defaultValueIsSet = parseBracketsValue(ts, "default")
			if t.defaultValueIsSet {
				continue
			}
		}
	}
	return
}
