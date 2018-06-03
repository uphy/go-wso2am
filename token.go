package wso2am

import (
	"net/http"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (c *Client) GenerateAccessToken(scope string) (*AccessToken, error) {
	body := newFormRequestBody()
	body.Add("grant_type", "password")
	body.Add("username", c.config.UserName)
	body.Add("password", c.config.Password)
	body.Add("scope", scope)
	req, _ := http.NewRequest("POST", c.endpointToken("token"), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)

	var v AccessToken
	if err := c.do(req, body, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
