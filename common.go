package wso2am

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	PageResponse struct {
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
	PageQuery struct {
		Limit  int
		Offset int
	}
	SearchFunc func(entryc chan<- interface{}, errc chan<- error, done <-chan struct{})
)

func (c *Client) RegisterClient(clientInfo *ClientInfo) (clientID string, clientSecret string, err error) {
	req, err := http.NewRequest("POST", c.endpointCarbon(fmt.Sprintf("client-registration/%s/register", c.config.APIVersion)), nil)
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
	}
	return &params
}

func (c *Client) SearchResultToSlice(searchFunc SearchFunc) ([]interface{}, error) {
	var (
		entryc = make(chan interface{})
		errc   = make(chan error)
		done   = make(chan struct{})
	)
	result := []interface{}{}
	go func() {
		defer func() {
			close(entryc)
			close(errc)
			close(done)
		}()
		searchFunc(entryc, errc, done)
	}()
l:
	for {
		select {
		case entry, ok := <-entryc:
			if ok {
				result = append(result, entry)
			} else {
				break l
			}
		case err, ok := <-errc:
			if ok {
				done <- struct{}{}
				return nil, err
			}
			break l
		}
	}
	return result, nil
}

func (c *Client) search(entryc chan<- interface{}, errc chan<- error, done <-chan struct{}, searchFunc func(*PageQuery) (*PageResponse, error)) {
	q := &PageQuery{
		Offset: 0,
		Limit:  100,
	}
	for {
		resp, err := searchFunc(q)
		if resp.Count == 0 {
			return
		}
		if err != nil {
			select {
			case errc <- err:
			case <-done:
				return
			}
		}
		for _, a := range resp.List {
			select {
			case entryc <- a:
			case <-done:
				return
			}
		}
		q.Offset += resp.Count
		select {
		case <-done:
			return
		default:
		}
	}
}

func convert(from interface{}, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}

func (c *Client) publisherURL(path string) string {
	return fmt.Sprintf("api/am/publisher/%s/%s", c.config.APIVersion, path)
}
