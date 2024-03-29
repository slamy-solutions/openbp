package accesscontrol

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type IdentityRouter struct {
	nativeStub *native.NativeStub
}

type identityManagementInformation struct {
	ManagementType    string  `json:"type"`
	Reason            *string `json:"reason"`
	Service           *string `json:"service"`
	ManagementId      *string `json:"managementId"`
	IdentityNamespace *string `json:"identityNamespace"`
	IdentityUUID      *string `json:"identityUUID"`
}

type identityRole struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
}

type identityPolicy struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
}

type formatedIdentity struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`

	Managed identityManagementInformation `json:"managed"`

	Roles    []*identityRole   `json:"roles"`
	Policies []*identityPolicy `json:"policies"`

	Created string `json:"created"`
	Updated string `json:"updated"`
	Version uint64 `json:"version"`
}

func FormatIdentity(i *identity.Identity) *formatedIdentity {
	managed := identityManagementInformation{}
	if m, ok := i.Managed.(*identity.Identity_Service); ok {
		managed.ManagementType = "service"
		managed.Service = &m.Service.Service
		managed.Reason = &m.Service.Reason
		managed.ManagementId = &m.Service.ManagementId
	} else if m, ok := i.Managed.(*identity.Identity_Identity); ok {
		managed.ManagementType = "identity"
		managed.IdentityNamespace = &m.Identity.IdentityNamespace
		managed.IdentityUUID = &m.Identity.IdentityUUID
	} else {
		managed.ManagementType = "none"
	}

	roles := make([]*identityRole, 0, len(i.Roles))
	for _, role := range i.Roles {
		roles = append(roles, &identityRole{Namespace: role.Namespace, UUID: role.Uuid})
	}

	policies := make([]*identityPolicy, 0, len(i.Policies))
	for _, policy := range i.Policies {
		policies = append(policies, &identityPolicy{Namespace: policy.Namespace, UUID: policy.Uuid})
	}

	return &formatedIdentity{
		Namespace: i.Namespace,
		UUID:      i.Uuid,
		Name:      i.Name,
		Active:    i.Active,
		Managed:   managed,
		Roles:     roles,
		Policies:  policies,
		Created:   i.Created.AsTime().UTC().Format(time.RFC3339),
		Updated:   i.Updated.AsTime().UTC().Format(time.RFC3339),
		Version:   i.Version,
	}
}

type getListRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	Skip      uint32 `form:"skip" binding:"gte=0"`
	Limit     uint32 `form:"limit" binding:"lte=100,gte=0"`
}

type getListResponse struct {
	Identities []*formatedIdentity `json:"identities"`
	TotalCount uint64              `json:"totalCount"`
}

func (r *IdentityRouter) List(ctx *gin.Context) {
	var requestData getListRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity.*"},
			Actions:              []string{"native.iam.identity.get"},
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

	listStream, err := r.nativeStub.Services.IAM.Identity.List(ctx.Request.Context(), &identity.ListIdentityRequest{
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
	identities := make([]*formatedIdentity, 0, capacity)

	for {
		response, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		identities = append(identities, FormatIdentity(response.Identity))
	}

	countResponse, err := r.nativeStub.Services.IAM.Identity.Count(ctx.Request.Context(), &identity.CountIdentityRequest{
		Namespace: requestData.Namespace, UseCache: true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &getListResponse{Identities: identities, TotalCount: countResponse.Count})
}

type createRequest struct {
	Namespace       string `json:"namespace" binding:"lte=32"`
	Name            string `json:"name" binding:"lte=64"`
	InitiallyActive bool   `json:"initiallyActive" binding:"required"`
}

type createResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) Create(ctx *gin.Context) {
	var requestData createRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity"},
			Actions:              []string{"native.iam.identity.create"},
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

	response, err := r.nativeStub.Services.IAM.Identity.Create(ctx.Request.Context(), &identity.CreateIdentityRequest{
		Namespace:       requestData.Namespace,
		Name:            requestData.Name,
		InitiallyActive: requestData.InitiallyActive,
		Managed: &identity.CreateIdentityRequest_Identity{
			Identity: &identity.IdentityManagedData{
				IdentityNamespace: authData.Namespace,
				IdentityUUID:      authData.IdentityUUID,
			},
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &createResponse{Identity: FormatIdentity(response.Identity)})
}

type getRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}

type getResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) Get(ctx *gin.Context) {
	var requestData getRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity." + requestData.UUID},
			Actions:              []string{"native.iam.identity.get"},
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

	// Get identity
	identityResponse, err := r.nativeStub.Services.IAM.Identity.Get(ctx.Request.Context(), &identity.GetIdentityRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		UseCache:  true,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity UUID or Namespace arguments."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &getResponse{Identity: FormatIdentity(identityResponse.Identity)})
}

type deleteRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}

func (r *IdentityRouter) Delete(ctx *gin.Context) {
	var requestData deleteRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity." + requestData.UUID},
			Actions:              []string{"native.iam.identity.delete"},
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

	// Delete identity
	_, err = r.nativeStub.Services.IAM.Identity.Delete(ctx.Request.Context(), &identity.DeleteIdentityRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, gin.H{})
}

type setActiveIdentityRequest struct {
	Namespace string `json:"namespace" binding:"lte=32"`
	UUID      string `json:"uuid" binding:"required,lte=64"`
	Active    bool   `json:"active"`
}

type setActiveIdentityResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) SetActive(ctx *gin.Context) {
	var requestData setActiveIdentityRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	requiredAction := "native.iam.identity.disable"
	if !requestData.Active {
		requiredAction = "native.iam.identity.enable"
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity." + requestData.UUID},
			Actions:              []string{requiredAction},
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

	re, err := r.nativeStub.Services.IAM.Identity.SetActive(ctx.Request.Context(), &identity.SetIdentityActiveRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		Active:    requestData.Active,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity UUID."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, setActiveIdentityResponse{Identity: FormatIdentity(re.Identity)})
}

type updateIdentityRequest struct {
	Namespace string `json:"namespace" binding:"lte=32"`
	UUID      string `json:"uuid" binding:"required,lte=64"`
	NewName   string `json:"newName" binding:"required,lte=64"`
}

type updateIdentityResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) Update(ctx *gin.Context) {
	var requestData updateIdentityRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity." + requestData.UUID},
			Actions:              []string{"native.iam.identity.update"},
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

	updateResponse, err := r.nativeStub.Services.IAM.Identity.Update(ctx.Request.Context(), &identity.UpdateIdentityRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		NewName:   requestData.NewName,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity UUID."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &updateIdentityResponse{Identity: FormatIdentity(updateResponse.Identity)})
}

type addPolicyToIdentityRequest struct {
	PolicyNamespace   string `json:"policyNamespace" binding:"lte=32"`
	PolicyUUID        string `json:"policyUUID" binding:"required,lte=64"`
	IdentityNamespace string `json:"identityNamespace" binding:"lte=64"`
	IdentityUUID      string `json:"identityUUID" binding:"required,lte=64"`
}

type addPolicyToIdentityResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) AddPolicy(ctx *gin.Context) {
	var requestData addPolicyToIdentityRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	scopes := []*auth.Scope{
		{
			Namespace:            requestData.IdentityNamespace,
			Resources:            []string{"native.iam.role." + requestData.IdentityUUID},
			Actions:              []string{"native.iam.role.update"},
			NamespaceIndependent: false,
		},
	}
	if requestData.PolicyNamespace != requestData.IdentityNamespace {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.PolicyNamespace,
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

	addPolicyResponse, err := r.nativeStub.Services.IAM.Identity.AddPolicy(ctx.Request.Context(), &identity.AddPolicyRequest{
		IdentityNamespace: requestData.IdentityNamespace,
		IdentityUUID:      requestData.IdentityUUID,
		PolicyNamespace:   requestData.PolicyNamespace,
		PolicyUUID:        requestData.PolicyUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity or policy UUID."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity or policy not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &addPolicyToIdentityResponse{Identity: FormatIdentity(addPolicyResponse.Identity)})
}

type removePolicyFromIdentityRequest struct {
	PolicyNamespace   string `json:"policyNamespace" binding:"lte=32"`
	PolicyUUID        string `json:"policyUUID" binding:"required,lte=64"`
	IdentityNamespace string `json:"identityNamespace" binding:"lte=64"`
	IdentityUUID      string `json:"identityUUID" binding:"required,lte=64"`
}

type removePolicyFromIdentityResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) RemovePolicy(ctx *gin.Context) {
	var requestData removePolicyFromIdentityRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	scopes := []*auth.Scope{
		{
			Namespace:            requestData.IdentityNamespace,
			Resources:            []string{"native.iam.role." + requestData.IdentityUUID},
			Actions:              []string{"native.iam.role.update"},
			NamespaceIndependent: false,
		},
	}
	if requestData.PolicyNamespace != requestData.IdentityNamespace {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.PolicyNamespace,
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

	removePolicyResponse, err := r.nativeStub.Services.IAM.Identity.RemovePolicy(ctx.Request.Context(), &identity.RemovePolicyRequest{
		IdentityNamespace: requestData.IdentityNamespace,
		IdentityUUID:      requestData.IdentityUUID,
		PolicyNamespace:   requestData.PolicyNamespace,
		PolicyUUID:        requestData.PolicyUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity or policy UUID."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity or policy not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &removePolicyFromIdentityResponse{Identity: FormatIdentity(removePolicyResponse.Identity)})
}

type addRoleToIdentityRequest struct {
	RoleNamespace     string `json:"roleNamespace" binding:"lte=32"`
	RoleUUID          string `json:"roleUUID" binding:"required,lte=64"`
	IdentityNamespace string `json:"identityNamespace" binding:"lte=64"`
	IdentityUUID      string `json:"identityUUID" binding:"required,lte=64"`
}

type addRoleToIdentityResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) AddRole(ctx *gin.Context) {
	var requestData addRoleToIdentityRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	scopes := []*auth.Scope{
		{
			Namespace:            requestData.IdentityNamespace,
			Resources:            []string{"native.iam.role." + requestData.IdentityUUID},
			Actions:              []string{"native.iam.role.update"},
			NamespaceIndependent: false,
		},
	}
	if requestData.RoleNamespace != requestData.IdentityNamespace {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.RoleNamespace,
			Resources:            []string{"native.iam.policy." + requestData.RoleUUID},
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

	addRoleyResponse, err := r.nativeStub.Services.IAM.Identity.AddRole(ctx.Request.Context(), &identity.AddRoleRequest{
		IdentityNamespace: requestData.IdentityNamespace,
		IdentityUUID:      requestData.IdentityUUID,
		RoleNamespace:     requestData.RoleNamespace,
		RoleUUID:          requestData.RoleUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity or role UUID."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity or role not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &addRoleToIdentityResponse{Identity: FormatIdentity(addRoleyResponse.Identity)})
}

type removeRoleFromIdentityRequest struct {
	RoleNamespace     string `json:"roleNamespace" binding:"lte=32"`
	RoleUUID          string `json:"roleUUID" binding:"required,lte=64"`
	IdentityNamespace string `json:"identityNamespace" binding:"lte=64"`
	IdentityUUID      string `json:"identityUUID" binding:"required,lte=64"`
}

type removeRoleFromIdentityResponse struct {
	Identity *formatedIdentity `json:"identity"`
}

func (r *IdentityRouter) RemoveRole(ctx *gin.Context) {
	var requestData removeRoleFromIdentityRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	scopes := []*auth.Scope{
		{
			Namespace:            requestData.IdentityNamespace,
			Resources:            []string{"native.iam.role." + requestData.IdentityUUID},
			Actions:              []string{"native.iam.role.update"},
			NamespaceIndependent: false,
		},
	}
	if requestData.RoleNamespace != requestData.IdentityNamespace {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.RoleNamespace,
			Resources:            []string{"native.iam.policy." + requestData.RoleUUID},
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

	removeRoleyResponse, err := r.nativeStub.Services.IAM.Identity.RemoveRole(ctx.Request.Context(), &identity.RemoveRoleRequest{
		IdentityNamespace: requestData.IdentityNamespace,
		IdentityUUID:      requestData.IdentityUUID,
		RoleNamespace:     requestData.RoleNamespace,
		RoleUUID:          requestData.RoleUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(422, gin.H{"message": "Invalid identity or role UUID."})
				return
			}
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Identity or role not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &removeRoleFromIdentityResponse{Identity: FormatIdentity(removeRoleyResponse.Identity)})
}
