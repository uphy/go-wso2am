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

func (c *Client) SubscriptionsByAPI(id string, subc chan<- Subscription, errc chan<- error, done <-chan struct{}) {
	var entryc = make(chan interface{})
	go func() {
		for {
			for v := range entryc {
				subc <- *c.ConvertToSubscription(v)
			}
		}
	}()
	c.SubscriptionsByAPIRaw(id, entryc, errc, done)
}

func (c *Client) ConvertToSubscription(v interface{}) *Subscription {
	var a Subscription
	convert(v, &a)
	return &a
}

func (c *Client) SubscriptionsByAPIRaw(id string, entryc chan<- interface{}, errc chan<- error, done <-chan struct{}) {
	c.search(entryc, errc, done, func(q *PageQuery) (*PageResponse, error) {
		return c.subscriptionsByAPI(id, q)
	})
}

func (c *Client) subscriptionsByAPI(id string, q *PageQuery) (*PageResponse, error) {
	var v PageResponse
	var params = pageQueryParams(q)
	if id != "" {
		params.Add("apiId", id)
	}
	if err := c.get(c.publisherURL("subscriptions?"+params.Encode()), "apim:subscription_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) Subscription(id string) (*Subscription, error) {
	var v Subscription
	if err := c.get(c.publisherURL("subscriptions/"+id), "apim:subscription_view", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) BlockSubscription(id string, state SubscriptionBlockState) (*Subscription, error) {
	var v Subscription
	if err := c.post(c.publisherURL(fmt.Sprintf("subscriptions/block-subscription?subscriptionId=%s&blockState=%v", id, state)), "apim:subscription_block", nil, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) UnblockSubscription(id string) (*Subscription, error) {
	var v Subscription
	if err := c.post(c.publisherURL("subscriptions/unblock-subscription?subscriptionId="+id), "apim:subscription_block", nil, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
