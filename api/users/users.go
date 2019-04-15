package users

import (
	"fmt"
	"github.com/dghubble/sling"
)

type Users struct {
	sling *sling.Sling
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type DeleteUserRequest struct {
	Username string `json:"username"`
}

func New(client *sling.Sling) *Users {
	return &Users{
		sling: client,
	}
}

func (users *Users) Enable(deploymentId string) error {
	_, err := users.sling.New().
		Get(fmt.Sprintf("/deployment/%s/solr/auth/enable", deploymentId)).
		ReceiveSuccess(nil)

	return err
}

func (users *Users) Disable(deploymentId string) error {
	_, err := users.sling.New().
		Get(fmt.Sprintf("/deployment/%s/solr/auth/disable", deploymentId)).
		ReceiveSuccess(nil)

	return err
}

func (users *Users) Create(deploymentId string, user CreateUserRequest) error {
	_, err := users.sling.New().
		Post(fmt.Sprintf("/deployment/%s/solr/auth/add-user/", deploymentId)).
		BodyJSON(user).
		ReceiveSuccess(nil)

	return err
}

func (users *Users) Delete(deploymentId, username string) error {
	_, err := users.sling.New().
		Post(fmt.Sprintf("/deployment/%s/solr/auth/delete-user/", deploymentId)).
		BodyJSON(DeleteUserRequest{
			Username: username,
		}).
		ReceiveSuccess(nil)

	return err
}
