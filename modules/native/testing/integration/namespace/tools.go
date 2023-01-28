package namespace

import (
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

var createParamsValidationTestCases = []struct {
	testName    string
	name        string
	fullName    string
	description string
	ok          bool
}{
	{"All OK", tools.GetRandomString(20), "Some big name", "123", true},
	{"Max name length is 32", tools.GetRandomString(33), "Some big name", "123", false},
	{"Invalid characters in name", tools.GetRandomString(16) + "%", "Some big name", "123", false},
	{"Name must not be empty", "", "Some big name", "123", false},
	{"Empty fullName is ok", tools.GetRandomString(20), "", "123", true},
	{"Max fullName length is 128", tools.GetRandomString(20), tools.GetRandomString(130), "123", false},
	{"Empty description is ok", tools.GetRandomString(20), "alalal", "", true},
	{"Max description length is 512", tools.GetRandomString(20), tools.GetRandomString(10), tools.GetRandomString(530), false},
}
