package wso2am

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	Config struct {
		EndpointToken  string
		EndpointCarbon string

		ClientName   string
		ClientID     string
		ClientSecret string
		UserName     string
		Password     string
	}
	Client struct {
		config *Config
		client *http.Client
	}
)

func New(config *Config) (*Client, error) {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client := &Client{config, c}
	if config.ClientID == "" || config.ClientSecret == "" {
		id, secret, err := client.RegisterClient(NewClientInfo(config.ClientName, config.UserName))
		if err != nil {
			return nil, err
		}
		config.ClientID = id
		config.ClientSecret = secret
	}
	return client, nil
}

type ClientInfo struct {
	CallbackURL string `json:"callbackUrl"`
	ClientName  string `json:"clientName"`
	Owner       string `json:"owner"`
	GrantType   string `json:"grantType"`
	SaaSApp     bool   `json:"saasApp"`
}

func NewClientInfo(clientName string, owner string) *ClientInfo {
	return &ClientInfo{
		CallbackURL: "www.google.lk",
		ClientName:  clientName,
		Owner:       owner,
		GrantType:   "password refresh_token",
		SaaSApp:     true,
	}
}

func (c *Client) endpointCarbon(path string) string {
	return c.config.EndpointCarbon + path
}

func (c *Client) endpointToken(path string) string {
	return c.config.EndpointToken + path
}

func (c *Client) get(path string, scope string, v interface{}) error {
	req, _ := http.NewRequest("GET", c.endpointCarbon(path), nil)
	if err := c.auth(scope, req); err != nil {
		return err
	}
	return c.do(req, nil, v)
}

func (c *Client) post(path string, scope string, body requestBody, v interface{}) error {
	req, _ := http.NewRequest("POST", c.endpointCarbon(path), nil)
	if err := c.auth(scope, req); err != nil {
		return err
	}
	return c.do(req, body, v)
}

func (c *Client) put(path string, scope string, body requestBody, v interface{}) error {
	req, _ := http.NewRequest("PUT", c.endpointCarbon(path), nil)
	if err := c.auth(scope, req); err != nil {
		return err
	}
	return c.do(req, body, v)
}

func (c *Client) delete(path string, scope string, v interface{}) error {
	req, _ := http.NewRequest("DELETE", c.endpointCarbon(path), nil)
	if err := c.auth(scope, req); err != nil {
		return err
	}
	return c.do(req, nil, &v)
}

func (c *Client) auth(scope string, req *http.Request) error {
	token, err := c.GenerateAccessToken(scope)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	return nil
}

func (c *Client) do(req *http.Request, body requestBody, v interface{}) error {
	if body != nil {
		body.writeTo(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return c.apiError(req, resp, err)
	}
	var b []byte
	if writer, ok := v.(io.Writer); ok {
		if _, err := io.Copy(writer, resp.Body); err != nil {
			return c.apiError(req, resp, errors.New("failed to read the response body"))
		}
	} else {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.apiError(req, resp, errors.New("failed to read response body"))
		}
		b = data
		if resp.Header.Get("Content-Type") == "application/json" {
			if err := json.Unmarshal(data, &v); err != nil {
				return c.apiError(req, resp, errors.New("failed to unmarshal json response body"))
			}
		}
	}
	if resp.StatusCode/100 != 2 {
		if b != nil {
			return c.apiErrorWithResponseBody(req, resp, b)
		} else {
			return c.apiError(req, resp, errors.New("failed to execute the API"))
		}
	}
	return nil
}
