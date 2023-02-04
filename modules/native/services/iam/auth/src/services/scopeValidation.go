package services

import (
	"strings"

	nativeIAmAuthGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmTokenGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
)

/*type Policy nativeIAmPolicyGRPC.Policy

type PolicyLike interface {
	GetNamespace() string
	GetResources() []string
	GetActions() []string
}

func (p *Policy) GetNamespace() string {
	return p.Namespace
}
func (p *Policy) GetResources() []string {
	return p.Resources
}
func (p *Policy) GetActions() []string {
	return p.Actions
}*/

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

func arePoliciesAllowScopes(policies []*nativeIAmPolicyGRPC.Policy, scopes []*nativeIAmAuthGRPC.Scope) bool {
	for _, scope := range scopes {
		scopeIsOk := false
		for _, policy := range policies {
			if scope.Namespace == policy.Namespace && policy.NamespaceIndependent {
				if compareStringList(policy.Resources, scope.Resources) && compareStringList(policy.Actions, scope.Actions) {
					scopeIsOk = true
					break
				}
			}
		}
		if !scopeIsOk {
			return false
		}
	}

	return true
}

func areTokenScopesValidForIdentityScopes(policies []*nativeIAmPolicyGRPC.Policy, scopes []*nativeIAmTokenGRPC.Scope) bool {
	for _, scope := range scopes {
		scopeIsOk := false
		for _, policy := range policies {
			if scope.Namespace == policy.Namespace {
				if compareStringList(policy.Resources, scope.Resources) && compareStringList(policy.Actions, scope.Actions) {
					scopeIsOk = true
					break
				}
			}
		}
		if !scopeIsOk {
			return false
		}
	}

	return true
}

func areTokenScopesAllowAccess(tokenScopes []*nativeIAmTokenGRPC.Scope, systemScopes []*nativeIAmAuthGRPC.Scope) bool {
	for _, systemScope := range systemScopes {
		scopeIsOk := false
		for _, tokenScope := range tokenScopes {
			if systemScope.Namespace == tokenScope.Namespace {
				if compareStringList(tokenScope.Resources, systemScope.Resources) && compareStringList(tokenScope.Actions, systemScope.Actions) {
					scopeIsOk = true
					break
				}
			}
		}
		if !scopeIsOk {
			return false
		}
	}

	return true
}
