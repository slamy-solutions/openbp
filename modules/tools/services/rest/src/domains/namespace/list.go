package namespace

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	nativeNamespace "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type ListRouter struct {
	nativeStub *native.NativeStub
}

type getListResponse struct {
	Namespaces []*nativeNamespace.Namespace `json:"namespaces"`
}

func (r *ListRouter) List(ctx *gin.Context) {
	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"native.namespace"},
			Actions:              []string{"native.namespace.get"},
			NamespaceIndependent: true,
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

	response := getListResponse{
		Namespaces: []*nativeNamespace.Namespace{},
	}

	namespaceClient, err := r.nativeStub.Services.Namespace.GetAll(ctx.Request.Context(), &nativeNamespace.GetAllNamespacesRequest{UseCache: true})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer namespaceClient.CloseSend()

	for {
		n, err := namespaceClient.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		response.Namespaces = append(response.Namespaces, n.Namespace)
	}

	ctx.JSON(http.StatusOK, response)
}

type createRequest struct {
	Name        string `json:"name" binding:"required,lte=32,gt=0"`
	FullName    string `json:"fullName" binding:"required,lte=128"`
	Description string `json:"description" binding:"required,lte=512"`
}
type createResponse struct {
	Namespace *nativeNamespace.Namespace `json:"namespace"`
}

func (r *ListRouter) Create(ctx *gin.Context) {
	var requestData createRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{fmt.Sprintf("native.namespace.%s", requestData.Name)},
			Actions:              []string{"native.namespace.create"},
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

	// Create new namespace
	response, err := r.nativeStub.Services.Namespace.Create(ctx.Request.Context(), &nativeNamespace.CreateNamespaceRequest{
		Name:        requestData.Name,
		FullName:    requestData.FullName,
		Description: requestData.Description,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid namespace name"})
				return
			}

			if st.Code() == codes.AlreadyExists {
				ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Namespace with this name already exists"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &createResponse{Namespace: response.Namespace})
}

type deleteRequest struct {
	Name string `form:"name" binding:"required,lte=32,gt=0"`
}

func (r *ListRouter) Delete(ctx *gin.Context) {
	var requestData deleteRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{fmt.Sprintf("native.namespace.%s", requestData.Name)},
			Actions:              []string{"native.namespace.delete"},
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

	// Delete namespace
	_, err = r.nativeStub.Services.Namespace.Delete(ctx.Request.Context(), &nativeNamespace.DeleteNamespaceRequest{
		Name: requestData.Name,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
