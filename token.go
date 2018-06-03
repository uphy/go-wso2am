package wso2am

import (
	"fmt"
	"net/http"
	"strings"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (c *Client) GenerateAccessToken(scope string) (*AccessToken, error) {
	body := fmt.Sprintf("grant_type=password&username=%s&password=%s&scope=%s", c.config.UserName, c.config.Password, scope)
	req, _ := http.NewRequest("POST", c.endpointToken("token"), strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)

	var v AccessToken
	if err := c.do(req, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
