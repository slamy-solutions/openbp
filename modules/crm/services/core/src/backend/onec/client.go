package onec

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
)

type clientRepository struct {
	logger    *slog.Logger
	connector *connector.OneCConnector
	namespace string
}

type client struct {
	UUID         string    `json:"id"`
	Name         string    `json:"name"`
	ActivityDate time.Time `json:"activityDate"`
	Telephone    string    `json:"telephone"`
}

type clientContact struct {
	UUID        string `json:"id"`
	Name        string `json:"name"`
	Tel1        string `json:"tel1"`
	Tel2        string `json:"tel2"`
	NotRelevant bool   `json:"notRelevant"`
}

func (r *clientRepository) Create(ctx context.Context, name string) (*models.Client, error) {
	clientCreateRequest := struct {
		Name         string    `json:"name"`
		ActivityDate time.Time `json:"activityDate"`
		Telephone    string    `json:"telephone"`
	}{
		Name:         name,
		ActivityDate: time.Now(),
		Telephone:    "",
	}

	response, statusCode, err := connector.MakeRequest[client](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/client/%s", r.connector.ServerURL, r.connector.ServerToken),
		clientCreateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create client"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to create client. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return &models.Client{
		Namespace:      r.backend.Namespace,
		UUID:           response.UUID,
		Name:           response.Name,
		ContactPersons: []models.ContactPerson{},
		LastUpdateTime: time.Now(),
		CreationTime:   time.Now(),
		Version:        -1,
	}, nil
}
func (r *clientRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Client, error) {
	uuid = url.QueryEscape(uuid)

	response, statusCode, err := connector.MakeRequest[[]client](
		ctx,
		r.connector,
		"GET",
		fmt.Sprintf("%s/client/%s?clientId=%s", r.connector.ServerURL, r.connector.ServerToken, uuid),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get client"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, models.ErrClientNotFound
		}

		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to get client. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	if len(*response) == 0 {
		return nil, models.ErrClientNotFound
	}

	contactPersons, err := r.GetContactPersonsForClient(ctx, uuid, useCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get client contact persons"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	return &models.Client{
		Namespace:      r.connector.Namespace,
		UUID:           (*response)[0].UUID,
		Name:           (*response)[0].Name,
		ContactPersons: contactPersons,
		LastUpdateTime: time.Now(),
		CreationTime:   time.Now(),
		Version:        -1,
	}, nil
}
func (r *clientRepository) GetAll(ctx context.Context, useCache bool) ([]models.Client, error) {
	response, statusCode, err := connector.MakeRequest[[]client](
		ctx,
		r.connector,
		"GET",
		fmt.Sprintf("%s/client/%s?clientId=", r.connector.ServerURL, r.connector.ServerToken),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get clients"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to get clients. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	clients := make([]models.Client, len(*response))
	for i, client := range *response {

		// TODO: move this loop to goroutines

		contactPersons, err := r.GetContactPersonsForClient(ctx, client.UUID, useCache)
		if err != nil {
			err := errors.Join(errors.New("failed to get client contact persons"), err)
			if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
				r.logger.Error(err.Error())
			}
			return nil, err
		}

		clients[i] = models.Client{
			Namespace:      r.backend.Namespace,
			UUID:           client.UUID,
			Name:           client.Name,
			ContactPersons: contactPersons,
			LastUpdateTime: time.Now(),
			CreationTime:   time.Now(),
			Version:        -1,
		}
	}

	return clients, nil
}
func (r *clientRepository) Update(ctx context.Context, uuid string, name string) (*models.Client, error) {
	clientUpdateRequest := struct {
		UUID         string    `json:"id"`
		Name         string    `json:"name"`
		ActivityDate time.Time `json:"activityDate"`
		Telephone    string    `json:"telephone"`
	}{
		UUID:         uuid,
		Name:         name,
		ActivityDate: time.Now(),
		Telephone:    "",
	}

	_, statusCode, err := MakeRequest[interface{}](
		ctx,
		r.backend,
		"POST",
		fmt.Sprintf("%s/client/%s?clientId=", r.backend.ServerURL, r.backend.ServerToken),
		clientUpdateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update client"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to update client. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return &models.Client{
		Namespace:      r.backend.Namespace,
		UUID:           uuid,
		Name:           name,
		ContactPersons: []models.ContactPerson{},
		LastUpdateTime: time.Now(),
		CreationTime:   time.Now(),
		Version:        -1,
	}, nil
}
func (r *clientRepository) Delete(ctx context.Context, uuid string) (*models.Client, error) {
	client, err := r.Get(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	_, statusCode, err := MakeRequest[interface{}](
		ctx,
		r.backend,
		"DELETE",
		fmt.Sprintf("%s/client/%s?clientId=", r.backend.ServerURL, r.backend.ServerToken),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete client"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to delete client. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return client, nil
}

func (r *clientRepository) AddContactPerson(ctx context.Context, clientUUID string, name string, email string, phone []string, comment string) (*models.ContactPerson, error) {
	tel1 := ""
	tel2 := ""

	if len(phone) >= 1 {
		tel1 = phone[0]
	}
	if len(phone) >= 2 {
		tel2 = phone[1]
	}

	contactAddRequest := struct {
		ClientUUID  string `json:"clientId"`
		Name        string `json:"name"`
		Tel1        string `json:"tel1"`
		Tel2        string `json:"tel2"`
		NotRelevant bool   `json:"notRelevant"`
	}{
		ClientUUID:  clientUUID,
		Name:        name,
		Tel1:        tel1,
		Tel2:        tel2,
		NotRelevant: false,
	}

	createdContact, statusCode, err := MakeRequest[clientContact](
		ctx,
		r.backend,
		"POST",
		fmt.Sprintf("%s/contact/%s", r.backend.ServerURL, r.backend.ServerToken),
		contactAddRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to add client contact"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to add client contact. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return &models.ContactPerson{
		Namespace:  r.backend.Namespace,
		UUID:       createdContact.UUID,
		ClientUUID: clientUUID,
		Name:       createdContact.Name,
		Email:      email,
	}, nil
}
func (r *clientRepository) UpdateContactPerson(ctx context.Context, clientUUID string, contactPersonUUID string, name string, email string, phone []string, notRelevant bool, comment string) (*models.ContactPerson, error) {
	tel1 := ""
	tel2 := ""

	if len(phone) >= 1 {
		tel1 = phone[0]
	}
	if len(phone) >= 2 {
		tel2 = phone[1]
	}

	contactUpdateRequest := struct {
		UUID        string `json:"id"`
		ClientUUID  string `json:"clientId"`
		Name        string `json:"name"`
		Tel1        string `json:"tel1"`
		Tel2        string `json:"tel2"`
		NotRelevant bool   `json:"notRelevant"`
	}{
		UUID:        contactPersonUUID,
		ClientUUID:  clientUUID,
		Name:        name,
		Tel1:        tel1,
		Tel2:        tel2,
		NotRelevant: notRelevant,
	}

	_, statusCode, err := MakeRequest[interface{}](
		ctx,
		r.backend,
		"POST",
		fmt.Sprintf("%s/contact/%s", r.backend.ServerURL, r.backend.ServerToken),
		contactUpdateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to add client contact"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to add client contact. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return &models.ContactPerson{
		Namespace:   r.backend.Namespace,
		UUID:        contactPersonUUID,
		ClientUUID:  clientUUID,
		Name:        name,
		Email:       email,
		Phone:       phone,
		NotRelevant: notRelevant,
		Comment:     comment,
	}, nil
}
func (r *clientRepository) DeleteContactPerson(ctx context.Context, contactPersonUUID string) (*models.ContactPerson, error) {
	contactPersonUUID = url.QueryEscape(contactPersonUUID)

	contactPersons, err := r.GetContactPersonsForClient(ctx, contactPersonUUID, false)
	if err != nil {
		return nil, err
	}
	var contactPerson *models.ContactPerson
	for _, cp := range contactPersons {
		if cp.UUID == contactPersonUUID {
			contactPerson = &cp
			break
		}
	}
	if contactPerson == nil {
		return nil, models.ErrClientContactPersonNotFound
	}

	_, statusCode, err := MakeRequest[struct{}](
		ctx,
		r.backend,
		"DELETE",
		fmt.Sprintf("%s/contact/%s&id=%s", r.backend.ServerURL, r.backend.ServerToken, contactPersonUUID),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete client contact"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to delete client contact. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return contactPerson, nil
}
func (r *clientRepository) GetContactPersonsForClient(ctx context.Context, clientUUID string, useCache bool) ([]models.ContactPerson, error) {
	clientUUID = url.QueryEscape(clientUUID)

	response, statusCode, err := MakeRequest[[]clientContact](
		ctx,
		r.backend,
		"GET",
		fmt.Sprintf("%s/contact/%s&clientId=%s", r.backend.ServerURL, r.backend.ServerToken, clientUUID),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get client contacts"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to get client contacts. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	contactPersons := make([]models.ContactPerson, len(*response))
	for i, contact := range *response {
		phones := []string{}
		if contact.Tel1 != "" {
			phones = append(phones, contact.Tel1)
		}
		if contact.Tel2 != "" {
			phones = append(phones, contact.Tel2)
		}

		contactPersons[i] = models.ContactPerson{
			UUID:        contact.UUID,
			Comment:     "",
			Name:        contact.Name,
			Email:       "",
			Phone:       phones,
			NotRelevant: contact.NotRelevant,
		}
	}

	return contactPersons, nil
}
