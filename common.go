package wso2am

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	req, err := http.NewRequest("POST", c.endpointCarbon("client-registration/v0.12/register"), nil)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.config.UserName, c.config.Password)
	if err != nil {
		return "", "", err
	}
	r := struct {
		CallbackURL  string  `json:"callBackURL"`
		JSONString   string  `json:"jsonString"`
		ClientName   *string `json:"clint_name"`
		ClientID     string  `json:"clientId"`
		ClientSecret string  `json:"clientSecret"`
	}{}
	if err := c.do(req, newJSONRequestBody(clientInfo), &r); err != nil {
		return "", "", err
	}
	return r.ClientID, r.ClientSecret, nil
}

func pageQueryParams(p *PageQuery) *url.Values {
	params := url.Values{}
	if p != nil {
		params.Add("limit", fmt.Sprint(p.Limit))
		params.Add("offset", fmt.Sprint(p.Offset))
		params.Add("query", p.Query)
	}
	return &params
}

func convert(from interface{}, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}
