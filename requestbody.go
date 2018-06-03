package wso2am

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
	requestBody interface {
		writeTo(*http.Request) error
	}
	jsonRequestBody struct {
		value interface{}
	}
	formRequestBody struct {
		url.Values
	}
)

func newFormRequestBody() *formRequestBody {
	return &formRequestBody{url.Values{}}
}

func (j *formRequestBody) writeTo(req *http.Request) error {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Body = ioutil.NopCloser(strings.NewReader(j.Encode()))
	return nil
}

func newJSONRequestBody(v interface{}) *jsonRequestBody {
	return &jsonRequestBody{v}
}

func (j *jsonRequestBody) writeTo(req *http.Request) error {
	if j.value == nil {
		return errors.New("body == nil")
	}
	req.Header.Add("Content-Type", "application/json")
	body, err := json.Marshal(j.value)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	return nil
}
