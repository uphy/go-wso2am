package wso2am

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PageResponse struct {
	Count      int           `json:"count"`
	Next       string        `json:"next"`
	Previous   string        `json:"previous"`
	List       []interface{} `json:"list"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

type PageQuery struct {
	Limit  int
	Offset int
	Query  string
}

func (c *Client) RegisterClient(clientInfo *ClientInfo) (clientID string, clientSecret string, err error) {
	bodyBytes, err := json.Marshal(clientInfo)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest("POST", c.endpointCarbon("client-registration/v0.12/register"), bytes.NewReader(bodyBytes))
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.config.UserName, c.config.Password)
	if err != nil {
		return "", "", err
	}
	r := struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}{}
	if err := c.do(req, &r); err != nil {
		return "", "", err
	}
	return r.ClientID, r.ClientSecret, nil
}
