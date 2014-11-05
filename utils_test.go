package yorm

import "testing"

func TestCamel2Unserscore_1(t *testing.T) {
	if camel2underscore("") != "" {
		t.Fail()
	}
	if camel2underscore("U") != "u" {
		t.Fail()
	}
	if camel2underscore("UserName") != "user_name" {
		t.Fail()
	}
}
func TestUnderscore2Camel_2(t *testing.T) {
	if underscore2camel("user_name") != "UserName" {
		t.Fail()
	}
}
