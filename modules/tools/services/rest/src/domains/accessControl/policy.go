package accesscontrol

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type PolicyRouter struct {
	nativeStub *native.NativeStub
}

type policyManagementInformation struct {
	ManagementType    string  `json:"type"`
	Reason            *string `json:"reason"`
	Service           *string `json:"service"`
	ManagementId      *string `json:"managementId"`
	IdentityNamespace *string `json:"identityNamespace"`
	IdentityUUID      *string `json:"identityUUID"`
}

type formatedPolicy struct {
	Namespace   string `json:"namespace"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Managed policyManagementInformation `json:"managed"`

	Resources            []string `json:"resources"`
	Actions              []string `json:"actions"`
	NamespaceIndependent bool     `json:"namespaceIndependent"`

	Tags    []string `json:"tags"`
	Created string   `json:"created"`
	Updated string   `json:"updated"`
	Version uint64   `json:"version"`
}

func FormatPolicy(i *policy.Policy) *formatedPolicy {
	managed := policyManagementInformation{}
	if m, ok := i.Managed.(*policy.Policy_Service); ok {
		managed.ManagementType = "service"
		managed.Service = &m.Service.Service
		managed.Reason = &m.Service.Reason
		managed.ManagementId = &m.Service.ManagementId
	} else if m, ok := i.Managed.(*policy.Policy_Identity); ok {
		managed.ManagementType = "identity"
		managed.IdentityNamespace = &m.Identity.IdentityNamespace
		managed.IdentityUUID = &m.Identity.IdentityUUID
	} else if _, ok := i.Managed.(*policy.Policy_BuiltIn); ok {
		managed.ManagementType = "builtIn"
	} else {
		managed.ManagementType = "none"
	}

	resources := make([]string, 0)
	if i.Resources != nil && len(i.Resources) > 0 {
		resources = i.Resources
	}
	actions := make([]string, 0)
	if i.Actions != nil && len(i.Actions) > 0 {
		actions = i.Actions
	}
	tags := make([]string, 0)
	if i.Tags != nil && len(i.Tags) > 0 {
		tags = i.Tags
	}

	return &formatedPolicy{
		Namespace:            i.Namespace,
		UUID:                 i.Uuid,
		Name:                 i.Name,
		Managed:              managed,
		Description:          i.Description,
		Resources:            resources,
		Actions:              actions,
		NamespaceIndependent: i.NamespaceIndependent,
		Tags:                 tags,
		Created:              i.Created.AsTime().UTC().Format(time.RFC3339),
		Updated:              i.Updated.AsTime().UTC().Format(time.RFC3339),
		Version:              i.Version,
	}
}

type getPolicyListRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	Skip      uint32 `form:"skip" binding:"gte=0"`
	Limit     uint32 `form:"limit" binding:"lte=100,gte=0"`
}

type getPolicyListResponse struct {
	Policies   []*formatedPolicy `json:"policies"`
	TotalCount uint64            `json:"totalCount"`
}

func (r *PolicyRouter) List(ctx *gin.Context) {
	var requestData getPolicyListRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.policy.*"},
			Actions:              []string{"native.iam.policy.get"},
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

	listStream, err := r.nativeStub.Services.IamPolicy.List(ctx.Request.Context(), &policy.ListPoliciesRequest{
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
	policies := make([]*formatedPolicy, 0, capacity)

	for {
		response, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		policies = append(policies, FormatPolicy(response.Policy))
	}

	countResponse, err := r.nativeStub.Services.IamPolicy.Count(ctx.Request.Context(), &policy.CountPoliciesRequest{
		Namespace: requestData.Namespace, UseCache: true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &getPolicyListResponse{Policies: policies, TotalCount: countResponse.Count})
}

type createPolicyRequest struct {
	Namespace            string   `json:"namespace" binding:"lte=32"`
	Name                 string   `json:"name" binding:"required,lte=64"`
	Description          string   `json:"description" binding:"lte=256"`
	Resources            []string `json:"resources" binding:"lte=64"`
	Actions              []string `json:"actions" binding:"lte=64"`
	NamespaceIndependent bool     `json:"namespaceIndependent" binding:""`
}

type createPolicyResponse struct {
	Policy *formatedPolicy `json:"policy"`
}

func (r *PolicyRouter) Create(ctx *gin.Context) {
	var requestData createPolicyRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.policy"},
			Actions:              []string{"native.iam.policy.create"},
			NamespaceIndependent: requestData.NamespaceIndependent, // If user will create namespace independant policy, then this user must be an administrator with full namespace-independant access.
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

	response, err := r.nativeStub.Services.IamPolicy.Create(ctx.Request.Context(), &policy.CreatePolicyRequest{
		Namespace:   requestData.Namespace,
		Name:        requestData.Name,
		Description: requestData.Description,
		Managed: &policy.CreatePolicyRequest_Identity{Identity: &policy.IdentityManagedData{
			IdentityNamespace: authData.Namespace,
			IdentityUUID:      authData.IdentityUUID,
		}},
		NamespaceIndependent: requestData.NamespaceIndependent,
		Resources:            requestData.Resources,
		Actions:              requestData.Actions,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Policy name/resource/action has bad format"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &createPolicyResponse{Policy: FormatPolicy(response.Policy)})
}

type deletePolicyRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}

func (r *PolicyRouter) Delete(ctx *gin.Context) {
	var requestData deletePolicyRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.policy." + requestData.UUID},
			Actions:              []string{"native.iam.policy.delete"},
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
	_, err = r.nativeStub.Services.IamPolicy.Delete(ctx.Request.Context(), &policy.DeletePolicyRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, gin.H{})
}

type getPolicyRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}

type getPolicyResponse struct {
	Policy *formatedPolicy `json:"policy"`
}

func (r *PolicyRouter) Get(ctx *gin.Context) {
	var requestData getPolicyRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.policy." + requestData.UUID},
			Actions:              []string{"native.iam.policy.get"},
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
	response, err := r.nativeStub.Services.IamPolicy.Get(ctx.Request.Context(), &policy.GetPolicyRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		UseCache:  true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	etag := strconv.FormatUint(response.Policy.Version, 10)
	if etag == ctx.Request.Header.Get("If-None-Match") {
		ctx.Status(http.StatusNotModified)
		return
	}

	ctx.Header("ETag", etag)
	ctx.Header("Cache-Control", "private")

	ctx.JSON(http.StatusOK, getPolicyResponse{Policy: FormatPolicy(response.Policy)})
}

type updatePolicyRequest struct {
	Namespace            string   `json:"namespace" binding:"lte=32"`
	UUID                 string   `json:"uuid" binding:"required,lte=64"`
	Name                 string   `json:"name" binding:"required,lte=64"`
	Description          string   `json:"description" binding:"lte=256"`
	Resources            []string `json:"resources" binding:"lte=64"`
	Actions              []string `json:"actions" binding:"lte=64"`
	NamespaceIndependent bool     `json:"namespaceIndependent"`
}

type updatePolicyResponse struct {
	Policy *formatedPolicy `json:"policy"`
}

func (r *PolicyRouter) Update(ctx *gin.Context) {
	var requestData updatePolicyRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.policy." + requestData.UUID},
			Actions:              []string{"native.iam.policy.update"},
			NamespaceIndependent: requestData.NamespaceIndependent,
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

	// Update policy
	response, err := r.nativeStub.Services.IamPolicy.Update(ctx.Request.Context(), &policy.UpdatePolicyRequest{
		Namespace:            requestData.Namespace,
		Uuid:                 requestData.UUID,
		Name:                 requestData.Name,
		Description:          requestData.Description,
		NamespaceIndependent: requestData.NamespaceIndependent,
		Resources:            requestData.Resources,
		Actions:              requestData.Actions,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Policy name/resource/action has bad format"})
				return
			}

			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Policy not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, &updatePolicyResponse{Policy: FormatPolicy(response.Policy)})
}
