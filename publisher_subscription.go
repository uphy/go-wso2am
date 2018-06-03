package wso2am

import "fmt"

type (
	Subscription struct {
		ID            string `json:"subscriptionId"`
		Tier          string `json:"tier"`
		APIIdentifier string `json:"apiIdentifier"`
		ApplicationID string `json:"applicationId"`
		Status        string `json:"status"`
	}
	SubscriptionResponse struct {
		PageResponse
	}
	SubscriptionBlockState string
)

const (
	SubscriptionBlockStateBlocked         SubscriptionBlockState = "BLOCKED"
	SubscriptionBlockStateProdOnlyBlocked SubscriptionBlockState = "PROD_ONLY_BLOCKED"
)

func (a *SubscriptionResponse) Subscriptions() []Subscription {
	s := []Subscription{}
	for _, elm := range a.List {
		var v Subscription
		convert(elm, &v)
		s = append(s, v)
	}
	return s
}

func (c *Client) Subscriptions(q *PageQuery) (*SubscriptionResponse, error) {
	return c.SubscriptionsByAPI("", q)
}

func (c *Client) SubscriptionsByAPI(id string, q *PageQuery) (*SubscriptionResponse, error) {
	var v SubscriptionResponse
	var params = pageQueryParams(q)
	if id != "" {
		params.Add("apiId", id)
	}
	if err := c.get("api/am/publisher/v0.12/subscriptions?"+params.Encode(), "apim:subscription_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Subscription(id string) (*Subscription, error) {
	var v Subscription
	if err := c.get("api/am/publisher/v0.12/subscriptions/"+id, "apim:subscription_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) BlockSubscription(id string, state SubscriptionBlockState) (*Subscription, error) {
	var v Subscription
	if err := c.post(fmt.Sprintf("api/am/publisher/v0.12/subscriptions/block-subscription?subscriptionId=%s&blockState=%v", id, state), "apim:subscription_block", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) UnblockSubscription(id string) (*Subscription, error) {
	var v Subscription
	if err := c.post("api/am/publisher/v0.12/subscriptions/unblock-subscription?subscriptionId="+id, "apim:subscription_block", &v); err != nil {
		return nil, err
	}
	return &v, nil
}
