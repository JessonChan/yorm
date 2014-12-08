package yorm

import "testing"

func TestParseValue_1(t *testing.T) {
	v, set := parseBracketsValue("column()", "column")
	if v != "" && set {
		t.Log(v)
		t.FailNow()
	}
}
func TestParseValue_2(t *testing.T) {
	v, set := parseBracketsValue("column(tag)", "column")
	if v != "tag" && set {
		t.Log(v)
		t.FailNow()
	}
}

func TestParseValue_3(t *testing.T) {
	v, set := parseBracketsValue("name(jc)", "nam")
	if v == "" && !set {
	} else {
		t.Fail()
	}
}

func TestParseTag_1(t *testing.T) {
	tag := parseTag("-")
	if tag.skip && tag.columnName == "" && tag.defaultValue == "" {
	} else {
		t.Log(tag)
		t.Fail()
	}
}
func TestParseTag_2(t *testing.T) {
	tag := parseTag("column(tag)")
	assert(t, !tag.skip, tag.columnName == "tag", tag.defaultValue == "")
}
func TestParseTag_3(t *testing.T) {
	tag := parseTag("default(1)")
	assert(t, tag.defaultValue == "1")
}
func TestParseTag_4(t *testing.T) {
	tag := parseTag("-;column(tag)")
	assertTrue(tag.skip, t)
	assertTrue(tag.columnName == "", t)
}

func TestParseTag_5(t *testing.T) {
	tag := parseTag("column(tag);-")
	assertTrue(tag.skip, t)
	assertTrue(tag.columnName == "", t)
}

func TestParseTag_6(t *testing.T) {
	tag := parseTag("column(tag1);default(1)")
	assertTrue(!tag.skip, t)
	assertTrue(tag.columnName == "tag1", t)
	assertTrue(tag.defaultValue == "1", t)
}
func assert(t *testing.T, b ...bool) {
	log := ""
	for _, v := range b {
		if !v {
			log += "false"
			t.Log(log)
			t.Fail()
		}
		log += "true,"
	}
}

func assertTrue(b bool, t *testing.T) {
	if !b {
		t.FailNow()
	}
}
