package apikeys

import (
	"github.com/dghubble/sling"
)

type ApiKeys struct {
	sling *sling.Sling
}

type CreateApiKeyRequest struct {
	Scope []string `json:"scope"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"apikey"`
}

type ApiKeyAssociationRequest struct {
	ApiKey     string `json:"apikey"`
	Deployment string `json:"deployment"`
}

type ApiKeyAssociationResponse struct {
	Deployments []string
}

type GetApiKeyAssociationsRequest struct {
	Deployment string `json:"deployment"`
}

type GetApiKeyAssociationsRespone struct {
	ApiKeys []string `json:"Apikey"`
}

type RevokeApiKeyRequest struct {
	ApiKey string `json:"apikey"`
}

type RevokeApiKeySuccess struct {
	Success string `json:"success"`
}

type RevokeApiKeyFailure struct {
	Failure string `json:"failure"`
}

type GetDeploymentsRequest struct {
	ApiKey string `json:"apikey"`
}

type GetDeploymentsResponse struct {
	Deployments []string `json:"deployments"`
}

func (failure RevokeApiKeyFailure) Error() string {
	return failure.Failure
}

func New(client *sling.Sling) *ApiKeys {
	return &ApiKeys{
		sling: client.New().
			Path("/apikey"),
	}
}

func (apikeys *ApiKeys) Create(request CreateApiKeyRequest) (*CreateApiKeyResponse, error) {
	response := new(CreateApiKeyResponse)

	_, err := apikeys.sling.New().
		Post("/").
		BodyJSON(request).
		ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (apikeys *ApiKeys) Associate(request ApiKeyAssociationRequest) (*ApiKeyAssociationResponse, error) {
	response := new(ApiKeyAssociationResponse)

	_, err := apikeys.sling.New().
		Post("/associate").
		BodyJSON(request).
		ReceiveSuccess(&response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (apikeys *ApiKeys) Disassociate(request ApiKeyAssociationRequest) (*ApiKeyAssociationResponse, error) {
	response := new(ApiKeyAssociationResponse)

	_, err := apikeys.sling.New().
		Post("/disassociate").
		BodyJSON(request).
		ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (apikeys *ApiKeys) Revoke(request RevokeApiKeyRequest) (*RevokeApiKeySuccess, *RevokeApiKeyFailure) {
	success := new(RevokeApiKeySuccess)
	failure := new(RevokeApiKeyFailure)

	_, err := apikeys.sling.New().
		Post("/revoke").
		BodyJSON(request).
		Receive(success, failure)

	if err != nil {
		return nil, failure
	}

	return success, nil
}

func (apikeys *ApiKeys) GetForDeployment(deploymentId string) ([]string, error) {
	success := new(GetApiKeyAssociationsRespone)

	_, err := apikeys.sling.New().
		Post("/list").
		BodyJSON(GetApiKeyAssociationsRequest{
			Deployment: deploymentId,
		}).
		ReceiveSuccess(success)

	if err != nil {
		return nil, err
	}

	return success.ApiKeys, nil
}

func (apikeys *ApiKeys) Deployments(apikey string) ([]string, error) {
	response := new(GetDeploymentsResponse)

	_, err := apikeys.sling.New().
		Post("/deployments/").
		BodyJSON(GetDeploymentsRequest{
			ApiKey: apikey,
		}).
		ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	return response.Deployments, nil
}
