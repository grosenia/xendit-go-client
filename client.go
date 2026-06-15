package xenditgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// Client struct
type Client struct {
	APIEnvType               EnvironmentType
	SecretAPIKey             string
	ApiVersion               string
	PaymentsAPIVersion       string
	InvoiceDurationInSeconds int

	LogLevel int
}

// NewClient : this function will always be called when the library is in use
func NewClient() Client {
	return Client{
		APIEnvType: Sandbox,

		// LogLevel is the logging level used by the library
		// 0: No logging
		// 1: Errors only
		// 2: Errors + informational (default)
		// 3: Errors + informational + debug
		LogLevel: 3,
	}
}

// ===================== HTTP CLIENT ================================================
var defHTTPTimeout = 80 * time.Second
var httpClient = &http.Client{Timeout: defHTTPTimeout}

// NewRequest : send new request
func (c *Client) NewRequest(method string, fullPath string, body io.Reader) (*http.Request, error) {
	logLevel := c.LogLevel
	log := clog.Get()

	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		if logLevel > 0 {
			log.Error("Request creation failed ", err)
		}
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.SecretAPIKey, "")
	req.Header.Add("api-version", c.ApiVersion)

	return req, nil
}

// NewRequestWithoutAPIVersion is for Customer API and other endpoints that reject legacy api-version headers.
func (c *Client) NewRequestWithoutAPIVersion(method string, fullPath string, body io.Reader) (*http.Request, error) {
	logLevel := c.LogLevel
	log := clog.Get()

	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		if logLevel > 0 {
			log.Error("Request creation failed ", err)
		}
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.SecretAPIKey, "")

	return req, nil
}

// ExecuteRequest : execute request
func (c *Client) ExecuteRequest(req *http.Request, v interface{}) (httpStatus int, err error) {
	logLevel := c.LogLevel
	log := clog.Get()

	start := time.Now()

	res, err := httpClient.Do(req)
	if err != nil {
		if logLevel > 0 {
			log.Error("Cannot send request ", err)
		}
		return httpStatus, err
	}

	defer res.Body.Close()

	_ = start // retained for potential future latency metrics

	if err != nil {
		if logLevel > 0 {
			log.Error("Request failed ", err)
		}
		return httpStatus, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if logLevel > 0 {
			log.Error("Cannot read response body ", err)
		}
		return res.StatusCode, err
	}

	if v != nil {
		// If response body is empty, skip unmarshalling (some endpoints like expire may return empty body)
		if len(resBody) > 0 {
		if err = json.Unmarshal(resBody, v); err != nil {
				if logLevel > 0 {
					log.Warnf("Failed to unmarshal response body: %s, error: %v", string(resBody), err)
				}
				// Don't return error for unmarshal failures, as some endpoints may return non-JSON responses
				// The caller can check httpStatus to determine success
			}
		} else {
			// no body returned; acceptable for some endpoints
		}

		// we're safe to reflect status_code if response not an array
		// if reflect.ValueOf(v).Elem().Kind() != reflect.Slice {
		//	reflect.ValueOf(v).Elem().FieldByName("StatusCode").SetString(strconv.Itoa(res.StatusCode))
		// }
	}

	return res.StatusCode, nil
}

// NewRequest for batch disbursement : send new request
func (c *Client) NewDisbBatchRequest(key string, method string, fullPath string, body io.Reader) (*http.Request, error) {
	logLevel := c.LogLevel
	log := clog.Get()

	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		if logLevel > 0 {
			log.Error("Request creation failed ", err)
		}
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-IDEMPOTENCY-KEY", key)
	req.SetBasicAuth(c.SecretAPIKey, "")
	return req, nil
}

// Call the Xendit API at specific `path` using the specified HTTP `method`. The result will be
// given to `v` if there is no error. If any error occurred, the return of this function is the error
// itself, otherwise nil.
/*
func (c *Client) Call(method, path string, body io.Reader, v interface{}) (httpStatus int, err error) {
	req, err := c.NewRequest(method, path, body)

	if err != nil {
		return httpStatus, err
	}

	return c.ExecuteRequest(req, v)
}
*/
// ===================== END HTTP CLIENT ================================================
