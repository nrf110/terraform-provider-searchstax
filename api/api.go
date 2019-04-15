package api

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/http"
	"github.com/nrf110/terraform-provider-searchstax/api/apikeys"
	"github.com/nrf110/terraform-provider-searchstax/api/deployments"
	"github.com/nrf110/terraform-provider-searchstax/api/users"
)

type API struct {
	Deployments *deployments.Deployments
	Users       *users.Users
	ApiKeys     *apikeys.ApiKeys
}

func New(accountName, username, password, baseUrl string) (*API, error) {
	client := sling.New().
		Client(http.DefaultClient).
		Base(fmt.Sprintf("%s/account/%s", baseUrl, accountName))

	token, err := GetToken(baseUrl, username, password)

	if err != nil {
		return nil, err
	}

	client.Set("Authorization", fmt.Sprintf("Token %s", token.Token))

	return &API{
		Deployments: deployments.New(client),
		Users:       users.New(client),
		ApiKeys: 	 apikeys.New(client),
	}, nil
}
