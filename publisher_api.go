package wso2am

import (
	"net/url"
)

type (
	APIsResponse struct {
		PageResponse
	}

	API struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Context      string `json:"context"`
		Version      string `json:"version"`
		Provider     string `json:"provider"`
		Status       string `json:"status"`
		ThumbnailURI string `json:"thumbnailUri"`
	}
	APIDetail struct {
		API
		APIDefinition           string         `json:"apiDefinition"`
		WSDLURI                 string         `json:"wsdlUri,omitempty"`
		ResponseCaching         string         `json:"responseCaching"`
		CacheTimeout            int            `json:"cacheTimeout"`
		DestinationStatsEnabled string         `json:"destinationStatsEnabled,omitempty"`
		DefaultVersion          bool           `json:"isDefaultVersion"`
		Type                    string         `json:"type"`
		APITransport            []APITransport `json:"transport"`
		Tags                    []string       `json:"tags"`
		Tiers                   []string       `json:"tiers"`
		MaxTPS                  *struct {
			Sandbox    int `json:"sandbox"`
			Production int `json:"production"`
		} `json:"maxTps,omitempty"`
		Visibility       string   `json:"visibility"`
		VisibleRoles     []string `json:"visibleRoles"`
		EndpointConfig   string   `json:"endpointConfig"`
		EndpointSecurity *struct {
			UserName string `json:"username"`
			Type     string `json:"type"`
			Password string `json:"password"`
		} `json:"endpointSecurity"`
		GatewayEnvironments          string        `json:"gatewayEnvironments"`
		Sequences                    []interface{} `json:"sequences,omitempty"`
		SubscriptionAvailability     interface{}   `json:"subscriptionAvailability,omitempty"`
		SubscriptionAvailableTenants interface{}   `json:"subscriptionAvailableTenants,omitempty"`
		BusinessInformation          *struct {
			BusinessOwnerEmail  string `json:"businessOwnerEmail"`
			TechnicalOwnerEmail string `json:"technicalOwnerEmail"`
			TechnicalOwner      string `json:"technicalOwner"`
			BusinessOwner       string `json:"businessOwner"`
		} `json:"businessInformation"`
		CORSConfiguration *struct {
			AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins"`
			AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders"`
			AccessControlAllowMethods     []string `json:"accessControlAllowMethods"`
			AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials"`
			CORSConfigurationEnabled      bool     `json:"corsConfigurationEnabled"`
		} `json:"corsConfiguration"`
	}
	APITransport string
	APIAction    string
)

const (
	APITransportHTTP  APITransport = "http"
	APITransportHTTPs APITransport = "https"

	APIActionPublish            APIAction = "Publish"
	APIActionDeployAsPrototype  APIAction = "Deploy as a Prototype"
	APIActionDemoteToCreated    APIAction = "Demote to Created"
	APIActionDemoteToPrototyped APIAction = "Demote to Prototyped"
	APIActinBlock               APIAction = "Block"
	APIActinDeprecate           APIAction = "Deprecate"
	APIActionRePublish          APIAction = "Re-Publish"
	APIActionRetire             APIAction = "Retire"
)

func (a *APIsResponse) APIs() []API {
	s := []API{}
	for _, elm := range a.List {
		var v API
		convert(elm, &v)
		s = append(s, v)
	}
	return s
}

func (c *Client) APIs(q *PageQuery) (*APIsResponse, error) {
	return c.SearchAPIs("", q)
}

func (c *Client) SearchAPIs(query string, q *PageQuery) (*APIsResponse, error) {
	params := pageQueryParams(q)
	if query != "" {
		params.Add("query", query)
	}
	var v APIsResponse
	if err := c.get("api/am/publisher/v0.12/apis?"+params.Encode(), "apim:api_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) ChangeAPIStatus(id string, action APIAction) error {
	params := url.Values{}
	params.Add("apiId", id)
	params.Add("action", string(action))
	return c.post("api/am/publisher/v0.12/apis/change-lifecycle?"+params.Encode(), "apim:api_publish", nil)
}

func (c *Client) DeleteAPI(id string) error {
	return c.delete("api/am/publisher/v0.12/apis/"+id, "apim:api_create", nil)
}

func (c *Client) API(id string) (*APIDetail, error) {
	var v APIDetail
	if err := c.get("api/am/publisher/v0.12/apis/"+id, "apim:api_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}
