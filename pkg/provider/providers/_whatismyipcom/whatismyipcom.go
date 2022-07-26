/*
Package whatismyipcom provides an interface to the www.whatismyip.com web service.

Currently it is disabled in the command due to the following issue:
  - https://github.com/KEINOS/whereami/issues/2
*/
package whatismyipcom

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	urlDefault = "https://www.whatismyip.com/"
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
//  Package Functions
// ----------------------------------------------------------------------------

// IsIPv4 returns true if the ip is well formatted.
func IsIPv4(ip string) bool {
	exp := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`
	rex := regexp.MustCompile(exp)

	return rex.Match([]byte(ip))
}

// ScrapeIPv4 returns the first IPv4 address found from the given html.
func ScrapeIPv4(html string) string {
	exp := `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
	rex := regexp.MustCompile(exp)

	ip := rex.FindString(html) // Find IP

	// Trim leading zeroed IP "001.001.001.001" as "1.1.1.1"
	return info.NormalizeIPv4(ip)
}

// ----------------------------------------------------------------------------
//  Type: Client
// ----------------------------------------------------------------------------

// Client holds information to request www.whatismyip.com's URL.
type Client struct {
	EndpointURL string
}

// ----------------------------------------------------------------------------
//  Type: Response
// ----------------------------------------------------------------------------

// Response is the structure of JSON to hold info from www.whatismyip.com.
type Response struct {
	Provider string `json:"provider"`
	IP       string `json:"ip"`
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new Client for the www.whatismyip.com with default values.
func New() *Client {
	return &Client{
		EndpointURL: urlDefault,
	}
}

// GetResponse returns the Response object parsed from the www.whatismyip.com's content body.
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

	doc.Find("#ipv4").Each(func(_ int, s *goquery.Selection) {
		if ip := ScrapeIPv4(s.Text()); IsIPv4(ip) {
			result.IP = ip

			return
		}
	})

	return result, nil
}

// ----------------------------------------------------------------------------
//  Methods for Client
// ----------------------------------------------------------------------------

// GetIP returns the current IP address detected by www.whatismyip.com.
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
