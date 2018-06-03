package wso2am

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	APIError struct {
		StatusCode int
		Status     string
		URL        string
		Method     string
		Cause      error
	}
	ErrorResponse struct {
		Code        int           `json:"code"`
		Message     string        `json:"message"`
		Description string        `json:"description"`
		MoreInfo    string        `json:"moreInfo"`
		ErrorObject []interface{} `json:"error"`
	}
)

func (e *APIError) Error() string {
	return fmt.Sprintf("API error.  (status=%s, cause=%v, url=%v, method=%s)", e.Status, e.Cause, e.URL, e.Method)
}

func (e *ErrorResponse) Error() string {
	var errStr string
	if len(e.ErrorObject) == 0 {
		errStr = ""
	} else {
		errStr = fmt.Sprintf("%#v", e.ErrorObject)
	}
	return fmt.Sprintf("%s: %s (moreInfo=%v, error=%v)", e.Message, e.Description, e.MoreInfo, errStr)
}

func (c *Client) apiErrorWithResponseBody(req *http.Request, resp *http.Response, body []byte) *APIError {
	var err error
	if body != nil {
		var e ErrorResponse
		if e2 := json.Unmarshal(body, &e); e2 == nil {
			err = &e
		}
	}
	if err == nil {
		err = fmt.Errorf("Server haven't returned the error information we expected. (body=%s)", string(body))
	}
	return c.apiError(req, resp, err)
}

func (c *Client) apiError(req *http.Request, resp *http.Response, cause error) *APIError {
	var statusCode int
	var status string
	if resp != nil {
		statusCode = resp.StatusCode
		status = resp.Status
	}
	return &APIError{
		StatusCode: statusCode,
		Status:     status,
		URL:        req.URL.String(),
		Method:     req.Method,
		Cause:      cause,
	}
}
