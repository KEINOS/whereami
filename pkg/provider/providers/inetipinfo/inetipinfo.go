/*
Package inetipinfo provides an interface to the inet-ip.info.
*/
package inetipinfo

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/pkg/errors"
)

const urlDefault = "https://inet-ip.info/json"

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

// Response is the structure of JSON from the API response of inet-ip.info.
type Response struct {
	Provider        string `json:"provider"`
	IP              string `json:"IP"`
	HostName        string `json:"Hostname,omitempty"`
	CountryCode     string `json:"CountryCode,omitempty"`
	CountryName     string `json:"CountryName,omitempty"`
	Accept          string `json:"Accept,omitempty"`
	AcceptEncoding  string `json:"AcceptEncoding,omitempty"`
	AcceptLanguage  string `json:"AcceptLanguage,omitempty"`
	UserAgent       string `json:"UserAgent,omitempty"`
	Via             string `json:"Via,omitempty"`
	XForwardedFor   string `json:"XForwardedFor,omitempty"`
	XForwardedPort  string `json:"XForwardedPort,omitempty"`
	XForwardedProto string `json:"XForwardedProto,omitempty"`
	RequestURI      string `json:"RequestURI,omitempty"`
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new Client for the inet-ip.info API with default values.
func New() *Client {
	return &Client{
		EndpointURL: urlDefault,
	}
}

// ----------------------------------------------------------------------------
//  Methods for Client
// ----------------------------------------------------------------------------

// GetIP returns the current IP address detected by inet-ip.info.
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

	// Parse response. The inet-ip.info API returns in JSON.
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
