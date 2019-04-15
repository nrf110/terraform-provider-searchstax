package api

import (
	"fmt"
	"github.com/dghubble/sling"
)

type TokenRequest struct {
	Username string `url:"username"`
	Password string `url:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func GetToken(baseUrl, username, password string) (*Token, error) {
	token := new(Token)

	_, err := sling.New().
		Base(baseUrl).
		Path(fmt.Sprintf("%s/obtain-auth-token", baseUrl)).
		BodyForm(TokenRequest{
			Username: username,
			Password: password,
		}).
		ReceiveSuccess(&token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
