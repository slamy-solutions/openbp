package password

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/system/testing/tools"
)

type PasswordAuthTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *PasswordAuthTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(
		native.NewStubConfig().
			WithNamespaceService().
			WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *PasswordAuthTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestPasswordAuthTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordAuthTestSuite))
}

func (s *PasswordAuthTestSuite) TestAuth() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	testNamespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        testNamespaceName,
		FullName:    tools.GetRandomString(10),
		Description: tools.GetRandomString(10),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: testNamespaceName})

	type policyData struct {
		namespace            string
		Actions              []string
		Resources            []string
		NamespaceIndependant bool
	}

	type roleData struct {
		namespace string
		policies  []int // index of the created scope
	}

	tests := []struct {
		testName string

		createNamespaces []string
		createPolicies   []policyData
		createRoles      []roleData

		assignedPolicies []int
		assignedRoles    []int

		requestedScopes []*auth.Scope
		hasAccess       bool
	}{
		{"Empty access", []string{}, []policyData{}, []roleData{}, []int{}, []int{}, []*auth.Scope{}, true},
		{"Root policy", []string{}, []policyData{{testNamespaceName, []string{"*"}, []string{"*"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple policy", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple policy wildcard 1", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.*"}, []string{"resource.345345.*"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple policy wildcard 2", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.password"}, []string{"resource.345345.*"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple policy wildcard 3", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.*"}, []string{"resource.345345.auth.password"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple policy but not assigned", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, false}}, []roleData{}, []int{}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, false},
		{"Simple policy but wrong namespace", []string{"simpleroletestnamespace"}, []policyData{{"simpleroletestnamespace", []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, false},
		{"Simple policy, wrong namespace but namespace independent", []string{"simpleroletestnamespace"}, []policyData{{"simpleroletestnamespace", []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, true}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{
			"Multiple policies from multiple namespaces",
			[]string{},
			[]policyData{
				{
					testNamespaceName,
					[]string{"action.123123.auth.1"},
					[]string{"resource.345345.auth.2"},
					false,
				},
				{
					"",
					[]string{"action.123123.auth.3"},
					[]string{"resource.345345.auth.4"},
					false,
				},
			},
			[]roleData{},
			[]int{0, 1},
			[]int{},
			[]*auth.Scope{
				{
					Namespace:            testNamespaceName,
					Resources:            []string{"resource.345345.auth.2"},
					Actions:              []string{"action.123123.auth.1"},
					NamespaceIndependent: false,
				},
				{
					Namespace:            "",
					Resources:            []string{"resource.345345.auth.4"},
					Actions:              []string{"action.123123.auth.3"},
					NamespaceIndependent: false,
				},
			},
			true,
		},
		{"Namespace independent root policy", []string{}, []policyData{{testNamespaceName, []string{"*"}, []string{"*"}, true}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            "",
				Resources:            []string{"*"},
				Actions:              []string{"*"},
				NamespaceIndependent: true,
			},
		}, true},
		{"Namespace independent root policy fail", []string{}, []policyData{{testNamespaceName, []string{"*"}, []string{"*"}, false}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"*"},
				Actions:              []string{"*"},
				NamespaceIndependent: true,
			},
		}, false},
		{"Namespace independent policy", []string{}, []policyData{{testNamespaceName, []string{"1234.*"}, []string{"1234.*"}, true}}, []roleData{}, []int{0}, []int{}, []*auth.Scope{
			{
				Namespace:            tools.GetRandomString(20),
				Resources:            []string{"1234.55"},
				Actions:              []string{"1234.66"},
				NamespaceIndependent: true,
			},
		}, true},
		{"Root role", []string{}, []policyData{{testNamespaceName, []string{"*"}, []string{"*"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple role", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple role wildcard 1", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.*"}, []string{"resource.345345.*"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple role wildcard 2", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.password"}, []string{"resource.345345.*"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple role wildcard 3", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.*"}, []string{"resource.345345.auth.password"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{"Simple role but not assigned", []string{}, []policyData{{testNamespaceName, []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, false},
		{"Simple role but wrong namespace", []string{"simpleroletestnamespace"}, []policyData{{"simpleroletestnamespace", []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, false},
		{"Simple role, wrong namespace but namespace independant", []string{"simpleroletestnamespace"}, []policyData{{"simpleroletestnamespace", []string{"action.123123.auth.password"}, []string{"resource.345345.auth.password"}, true}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"resource.345345.auth.password"},
				Actions:              []string{"action.123123.auth.password"},
				NamespaceIndependent: false,
			},
		}, true},
		{
			"Multiple roles from multiple namespaces",
			[]string{},
			[]policyData{
				{
					testNamespaceName,
					[]string{"action.123123.auth.1"},
					[]string{"resource.345345.auth.2"},
					false,
				},
				{
					"",
					[]string{"action.123123.auth.3"},
					[]string{"resource.345345.auth.4"},
					false,
				},
			},
			[]roleData{
				{
					testNamespaceName,
					[]int{0},
				},
				{
					"",
					[]int{1},
				},
			},
			[]int{},
			[]int{0, 1},
			[]*auth.Scope{
				{
					Namespace:            testNamespaceName,
					Resources:            []string{"resource.345345.auth.2"},
					Actions:              []string{"action.123123.auth.1"},
					NamespaceIndependent: false,
				},
				{
					Namespace:            "",
					Resources:            []string{"resource.345345.auth.4"},
					Actions:              []string{"action.123123.auth.3"},
					NamespaceIndependent: false,
				},
			},
			true,
		},
		{"Namespace independent root role", []string{}, []policyData{{testNamespaceName, []string{"*"}, []string{"*"}, true}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            "",
				Resources:            []string{"*"},
				Actions:              []string{"*"},
				NamespaceIndependent: true,
			},
		}, true},
		{"Namespace independent root role fail", []string{}, []policyData{{testNamespaceName, []string{"*"}, []string{"*"}, false}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            testNamespaceName,
				Resources:            []string{"*"},
				Actions:              []string{"*"},
				NamespaceIndependent: true,
			},
		}, false},
		{"Namespace independent role", []string{}, []policyData{{testNamespaceName, []string{"1234.*"}, []string{"1234.*"}, true}}, []roleData{{testNamespaceName, []int{0}}}, []int{}, []int{0}, []*auth.Scope{
			{
				Namespace:            tools.GetRandomString(20),
				Resources:            []string{"1234.55"},
				Actions:              []string{"1234.55"},
				NamespaceIndependent: true,
			},
		}, true},
	}

	for _, tc := range tests {
		s.Run(tc.testName, func() {
			// Create namespaces
			defer func() {
				for _, namespaceName := range tc.createNamespaces {
					s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})
				}
			}()
			for _, namespaceName := range tc.createNamespaces {
				_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
					Name:        namespaceName,
					FullName:    tools.GetRandomString(10),
					Description: tools.GetRandomString(10),
				})
				require.Nil(s.T(), err)
			}

			// Create policies
			createdPolicies := make([]struct {
				uuid      string
				namespace string
			}, 0, len(tc.assignedPolicies))
			defer func() {
				for _, createdPolicy := range createdPolicies {
					s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: createdPolicy.namespace, Uuid: createdPolicy.uuid})
				}
			}()
			for _, p := range tc.createPolicies {
				createResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
					Namespace:            p.namespace,
					Name:                 tools.GetRandomString(10),
					Description:          tools.GetRandomString(10),
					Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
					NamespaceIndependent: p.NamespaceIndependant,
					Resources:            p.Resources,
					Actions:              p.Actions,
				})
				require.Nil(s.T(), err)
				createdPolicies = append(createdPolicies, struct {
					uuid      string
					namespace string
				}{uuid: createResponse.Policy.Uuid, namespace: createResponse.Policy.Namespace})
			}

			// Create roles
			createdRoles := make([]struct {
				uuid      string
				namespace string
			}, 0, len(tc.assignedPolicies))
			defer func() {
				for _, createdRole := range createdRoles {
					s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: createdRole.namespace, Uuid: createdRole.uuid})
				}
			}()
			for _, r := range tc.createRoles {
				createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
					Namespace:   r.namespace,
					Name:        tools.GetRandomString(10),
					Description: tools.GetRandomString(10),
					Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
				})
				require.Nil(s.T(), err)
				createdRoles = append(createdRoles, struct {
					uuid      string
					namespace string
				}{uuid: createResponse.Role.Uuid, namespace: createResponse.Role.Namespace})

				for _, assignedPolicyIndex := range r.policies {
					_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
						RoleNamespace:   createResponse.Role.Namespace,
						RoleUUID:        createResponse.Role.Uuid,
						PolicyNamespace: createdPolicies[assignedPolicyIndex].namespace,
						PolicyUUID:      createdPolicies[assignedPolicyIndex].uuid,
					})
					require.Nil(s.T(), err)
				}
			}

			// Create Identity. Assign policies and roles
			identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
				Namespace:       testNamespaceName,
				Name:            tools.GetRandomString(10),
				InitiallyActive: true,
				Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
			})
			require.Nil(s.T(), err)

			pwd := tools.GetRandomString(20)
			_, err = s.nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
				Namespace: identityCreateResponse.Identity.Namespace,
				Identity:  identityCreateResponse.Identity.Uuid,
				Password:  pwd,
			})
			require.Nil(s.T(), err)

			for _, policyIndexToAssign := range tc.assignedPolicies {
				_, err := s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
					IdentityNamespace: identityCreateResponse.Identity.Namespace,
					IdentityUUID:      identityCreateResponse.Identity.Uuid,
					PolicyNamespace:   createdPolicies[policyIndexToAssign].namespace,
					PolicyUUID:        createdPolicies[policyIndexToAssign].uuid,
				})
				require.Nil(s.T(), err)
			}

			for _, roleIndexToAssign := range tc.assignedRoles {
				_, err := s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
					IdentityNamespace: identityCreateResponse.Identity.Namespace,
					IdentityUUID:      identityCreateResponse.Identity.Uuid,
					RoleNamespace:     createdRoles[roleIndexToAssign].namespace,
					RoleUUID:          createdRoles[roleIndexToAssign].uuid,
				})
				require.Nil(s.T(), err)
			}

			// Actually verify access
			accessResponse, err := s.nativeStub.Services.IAM.Auth.CheckAccessWithPassword(ctx, &auth.CheckAccessWithPasswordRequest{
				Namespace: identityCreateResponse.Identity.Namespace,
				Identity:  identityCreateResponse.Identity.Uuid,
				Password:  pwd,
				Metadata:  "{}",
				Scopes:    tc.requestedScopes,
			})
			require.Nil(s.T(), err)
			require.Equal(s.T(), tc.hasAccess, accessResponse.Status == auth.CheckAccessWithPasswordResponse_OK)
		})
	}
}
