package wso2am

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type APIsResponse struct {
	PageResponse
}

func (a *APIsResponse) APIs() []API {
	apis := []API{}
	for _, v := range a.List {
		b, _ := json.Marshal(v)
		var api API
		json.Unmarshal(b, &api)
		apis = append(apis, api)
	}
	return apis
}

type API struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Context      string `json:"context"`
	Version      string `json:"version"`
	Provider     string `json:"provider"`
	Status       string `json:"status"`
	ThumbnailURI string `json:"thumbnailUri"`
}

func (c *Client) APIs(q *PageQuery) (*APIsResponse, error) {
	var v APIsResponse
	params := url.Values{}
	if q != nil {
		params.Add("limit", fmt.Sprint(q.Limit))
		params.Add("offset", fmt.Sprint(q.Offset))
		params.Add("query", q.Query)
	}
	if err := c.get("api/am/publisher/v0.12/apis?"+params.Encode(), "apim:api_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}
