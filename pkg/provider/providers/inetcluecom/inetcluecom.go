/*
Package inetcluecom provides an interface to the inetclue.com API.
*/
package inetcluecom

import (
	"io"
	"net"
	"net/http"
	"regexp"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/pkg/errors"
)

const (
	urlDefault = "http://inetclue.com/"
)

// IOReadAll is a copy of io.ReadAll function to ease mock it's behavior during
// test.
var IOReadAll = io.ReadAll

// LogInfo is a copy of info.Log function to ease mock it's behavior during test.
var LogInfo = info.Log

// ----------------------------------------------------------------------------
//  Type: Client
// ----------------------------------------------------------------------------

// Client holds information to request inetclue.com's URL.
type Client struct {
	EndpointURL string
}

// ----------------------------------------------------------------------------
//  Type: Response
// ----------------------------------------------------------------------------

// Response is the structure of JSON from the API response of inetclue.com.
type Response struct {
	Provider string `json:"provider"`
	IP       string `json:"origin"`
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new Client for the inetclue.com API with default values.
func New() *Client {
	return &Client{
		EndpointURL: urlDefault,
	}
}

// ScrapeIPv4 returns the first IPv4 address found from the given html.
func ScrapeIPv4(html []byte) string {
	exp := `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
	rex := regexp.MustCompile(exp)

	ip := rex.FindString(string(html)) // Find IP

	// Trim leading zeroed IP "001.001.001.001" as "1.1.1.1"
	return info.NormalizeIPv4(ip)
}

// ----------------------------------------------------------------------------
//  Methods for Client
// ----------------------------------------------------------------------------

// GetIP returns the current IP address detected by inetclue.com.
func (c *Client) GetIP() (net.IP, error) {
	// HTTP request
	response, err := http.Get(c.EndpointURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to GET HTTP request")
	}

	defer response.Body.Close()

	// Read response body
	resBody, err := IOReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response body")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"fail to GET response from: %v\nStatus: %v\nResponse body: %v",
			c.EndpointURL,
			response.Status,
			string(resBody),
		)
	}

	ip := ScrapeIPv4(resBody)

	resJSON := new(Response)

	// Add Provider
	resJSON.IP = ip
	resJSON.Provider = c.EndpointURL

	// Log for verbose output
	if _, err := LogInfo("Response info:\n" + resJSON.String()); err != nil {
		return nil, errors.Wrap(err, "failed to log response")
	}

	return net.ParseIP(resJSON.IP), nil
}

// Name returns the URL of the current provider as its name.
func (c *Client) Name() string {
	return c.EndpointURL
}

// URL overrides the default value of the API endpoint URL.
func (c *Client) URL(url string) {
	c.EndpointURL = url
}

// ----------------------------------------------------------------------------
//  Methods for Response
// ----------------------------------------------------------------------------

// String returns the struct pretty in JSON format.
func (r *Response) String() string {
	return util.FmtStructPretty(r)
}
