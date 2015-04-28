package yorm

import "strings"

type columnTag struct {
	skip              bool
	columnIsSet       bool
	defaultValueIsSet bool

	columnName   string
	defaultValue string
}

//parse value like column(name) ,return name
func parseBracketsValue(toParse, key string) (value string, isSet bool) {
	index := strings.Index(toParse, key)
	if index != 0 {
		return
	}
	keyLen := len(key)
	length := len(toParse)
	if keyLen+2 > length { // 2 == len("(") + len(")")
		return
	}
	if toParse[keyLen:keyLen+1] == "(" && toParse[length-1:] == ")" {
		return toParse[keyLen+1 : length-1], true
	}
	return
}

func parseTag(tagStr string) (t columnTag) {
	if tagStr == "" {
		return
	}

	defer func() {
		if t.skip {
			//ignore other filed value
			t = columnTag{skip: true}
		}
	}()

	tags := strings.Split(tagStr, ";")
	for _, tag := range tags {
		if tag == "-" {
			t.skip = true
			return
		}
		if !t.columnIsSet {
			t.columnName, t.columnIsSet = parseBracketsValue(tag, "column")
			if t.columnIsSet {
				continue
			}
		}
		if !t.defaultValueIsSet {
			t.defaultValue, t.defaultValueIsSet = parseBracketsValue(tag, "default")
			if t.defaultValueIsSet {
				continue
			}
		}
	}
	return
}
