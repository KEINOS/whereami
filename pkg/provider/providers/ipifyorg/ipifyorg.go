/*
Package ipifyorg implements an interface to the ipify.org API.

- Limits of the ipify.org API: 1000 request/day or 50,000 requests/month for free plan.
*/
package ipifyorg

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/KEINOS/whereami/pkg/netutil"
	"github.com/pkg/errors"
)

// This endpoint returns in JSON with IPv4/IPv6 address.
const urlDefault = "https://api64.ipify.org?format=json"

// IOReadAll is a copy of io.ReadAll function to ease mock it's behavior during
// test.
var IOReadAll = io.ReadAll

// LogInfo is a copy of info.Log function to ease mock it's behavior during test.
var LogInfo = info.Log

// ----------------------------------------------------------------------------
//  Type: Client
// ----------------------------------------------------------------------------

// Client holds information to request ipify.org API.
type Client struct {
	EndpointURL string
}

// ----------------------------------------------------------------------------
//  Type: Response
// ----------------------------------------------------------------------------

// Response is the structure of JSON from the API response of ipify.org.
type Response struct {
	Provider     string `json:"provider"`
	IP           string `json:"ip"`
	HostName     string `json:"hostname,omitempty"`
	City         string `json:"city,omitempty"`
	Region       string `json:"region,omitempty"`
	Location     string `json:"loc,omitempty"`
	Organization string `json:"org,omitempty"`
	PostalCode   string `json:"postal,omitempty"`
	TimeZone     string `json:"timezone,omitempty"`
	Readme       string `json:"readme,omitempty"`
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new Client for the ipify.org API with default values.
func New() *Client {
	return &Client{
		EndpointURL: urlDefault,
	}
}

// ----------------------------------------------------------------------------
//  Methods for Client
// ----------------------------------------------------------------------------

// GetIP returns the current IP address detected by ipify.org.
func (c *Client) GetIP() (net.IP, error) {
	// HTTP request
	response, err := netutil.HTTPGet(c.EndpointURL)
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

	// Parse response. The ipify.org API returns in JSON.
	resJSON := new(Response)
	_ = json.Unmarshal(resBody, resJSON)

	// Add Provider
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

// SetURL overrides the default value of the API endpoint URL.
func (c *Client) SetURL(url string) {
	c.EndpointURL = url
}

// ----------------------------------------------------------------------------
//  Methods for Response
// ----------------------------------------------------------------------------

// String returns the struct pretty in JSON format.
func (r *Response) String() string {
	return util.FmtStructPretty(r)
}
