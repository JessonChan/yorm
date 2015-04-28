package yorm

import "testing"

func TestCamelToUnserscore_1(t *testing.T) {
	if camelToUnderscore("") != "" {
		t.Fail()
	}
	if camelToUnderscore("U") != "u" {
		t.Fail()
	}
	if camelToUnderscore("UserName") != "user_name" {
		t.Fail()
	}
}
func TestUnderscoreToCamel_2(t *testing.T) {
	if underscoreToCamel("user_name") != "UserName" {
		t.Fail()
	}
}
