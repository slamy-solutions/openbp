package crm

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	client "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/client"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type clientRouter struct {
	crmStub    *crm.CRMStub
	nativeStub *native.NativeStub
}

type formatedContactPerson struct {
	Namespace   string   `json:"namespace"`
	UUID        string   `json:"uuid"`
	ClientUUID  string   `json:"clientUUID"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Phone       []string `json:"phone"`
	NotRelevant bool     `json:"notRelevant"`
	Comment     string   `json:"comment"`
}

func formatedContactPersonFromGRPC(c *client.ContactPerson) formatedContactPerson {
	return formatedContactPerson{
		Namespace:   c.Namespace,
		UUID:        c.Uuid,
		ClientUUID:  c.ClientUUID,
		Name:        c.Name,
		Email:       c.Email,
		Phone:       c.Phone,
		NotRelevant: c.NotRelevant,
		Comment:     c.Comment,
	}
}

type formatedClient struct {
	Namespace      string                  `json:"namespace"`
	UUID           string                  `json:"uuid"`
	Name           string                  `json:"name"`
	ContactPersons []formatedContactPerson `json:"contactPersons"`
	CreatedAt      time.Time               `json:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt"`
	Version        int64                   `json:"version"`
}

func formatedClientFromGRPC(c *client.Client) *formatedClient {
	var contactPersons []formatedContactPerson
	for _, cp := range c.ContactPersons {
		contactPersons = append(contactPersons, formatedContactPersonFromGRPC(cp))
	}

	return &formatedClient{
		Namespace:      c.Namespace,
		UUID:           c.Uuid,
		Name:           c.Name,
		ContactPersons: contactPersons,
		CreatedAt:      c.CreatedAt.AsTime(),
		UpdatedAt:      c.UpdatedAt.AsTime(),
		Version:        c.Version,
	}
}

type getAllClientsRequest struct {
	Namespace string `form:"namespace"`
}
type getAllClientsResponse struct {
	Clients []formatedClient `json:"clients"`
}

func (r *clientRouter) GetAllClients(ctx *gin.Context) {
	var requestData getAllClientsRequest
	if err := ctx.ShouldBindQuery(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.get"},
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

	response, err := r.crmStub.Core.Client.GetAll(ctx.Request.Context(), &client.GetAllRequest{
		Namespace: requestData.Namespace,
		UseCache:  true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var clients []formatedClient
	for _, c := range response.Clients {
		clients = append(clients, *formatedClientFromGRPC(c))
	}

	ctx.JSON(http.StatusOK, response)
}

type getClientRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid"`
}
type getClientResponse struct {
	Client *formatedClient `json:"client"`
}

func (r *clientRouter) GetClient(ctx *gin.Context) {
	var requestData getClientRequest
	if err := ctx.ShouldBindQuery(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.get"},
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

	response, err := r.crmStub.Core.Client.Get(ctx.Request.Context(), &client.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		UseCache:  true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getClientResponse{
		Client: formatedClientFromGRPC(response.Client),
	})
}

type createClientRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
}
type createClientResponse struct {
	Client *formatedClient `json:"client"`
}

func (r *clientRouter) CreateClient(ctx *gin.Context) {
	var requestData createClientRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.create"},
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

	response, err := r.crmStub.Core.Client.Create(ctx.Request.Context(), &client.CreateRequest{
		Namespace: requestData.Namespace,
		Name:      requestData.Name,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, createClientResponse{
		Client: formatedClientFromGRPC(response.Client),
	})
}

type updateClientRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	UUID      string `json:"uuid" binding:"required"`
	Name      string `json:"name" binding:"required"`
}
type updateClientResponse struct {
	Client *formatedClient `json:"client"`
}

func (r *clientRouter) UpdateClient(ctx *gin.Context) {
	var requestData updateClientRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.update"},
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

	response, err := r.crmStub.Core.Client.Update(ctx.Request.Context(), &client.UpdateRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		Name:      requestData.Name,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updateClientResponse{
		Client: formatedClientFromGRPC(response.Client),
	})
}

type deleteClientRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	UUID      string `json:"uuid" binding:"required"`
}
type deleteClientResponse struct{}

func (r *clientRouter) DeleteClient(ctx *gin.Context) {
	var requestData deleteClientRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.delete"},
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

	_, err = r.crmStub.Core.Client.Delete(ctx.Request.Context(), &client.DeleteRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, deleteClientResponse{})
}

type addContactPersonRequest struct {
	Namespace  string   `json:"namespace"`
	ClientUUID string   `json:"clientUUID" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Email      string   `json:"email"`
	Phone      []string `json:"phone"`
	Comment    string   `json:"comment"`
}
type addContactPersonResponse struct {
	ContactPerson formatedContactPerson `json:"contactPerson"`
}

func (r *clientRouter) AddContactPerson(ctx *gin.Context) {
	var requestData addContactPersonRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.contactPerson.add"},
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

	response, err := r.crmStub.Core.Client.AddContactPerson(ctx.Request.Context(), &client.AddContactPersonRequest{
		Namespace:  requestData.Namespace,
		ClientUUID: requestData.ClientUUID,
		Name:       requestData.Name,
		Email:      requestData.Email,
		Phone:      requestData.Phone,
		Comment:    requestData.Comment,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, addContactPersonResponse{
		ContactPerson: formatedContactPersonFromGRPC(response.ContactPerson),
	})
}

type updateContactPersonRequest struct {
	Namespace   string   `json:"namespace"`
	UUID        string   `json:"uuid" binding:"required"`
	ClientUUID  string   `json:"clientUUID" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Email       string   `json:"email"`
	Phone       []string `json:"phone"`
	NotRelevant bool     `json:"notRelevant"`
	Comment     string   `json:"comment"`
}
type updateContactPersonResponse struct {
	ContactPerson formatedContactPerson `json:"contactPerson"`
}

func (r *clientRouter) UpdateContactPerson(ctx *gin.Context) {
	var requestData updateContactPersonRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.contactPerson.update"},
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

	response, err := r.crmStub.Core.Client.UpdateContactPerson(ctx.Request.Context(), &client.UpdateContactPersonRequest{
		Namespace:         requestData.Namespace,
		ContactPersonUUID: requestData.UUID,
		ClientUUID:        requestData.ClientUUID,
		Name:              requestData.Name,
		Email:             requestData.Email,
		Phone:             requestData.Phone,
		NotRelevant:       requestData.NotRelevant,
		Comment:           requestData.Comment,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updateContactPersonResponse{
		ContactPerson: formatedContactPersonFromGRPC(response.ContactPerson),
	})
}

type deleteContactPersonRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid" binding:"required"`
}
type deleteContactPersonResponse struct{}

func (r *clientRouter) DeleteContactPerson(ctx *gin.Context) {
	var requestData deleteContactPersonRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.contactPerson.delete"},
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

	_, err = r.crmStub.Core.Client.DeleteContactPerson(ctx.Request.Context(), &client.DeleteContactPersonRequest{
		Namespace:         requestData.Namespace,
		ContactPersonUUID: requestData.UUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, deleteContactPersonResponse{})
}

type getContactPersonsForClientRequest struct {
	Namespace  string `form:"namespace"`
	ClientUUID string `form:"clientUUID" binding:"required"`
}
type getContactPersonsForClientResponse struct {
	ContactPersons []formatedContactPerson `json:"contactPersons"`
}

func (r *clientRouter) GetContactPersonsForClient(ctx *gin.Context) {
	var requestData getContactPersonsForClientRequest
	if err := ctx.ShouldBindQuery(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.client"},
			Actions:              []string{"crm.client.contactPerson.get"},
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

	response, err := r.crmStub.Core.Client.GetContactPersonsForClient(ctx.Request.Context(), &client.GetContactPersonsForClientRequest{
		Namespace:  requestData.Namespace,
		ClientUUID: requestData.ClientUUID,
		UseCache:   true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var contactPersons []formatedContactPerson
	for _, cp := range response.ContactPersons {
		contactPersons = append(contactPersons, formatedContactPersonFromGRPC(cp))
	}

	ctx.JSON(http.StatusOK, getContactPersonsForClientResponse{
		ContactPersons: contactPersons,
	})
}
