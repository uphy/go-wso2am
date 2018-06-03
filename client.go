package wso2am

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
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
	return c.do(req, &v)
}

func (c *Client) post(path string, scope string, v interface{}) error {
	req, _ := http.NewRequest("POST", c.endpointCarbon(path), nil)
	if err := c.auth(scope, req); err != nil {
		return err
	}
	return c.do(req, &v)
}

func (c *Client) delete(path string, scope string, v interface{}) error {
	req, _ := http.NewRequest("DELETE", c.endpointCarbon(path), nil)
	if err := c.auth(scope, req); err != nil {
		return err
	}
	return c.do(req, &v)
}

func (c *Client) auth(scope string, req *http.Request) error {
	token, err := c.GenerateAccessToken(scope)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	return nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if resp.Header.Get("Content-Type") == "application/json" {
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &v)
		if err != nil {
			return err
		}
	}
	if resp.StatusCode != 200 {
		return c.apiError(req, resp, b)
	}
	return err
}

type errorMessage struct {
	Code        int           `json:"code"`
	Message     string        `json:"message"`
	Description string        `json:"description"`
	MoreInfo    string        `json:"moreInfo"`
	Error       []interface{} `json:"error"`
}

func (e errorMessage) String() string {
	var errStr string
	if len(e.Error) == 0 {
		errStr = ""
	} else {
		errStr = fmt.Sprintf("%#v", e.Error)
	}
	return fmt.Sprintf("%s: %s (moreInfo=%v, error=%v)", e.Message, e.Description, e.MoreInfo, errStr)
}

func (c *Client) apiError(req *http.Request, resp *http.Response, body []byte) error {
	var detail string
	if len(body) == 0 {
		detail = ""
	} else {
		var e errorMessage
		if err := json.Unmarshal(body, &e); err != nil {
			detail = ""
		} else {
			detail = fmt.Sprint(e.String())
		}
	}
	return fmt.Errorf("API error.  (status=%s, detail=%s, url=%v, method=%s)", resp.Status, detail, req.URL, req.Method)
}
