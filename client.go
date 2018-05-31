package wso2am

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Config struct {
	EndpointToken  string
	EndpointCarbon string

	ClientName   string
	ClientID     string
	ClientSecret string
	UserName     string
	Password     string
}

type Client struct {
	config *Config
	client *http.Client
}

func New(config *Config) (*Client, error) {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client := &Client{config, c}
	if config.ClientID == "" || config.ClientSecret == "" {
		id, secret, err := client.RegisterClient(NewClientInfo(config.ClientName))
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

func NewClientInfo(clientName string) *ClientInfo {
	return &ClientInfo{
		CallbackURL: "www.google.lk",
		ClientName:  clientName,
		Owner:       "admin",
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
	token, err := c.GenerateAccessToken(scope)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	return c.do(req, &v)
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}
