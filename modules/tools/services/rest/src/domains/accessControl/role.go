package accesscontrol

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoleRouter struct {
	nativeStub *native.NativeStub
}

type roleManagementInformation struct {
	ManagementType    string  `json:"type"`
	Reason            *string `json:"reason"`
	Service           *string `json:"service"`
	ManagementId      *string `json:"managementId"`
	IdentityNamespace *string `json:"identityNamespace"`
	IdentityUUID      *string `json:"identityUUID"`
}

type rolePolicy struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
}

type formatedRole struct {
	Namespace   string `json:"namespace"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Managed roleManagementInformation `json:"managed"`

	Policies []rolePolicy `json:"policies"`

	Tags    []string `json:"tags"`
	Created string   `json:"created"`
	Updated string   `json:"updated"`
	Version uint64   `json:"version"`
}

func FormatRole(i *role.Role) *formatedRole {
	managed := roleManagementInformation{}
	if m, ok := i.Managed.(*role.Role_Service); ok {
		managed.ManagementType = "service"
		managed.Service = &m.Service.Service
		managed.Reason = &m.Service.Reason
		managed.ManagementId = &m.Service.ManagementId
	} else if m, ok := i.Managed.(*role.Role_Identity); ok {
		managed.ManagementType = "identity"
		managed.IdentityNamespace = &m.Identity.IdentityNamespace
		managed.IdentityUUID = &m.Identity.IdentityUUID
	} else if _, ok := i.Managed.(*role.Role_BuiltIn); ok {
		managed.ManagementType = "builtIn"
	} else {
		managed.ManagementType = "none"
	}

	policies := make([]rolePolicy, 0, len(i.Policies))
	for _, p := range i.Policies {
		policies = append(policies, rolePolicy{Namespace: p.Namespace, UUID: p.Uuid})
	}

	return &formatedRole{
		Namespace:   i.Namespace,
		UUID:        i.Uuid,
		Name:        i.Name,
		Managed:     managed,
		Description: i.Description,
		Policies:    policies,
		Tags:        i.Tags,
		Created:     i.Created.AsTime().UTC().Format(time.RFC3339),
		Updated:     i.Updated.AsTime().UTC().Format(time.RFC3339),
		Version:     i.Version,
	}
}

type getRoleListRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	Skip      uint32 `form:"skip" binding:"gte=0"`
	Limit     uint32 `form:"limit" binding:"lte=100,gte=0"`
}

type getRoleListResponse struct {
	Roles      []*formatedRole `json:"roles"`
	TotalCount uint64          `json:"totalCount"`
}

func (r *RoleRouter) List(ctx *gin.Context) {
	var requestData getRoleListRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.role.*"},
			Actions:              []string{"native.iam.role.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	listStream, err := r.nativeStub.Services.IamRole.List(ctx.Request.Context(), &role.ListRolesRequest{
		Namespace: requestData.Namespace,
		Skip:      uint64(requestData.Skip),
		Limit:     uint64(requestData.Limit),
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	capacity := 100
	if requestData.Limit != 0 {
		capacity = int(requestData.Limit)
	}
	roles := make([]*formatedRole, 0, capacity)

	for {
		response, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		roles = append(roles, FormatRole(response.Role))
	}

	countResponse, err := r.nativeStub.Services.IamRole.Count(ctx.Request.Context(), &role.CountRolesRequest{
		Namespace: requestData.Namespace, UseCache: true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &getRoleListResponse{Roles: roles, TotalCount: countResponse.Count})
}

type createRoleRequest struct {
	Namespace   string `json:"namespace" binding:"lte=32"`
	Name        string `json:"name" binding:"required,lte=64"`
	Description string `json:"description" binding:"lte=256"`
}

type createRoleResponse struct {
	Role *formatedRole `json:"role"`
}

func (r *RoleRouter) Create(ctx *gin.Context) {
	var requestData createRoleRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.role"},
			Actions:              []string{"native.iam.role.create"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	response, err := r.nativeStub.Services.IamRole.Create(ctx.Request.Context(), &role.CreateRoleRequest{
		Namespace:   requestData.Namespace,
		Name:        requestData.Name,
		Description: requestData.Description,
		Managed: &role.CreateRoleRequest_Identity{Identity: &role.IdentityManagedData{
			IdentityNamespace: authData.Namespace,
			IdentityUUID:      authData.IdentityUUID,
		}},
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Role name/resource/action has bad format"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &createRoleResponse{Role: FormatRole(response.Role)})
}

type getRoleRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}

type getRoleResponse struct {
	Role *formatedRole `json:"role"`
}

func (r *RoleRouter) Get(ctx *gin.Context) {
	var requestData getRoleRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.role." + requestData.UUID},
			Actions:              []string{"native.iam.role.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	// Delete policy
	response, err := r.nativeStub.Services.IamRole.Get(ctx.Request.Context(), &role.GetRoleRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		UseCache:  true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	etag := strconv.FormatUint(response.Role.Version, 10)
	if etag == ctx.Request.Header.Get("If-None-Match") {
		ctx.Status(http.StatusNotModified)
		return
	}

	ctx.Header("ETag", etag)
	ctx.Header("Cache-Control", "private")

	ctx.JSON(http.StatusOK, getRoleResponse{Role: FormatRole(response.Role)})
}

/*
type updateRoleRequest struct {
	Namespace            string   `json:"namespace" binding:"lte=32"`
	UUID                 string   `json:"uuid" binding:"required,lte=64"`
	Name                 string   `json:"name" binding:"required,lte=64"`
	Description          string   `json:"description" binding:"lte=256"`
}

type updateRoleResponse struct {
	Policy *formatedPolicy `json:"policy"`
}

func (r *RoleRouter) Update(ctx *gin.Context) {
	var requestData updatePolicyRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.role." + requestData.UUID},
			Actions:              []string{"native.iam.role.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	r.nativeStub.Services.IamRole.
}
*/

type addPolicyToRoleRequest struct {
	RoleNamespace  string `json:"roleNamespace" binding:"lte=32"`
	RoleUUID       string `json:"roleUUID" binding:"required,lte=64"`
	Policyamespace string `json:"policyNamespace" binding:"lte=32"`
	PolicyUUID     string `json:"policyUUID" binding:"required,lte=64"`
}

type addPolicyToResponse struct {
	Role *formatedRole `json:"role"`
}

func (r *RoleRouter) AddPolicy(ctx *gin.Context) {
	var requestData addPolicyToRoleRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	scopes := []*auth.Scope{
		{
			Namespace:            requestData.RoleNamespace,
			Resources:            []string{"native.iam.role." + requestData.RoleUUID},
			Actions:              []string{"native.iam.role.policy"},
			NamespaceIndependent: false,
		},
	}
	if requestData.RoleNamespace != requestData.Policyamespace {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.Policyamespace,
			Resources:            []string{"native.iam.policy." + requestData.PolicyUUID},
			Actions:              []string{"native.iam.policy.useInOtherNamespace"},
			NamespaceIndependent: false,
		})
	}

	authData, err := authTools.CheckAuth(ctx, r.nativeStub, scopes)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	response, err := r.nativeStub.Services.IamRole.AddPolicy(ctx.Request.Context(), &role.AddPolicyRequest{
		RoleNamespace:   requestData.RoleNamespace,
		RoleUUID:        requestData.RoleUUID,
		PolicyNamespace: requestData.Policyamespace,
		PolicyUUID:      requestData.PolicyUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Role not found"})
				return
			}
			if st.Code() == codes.FailedPrecondition {
				ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "Policy doesnt exist"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, addPolicyToResponse{Role: FormatRole(response.Role)})
}

type removePolicyFromRoleRequest struct {
	RoleNamespace  string `json:"roleNamespace" binding:"lte=32"`
	RoleUUID       string `json:"roleUUID" binding:"required,lte=64"`
	Policyamespace string `json:"policyNamespace" binding:"lte=32"`
	PolicyUUID     string `json:"policyUUID" binding:"required,lte=64"`
}

type removePolicyFromResponse struct {
	Role *formatedRole `json:"role"`
}

func (r *RoleRouter) RemovePolicy(ctx *gin.Context) {
	var requestData removePolicyFromRoleRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	scopes := []*auth.Scope{
		{
			Namespace:            requestData.RoleNamespace,
			Resources:            []string{"native.iam.role." + requestData.RoleUUID},
			Actions:              []string{"native.iam.role.policy"},
			NamespaceIndependent: false,
		},
	}
	if requestData.RoleNamespace != requestData.Policyamespace {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.Policyamespace,
			Resources:            []string{"native.iam.policy." + requestData.PolicyUUID},
			Actions:              []string{"native.iam.policy.useInOtherNamespace"},
			NamespaceIndependent: false,
		})
	}

	authData, err := authTools.CheckAuth(ctx, r.nativeStub, scopes)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	response, err := r.nativeStub.Services.IamRole.RemovePolicy(ctx.Request.Context(), &role.RemovePolicyRequest{
		RoleNamespace:   requestData.RoleNamespace,
		RoleUUID:        requestData.RoleUUID,
		PolicyNamespace: requestData.Policyamespace,
		PolicyUUID:      requestData.PolicyUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Role not found"})
				return
			}
			if st.Code() == codes.FailedPrecondition {
				ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "Policy doesnt exist"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, removePolicyFromResponse{Role: FormatRole(response.Role)})
}

type deleteRoleRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}

func (r *RoleRouter) Delete(ctx *gin.Context) {
	var requestData deleteRoleRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.role." + requestData.UUID},
			Actions:              []string{"native.iam.role.delete"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	// Delete policy
	_, err = r.nativeStub.Services.IamRole.Delete(ctx.Request.Context(), &role.DeleteRoleRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, gin.H{})
}
