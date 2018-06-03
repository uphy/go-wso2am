package wso2am

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type (
	APIsResponse struct {
		PageResponse
	}
	API struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		Context      string    `json:"context"`
		Version      string    `json:"version"`
		Provider     string    `json:"provider"`
		Status       APIStatus `json:"status"`
		ThumbnailURI string    `json:"thumbnailUri"`
	}
	// https://github.com/wso2/carbon-apimgt/blob/master/components/apimgt/org.wso2.carbon.apimgt.rest.api.publisher/src/gen/java/org/wso2/carbon/apimgt/rest/api/publisher/dto/APIDTO.java
	APIDetail struct {
		API
		Definition                   string                  `json:"apiDefinition,omitempty"`
		WSDLURI                      *string                 `json:"wsdlUri,omitempty"`
		ResponseCaching              string                  `json:"responseCaching"`
		CacheTimeout                 int                     `json:"cacheTimeout"`
		DestinationStatsEnabled      bool                    `json:"destinationStatsEnabled,omitempty"`
		DefaultVersion               bool                    `json:"isDefaultVersion"`
		Type                         APIType                 `json:"type"`
		Transport                    []APITransport          `json:"transport"`
		Tags                         []string                `json:"tags"`
		Tiers                        []string                `json:"tiers"`
		MaxTPS                       *APIMaxTPS              `json:"maxTps,omitempty"`
		Visibility                   APIVisibility           `json:"visibility"`
		VisibleRoles                 []string                `json:"visibleRoles"`
		EndpointConfig               string                  `json:"endpointConfig"`
		EndpointSecurity             *APIEndpointSecurity    `json:"endpointSecurity"`
		GatewayEnvironments          string                  `json:"gatewayEnvironments"`
		Sequences                    []APISequence           `json:"sequences,omitempty"`
		SubscriptionAvailability     *string                 `json:"subscriptionAvailability,omitempty"`
		SubscriptionAvailableTenants []string                `json:"subscriptionAvailableTenants,omitempty"`
		BusinessInformation          *APIBusinessInformation `json:"businessInformation"`
		CORSConfiguration            *APICORSConfiguration   `json:"corsConfiguration"`
	}
	APIMaxTPS struct {
		Sandbox    int `json:"sandbox"`
		Production int `json:"production"`
	}
	APIEndpointSecurity struct {
		UserName string `json:"username"`
		Type     string `json:"type"`
		Password string `json:"password"`
	}
	// https://github.com/wso2/carbon-apimgt/blob/master/components/apimgt/org.wso2.carbon.apimgt.rest.api.publisher/src/gen/java/org/wso2/carbon/apimgt/rest/api/publisher/dto/SequenceDTO.java
	APISequence struct {
		Name   string  `json:"name"`
		Config *string `json:"config"`
		Type   string  `json:"type"`
	}
	APIBusinessInformation struct {
		BusinessOwnerEmail  string `json:"businessOwnerEmail"`
		TechnicalOwnerEmail string `json:"technicalOwnerEmail"`
		TechnicalOwner      string `json:"technicalOwner"`
		BusinessOwner       string `json:"businessOwner"`
	}
	APICORSConfiguration struct {
		AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins"`
		AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders"`
		AccessControlAllowMethods     []string `json:"accessControlAllowMethods"`
		AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials"`
		CORSConfigurationEnabled      bool     `json:"corsConfigurationEnabled"`
	}
	APIEndpointConfig struct {
		ProductionEndpoints *APIEndpoint `json:"production_endpoints"`
		SandboxEndpoints    *APIEndpoint `json:"sandbox_endpoints"`
		Type                string       `json:"endpoint_type"`
	}
	APIEndpoint struct {
		URL    string      `json:"url"`
		Config interface{} `json:"config"`
	}
	APIVisibility string
	APITransport  string
	APIAction     string
	APIStatus     string
	APIType       string
)

const (
	APITransportHTTP  APITransport = "http"
	APITransportHTTPS APITransport = "https"

	APIActionPublish            APIAction = "Publish"
	APIActionDeployAsPrototype  APIAction = "Deploy as a Prototype"
	APIActionDemoteToCreated    APIAction = "Demote to Created"
	APIActionDemoteToPrototyped APIAction = "Demote to Prototyped"
	APIActinBlock               APIAction = "Block"
	APIActinDeprecate           APIAction = "Deprecate"
	APIActionRePublish          APIAction = "Re-Publish"
	APIActionRetire             APIAction = "Retire"

	APIStatusCreated     APIStatus = "Created"
	APIStatusPublished   APIStatus = "Published"
	APIStatusDeprecated  APIStatus = "Deprecated"
	APIStatusRetired     APIStatus = "Retired"
	APIStatusMaintenance APIStatus = "Maintenance"
	APIStatusPrototyped  APIStatus = "Prototyped"

	APIVisibilityPublic     APIVisibility = "PUBLIC"
	APIVisibilityPrivate    APIVisibility = "PRIVATE"
	APIVisibilityRestricted APIVisibility = "RESTRICTED"

	APITypeHTTP APIType = "HTTP"
	APITypeWS   APIType = "WS"
)

func (c *Client) NewAPI() *APIDetail {
	return &APIDetail{
		API: API{
			ID:           "",
			Name:         "", // required
			Description:  "",
			Context:      "", // required
			Version:      "", // required
			Provider:     c.config.UserName,
			Status:       APIStatusCreated,
			ThumbnailURI: "",
		},
		Definition:              "", // required
		WSDLURI:                 nil,
		ResponseCaching:         "Disabled",
		CacheTimeout:            300,
		DestinationStatsEnabled: false,
		DefaultVersion:          false,
		Type:                    APITypeHTTP,
		Transport:               []APITransport{APITransportHTTP, APITransportHTTPS},
		Tags:                    []string{},
		Tiers:                   []string{"Unlimited"},
		MaxTPS: &APIMaxTPS{
			Sandbox:    5000,
			Production: 1000,
		},
		Visibility:                   APIVisibilityPublic,
		VisibleRoles:                 []string{},
		EndpointConfig:               "", // required
		EndpointSecurity:             nil,
		GatewayEnvironments:          "Production and Sandbox", // required?
		Sequences:                    []APISequence{},
		SubscriptionAvailability:     nil,
		SubscriptionAvailableTenants: []string{},
		BusinessInformation: &APIBusinessInformation{
			BusinessOwner:       "",
			BusinessOwnerEmail:  "",
			TechnicalOwner:      "",
			TechnicalOwnerEmail: "",
		},
		CORSConfiguration: &APICORSConfiguration{
			AccessControlAllowOrigins:     []string{"*"},
			AccessControlAllowHeaders:     []string{"authorization", "Access-Control-Allow-Origin", "Content-Type", "SOAPAction"},
			AccessControlAllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
			AccessControlAllowCredentials: false,
			CORSConfigurationEnabled:      false,
		},
	}
}

func (a *APIDetail) SetEndpointConfig(endpointConfig *APIEndpointConfig) {
	data, _ := json.Marshal(endpointConfig)
	a.EndpointConfig = string(data)
}

func (a *APIDetail) SetDefinitionFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	switch ext {
	case ".json":
		return a.SetDefinitionFromJSON(f)
	case ".yaml", ".yml":
		return a.SetDefinitionFromYAML(f)
	default:
		return fmt.Errorf("unsupported swagger file format: %s", path)
	}
}

func (a *APIDetail) SetDefinitionFromJSON(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	a.Definition = string(data)
	return nil
}

func (a *APIDetail) SetDefinitionFromYAML(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var v map[string]interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return err
	}
	j, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.Definition = string(j)
	return nil
}

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
	return c.post("api/am/publisher/v0.12/apis/change-lifecycle?"+params.Encode(), "apim:api_publish", nil, nil)
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

func (c *Client) CreateAPI(api *APIDetail) (*APIDetail, error) {
	return c.createAPI(api, false)
}

func (c *Client) UpdateAPI(api *APIDetail) (*APIDetail, error) {
	return c.createAPI(api, true)
}

func (c *Client) createAPI(api *APIDetail, update bool) (*APIDetail, error) {
	if api == nil {
		return nil, errors.New("api == nil")
	}
	body, err := json.Marshal(api)
	if err != nil {
		return nil, err
	}

	var v APIDetail
	if update {
		if err := c.put("api/am/publisher/v0.12/apis/"+api.ID, "apim:api_create", bytes.NewReader(body), &v); err != nil {
			return nil, err
		}
	} else {
		if err := c.post("api/am/publisher/v0.12/apis", "apim:api_create", bytes.NewReader(body), &v); err != nil {
			return nil, err
		}
	}
	return &v, nil
}

func (c *Client) APIDefinition(id string) (map[string]interface{}, error) {
	var v map[string]interface{}
	if err := c.get("api/am/publisher/v0.12/apis/"+id+"/swagger", "apim:api_view", &v); err != nil {
		return nil, err
	}
	return v, nil
}
