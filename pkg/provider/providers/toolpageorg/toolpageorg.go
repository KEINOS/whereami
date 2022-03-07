/*
Package toolpageorg provides an interface to the en.toolpage.org web service.
*/
package toolpageorg

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	urlDefault = "https://en.toolpage.org/tool/ip-address"
)

// IOReadAll is a copy of io.ReadAll function to ease mock it's behavior during
// test.
var IOReadAll = io.ReadAll

// LogInfo is a copy of info.Log function to ease mock it's behavior during test.
var LogInfo = info.Log

// NewDocument is a copy of goquery.NewDocumentFromReader function to ease mock
// it's behavior during test.
var NewDocument = goquery.NewDocumentFromReader

// ----------------------------------------------------------------------------
//  Type: Client
// ----------------------------------------------------------------------------

// Client holds information to request en.toolpage.org's URL.
type Client struct {
	EndpointURL string
}

// ----------------------------------------------------------------------------
//  Type: Response
// ----------------------------------------------------------------------------

// Response is the structure of JSON from the API response of en.toolpage.org.
type Response struct {
	Provider   string `json:"provider"`
	IP         string `json:"ip"`
	Hostname   string `json:"hostname,omitempty"`
	IPVersion  string `json:"ip_version,omitempty"`
	RemotePort string `json:"remote_port,omitempty"`
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new Client for the en.toolpage.org API with default values.
func New() *Client {
	return &Client{
		EndpointURL: urlDefault,
	}
}

// GetResponse returns the Response object parsed from the en.toolpage.org's content body.
func GetResponse(urlProvider string) (*Response, error) {
	result := &Response{
		Provider: urlProvider,
	}

	// Validate URL to avoid gosec G107 vulnerability: Potential HTTP request made with variable url.
	parsedURL, err := url.Parse(urlProvider)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	// HTTP request
	response, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to GET HTTP request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// Read error response body
		resBody, resErr := IOReadAll(response.Body)
		if resErr != nil {
			return nil, errors.Wrap(resErr, "fail to read response body")
		}

		return nil, errors.Errorf(
			"fail to GET response from: %v\nStatus: %v\nResponse body: %v",
			urlProvider,
			response.Status,
			string(resBody),
		)
	}

	// Parse document from response body
	doc, err := NewDocument(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct goquery document")
	}

	doc.Find(".outputTableKey").Each(func(_ int, s *goquery.Selection) {
		switch strings.TrimSpace(s.Text()) {
		case "IP Address:":
			result.IP = s.Next().Text()
		case "Host Name:":
			result.Hostname = s.Next().Text()
		case "IP Version:":
			result.IPVersion = s.Next().Text()
		case "Remote Port:":
			result.RemotePort = s.Next().Text()
		}
	})

	return result, nil
}

// ----------------------------------------------------------------------------
//  Methods for Client
// ----------------------------------------------------------------------------

// GetIP returns the current IP address detected by en.toolpage.org.
func (c *Client) GetIP() (net.IP, error) {
	result, err := GetResponse(c.EndpointURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get IP address")
	}

	// Log for verbose output
	if _, err := LogInfo("Response info:\n" + result.String()); err != nil {
		return nil, errors.Wrap(err, "failed to log response")
	}

	return net.ParseIP(result.IP), nil
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
