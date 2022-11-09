/*
Package ipinfoio provides an interface to the ipinfo.io API.
*/
package ipinfoio

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/pkg/errors"
)

const urlDefault = "https://ipinfo.io/"

// IOReadAll is a copy of io.ReadAll function to ease mock it's behavior during
// test.
var IOReadAll = io.ReadAll

// LogInfo is a copy of info.Log function to ease mock it's behavior during test.
var LogInfo = info.Log

// ----------------------------------------------------------------------------
//  Type: Client
// ----------------------------------------------------------------------------

// Client holds information to request ipinfo.io API.
type Client struct {
	EndpointURL string
}

// ----------------------------------------------------------------------------
//  Type: Response
// ----------------------------------------------------------------------------

// Response is the structure of JSON from the API response of ipinfo.io.
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

// New returns a new Client for the ipinfo.io API with default values.
func New() *Client {
	return &Client{
		EndpointURL: urlDefault,
	}
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

func httpGet(url string) (*http.Response, error) {
	body := strings.NewReader("")

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := http.DefaultClient.Do(request)

	return resp, errors.Wrap(err, "failed to do HTTP request")
}

// ----------------------------------------------------------------------------
//  Methods for Client
// ----------------------------------------------------------------------------

// GetIP returns the current IP address detected by ipinfo.io.
func (c *Client) GetIP() (net.IP, error) {
	// HTTP request
	response, err := httpGet(c.EndpointURL)
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

	// Parse response. The ipinfo.io API returns in JSON.
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
