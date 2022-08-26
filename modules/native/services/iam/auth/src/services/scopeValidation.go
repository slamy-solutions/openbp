package services

import (
	"strings"
)

func compareString(has string, required string) bool {
	hasWildcard := strings.HasSuffix(has, "*")
	requiredWildcard := strings.HasSuffix(required, "*")
	if requiredWildcard && !hasWildcard {
		return false
	}
	if !requiredWildcard && !hasWildcard {
		return has == required
	}

	return strings.HasPrefix(required, has[:len(has)-1])
}

func compareStringList(hasList []string, requiredList []string) bool {
	for _, required := range requiredList {
		accessOk := false
		for _, has := range hasList {
			if compareString(has, required) {
				accessOk = true
				break
			}
		}
		if !accessOk {
			return false
		}
	}

	return true
}

func RequestScope(resources []string, actions []string, requiredResources []string, requiredActions []string) bool {
	return compareStringList(resources, requiredResources) && compareStringList(actions, requiredActions)
}
