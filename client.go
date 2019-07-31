package xenditgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Client struct
type Client struct {
	APIEnvType               EnvironmentType
	SecretAPIKey             string
	InvoiceDurationInSeconds int

	LogLevel int
	Logger   *log.Logger
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
		LogLevel: 2,
		Logger:   log.New(os.Stderr, "", log.LstdFlags),
	}
}

// ===================== HTTP CLIENT ================================================
var defHTTPTimeout = 80 * time.Second
var httpClient = &http.Client{Timeout: defHTTPTimeout}

// NewRequest : send new request
func (c *Client) NewRequest(method string, fullPath string, body io.Reader) (*http.Request, error) {
	logLevel := c.LogLevel
	logger := c.Logger

	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Request creation failed: ", err)
		}
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.SecretAPIKey, "")

	return req, nil
}

// ExecuteRequest : execute request
func (c *Client) ExecuteRequest(req *http.Request, v interface{}) error {
	logLevel := c.LogLevel
	logger := c.Logger

	if logLevel > 1 {
		logger.Println("Request ", req.Method, ": ", req.URL.Host, req.URL.Path)
	}

	start := time.Now()

	res, err := httpClient.Do(req)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Cannot send request: ", err)
		}
		return err
	}

	// Check http status code
	logger.Println("HTTP Status Code: ", res.StatusCode, http.StatusText(res.StatusCode))
	if res.StatusCode == 200 {
		// TODO
	}
	// TODO untuk yang 4XX
	// TODO untuk yang 5XX

	defer res.Body.Close()

	if logLevel > 2 {
		logger.Println("Completed in ", time.Since(start))
	}

	if err != nil {
		if logLevel > 0 {
			logger.Println("Request failed: ", err)
		}
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Cannot read response body: ", err)
		}
		return err
	}

	if logLevel > 2 {
		logger.Println("Payment response: ", string(resBody))
	}

	if v != nil {
		if err = json.Unmarshal(resBody, v); err != nil {
			return err
		}

		// we're safe to reflect status_code if response not an array
		// if reflect.ValueOf(v).Elem().Kind() != reflect.Slice {
		//	reflect.ValueOf(v).Elem().FieldByName("StatusCode").SetString(strconv.Itoa(res.StatusCode))
		// }
	}

	return nil
}

// Call the Xendit API at specific `path` using the specified HTTP `method`. The result will be
// given to `v` if there is no error. If any error occurred, the return of this function is the error
// itself, otherwise nil.
func (c *Client) Call(method, path string, body io.Reader, v interface{}) error {
	req, err := c.NewRequest(method, path, body)

	if err != nil {
		return err
	}

	return c.ExecuteRequest(req, v)
}

// ===================== END HTTP CLIENT ================================================
