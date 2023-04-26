package policy

import (
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

var writeParamsValidationTestCases = []struct {
	testName    string
	name        string
	description string
	resources   []string
	actions     []string
	ok          bool
}{
	{"All OK (actions and resources empty)", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{}, true},
	{"All OK (actions empty)", tools.GetRandomString(20), tools.GetRandomString(20), []string{"sonething"}, []string{}, true},
	{"All OK", tools.GetRandomString(20), tools.GetRandomString(20), []string{"sonething"}, []string{"action"}, true},
	{"All OK (resource with asterisk)", tools.GetRandomString(20), tools.GetRandomString(20), []string{"sonething*"}, []string{"action"}, true},
	{"All OK (action with asterisk)", tools.GetRandomString(20), tools.GetRandomString(20), []string{"sonething"}, []string{"action*"}, true},
	{"All OK (all with asterisk)", tools.GetRandomString(20), tools.GetRandomString(20), []string{"sonething*"}, []string{"action*"}, true},
	{"All OK (full resource with asterisk)", tools.GetRandomString(20), tools.GetRandomString(20), []string{"*"}, []string{}, true},
	{"All OK (full action with asterisk)", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"*"}, true},
	{"Fail when asterisk in the middle of resource)", tools.GetRandomString(20), tools.GetRandomString(20), []string{"sonething*lala"}, []string{}, false},
	{"Fail when asterisk in the middle of action)", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"sonething*lala"}, false},
	{"Fail when multiple asteriscs in resource. Type 1", tools.GetRandomString(20), tools.GetRandomString(20), []string{"some**"}, []string{}, false},
	{"Fail when multiple asteriscs in resource. Type 2", tools.GetRandomString(20), tools.GetRandomString(20), []string{"so*me*"}, []string{}, false},
	{"Fail when multiple asteriscs in resource. Type 3", tools.GetRandomString(20), tools.GetRandomString(20), []string{"*some*"}, []string{}, false},
	{"Fail when multiple asteriscs in resource. Type 4", tools.GetRandomString(20), tools.GetRandomString(20), []string{"*som*e"}, []string{}, false},
	{"Fail when multiple asteriscs in resource. Type 5", tools.GetRandomString(20), tools.GetRandomString(20), []string{"**some"}, []string{}, false},
	{"Fail when multiple asteriscs in resource. Type 6", tools.GetRandomString(20), tools.GetRandomString(20), []string{"**"}, []string{}, false},
	{"Fail when multiple asteriscs in actions. Type 1", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"some**"}, false},
	{"Fail when multiple asteriscs in actions. Type 2", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"so*me*"}, false},
	{"Fail when multiple asteriscs in actions. Type 3", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"*some*"}, false},
	{"Fail when multiple asteriscs in actions. Type 4", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"*som*e"}, false},
	{"Fail when multiple asteriscs in actions. Type 5", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"**some"}, false},
	{"Fail when multiple asteriscs in actions. Type 6", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"**"}, false},
	{"Fail when name is empty", "", tools.GetRandomString(20), []string{}, []string{}, false},
	{"OK when description is empty", tools.GetRandomString(20), "", []string{}, []string{}, true},

	{"Fail when invalid character in resource. Type 1", tools.GetRandomString(20), tools.GetRandomString(20), []string{"some#"}, []string{}, false},
	{"Fail when invalid character in resource. Type 2", tools.GetRandomString(20), tools.GetRandomString(20), []string{"some%"}, []string{}, false},
	{"Fail when invalid character in resource. Type 3", tools.GetRandomString(20), tools.GetRandomString(20), []string{"some thing"}, []string{}, false},
	{"Fail when invalid character in resource. Type 4", tools.GetRandomString(20), tools.GetRandomString(20), []string{"some!"}, []string{}, false},
	{"Fail when invalid character in resource. Type 5", tools.GetRandomString(20), tools.GetRandomString(20), []string{"some@"}, []string{}, false},

	{"Fail when invalid character in action. Type 1", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"some#"}, false},
	{"Fail when invalid character in action. Type 2", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"some%"}, false},
	{"Fail when invalid character in action. Type 3", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"some thing"}, false},
	{"Fail when invalid character in action. Type 4", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"some!"}, false},
	{"Fail when invalid character in action. Type 5", tools.GetRandomString(20), tools.GetRandomString(20), []string{}, []string{"some@"}, false},

	{"OK while there are dots in resource", tools.GetRandomString(20), "", []string{"some.level.resource"}, []string{}, true},
	{"OK while there are dots in action", tools.GetRandomString(20), "", []string{}, []string{"some.level.action"}, true},
	{"OK while there are usenrscores in resource", tools.GetRandomString(20), "", []string{"some.level_resource.thing"}, []string{}, true},
	{"OK while there are usenrscores in action", tools.GetRandomString(20), "", []string{}, []string{"some.level_action.thing"}, true},

	{"Failed while resource is empty", tools.GetRandomString(20), "", []string{""}, []string{}, false},
	{"Failed while action is empty", tools.GetRandomString(20), "", []string{}, []string{""}, false},
}
