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
	"github.com/KEINOS/whereami/pkg/netutil"
	"github.com/pkg/errors"
)

const urlDefault = "https://inet-ip.info/json"

// IOReadAll is a copy of io.ReadAll function to ease mock it's behavior during
// test.
var IOReadAll = io.ReadAll

// LogInfo is a copy of info.Log function to ease mock it's behavior during test.
var LogInfo = info.Log

// ============================================================================
//  Type: Client
// ============================================================================

// Client holds information to request ipinfo.io API.
type Client struct {
	EndpointURL string
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
	response, err := netutil.HTTPGet(c.EndpointURL)
	if err != nil || response == nil {
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

	if err := json.Unmarshal(resBody, resJSON); err != nil {
		return nil, errors.Wrap(err, "fail to parse JSON response: \n"+string(resBody))
	}

	// Add Provider
	resJSON.Provider = c.EndpointURL

	// Log for verbose output
	if _, err := LogInfo("Response info:\n" + resJSON.String()); err != nil {
		return nil, errors.Wrap(err, "failed to log response")
	}

	return net.ParseIP(resJSON.IPAddress), nil
}

// Name returns the URL of the current provider as its name.
func (c *Client) Name() string {
	return c.EndpointURL
}

// SetURL overrides the default value of the API endpoint URL.
func (c *Client) SetURL(url string) {
	c.EndpointURL = url
}

// ============================================================================
//  Type: Response
// ============================================================================

// Response is the structure of JSON from the API response of inet-ip.info.
//
//nolint:tagliatelle // Allow UpperCamelCase for JSON keys due to the API response.
type Response struct {
	Provider  string `json:"provider"`
	IPAddress string `json:"ipAddress"`
	License   string `json:"license,omitempty"`
	ASN       struct {
		AutonomousSystemOrganization string `json:"AutonomousSystemOrganization,omitempty"`
		AutonomousSystemNumber       int    `json:"AutonomousSystemNumber,omitempty"`
	} `json:"asn,omitempty"`
	City struct {
		Postal struct {
			Code string `json:"Code,omitempty"`
		} `json:"Postal,omitempty"`
		Continent struct {
			Names struct {
				De   string `json:"de,omitempty"`
				En   string `json:"en,omitempty"`
				Es   string `json:"es,omitempty"`
				Fr   string `json:"fr,omitempty"`
				Ja   string `json:"ja,omitempty"`
				PtBR string `json:"pt-BR,omitempty"`
				Ru   string `json:"ru,omitempty"`
				ZhCN string `json:"zh-CN,omitempty"`
			} `json:"Names,omitempty"`
			Code      string `json:"Code,omitempty"`
			GeoNameID int    `json:"GeoNameID,omitempty"`
		} `json:"Continent,omitempty"`
		City struct {
			Names struct {
				De   string `json:"de,omitempty"`
				En   string `json:"en,omitempty"`
				Es   string `json:"es,omitempty"`
				Fr   string `json:"fr,omitempty"`
				Ja   string `json:"ja,omitempty"`
				PtBR string `json:"pt-BR,omitempty"`
				Ru   string `json:"ru,omitempty"`
				ZhCN string `json:"zh-CN,omitempty"`
			} `json:"Names,omitempty"`
			GeoNameID int `json:"GeoNameID,omitempty"`
		} `json:"City,omitempty"`
		Subdivisions []struct {
			IsoCode string `json:"IsoCode,omitempty"`
			Names   struct {
				De string `json:"de,omitempty"`
				En string `json:"en,omitempty"`
				Es string `json:"es,omitempty"`
				Fr string `json:"fr,omitempty"`
				Ja string `json:"ja,omitempty"`
				Ru string `json:"ru,omitempty"`
			} `json:"Names,omitempty"`
			GeoNameID int `json:"GeoNameID,omitempty"`
		} `json:"Subdivisions,omitempty"`
		Country struct {
			Names struct {
				De   string `json:"de,omitempty"`
				En   string `json:"en,omitempty"`
				Es   string `json:"es,omitempty"`
				Fr   string `json:"fr,omitempty"`
				Ja   string `json:"ja,omitempty"`
				PtBR string `json:"pt-BR,omitempty"`
				Ru   string `json:"ru,omitempty"`
				ZhCN string `json:"zh-CN,omitempty"`
			} `json:"Names,omitempty"`
			IsoCode           string `json:"IsoCode,omitempty"`
			GeoNameID         int    `json:"GeoNameID,omitempty"`
			IsInEuropeanUnion bool   `json:"IsInEuropeanUnion,omitempty"`
		} `json:"Country,omitempty"`
		RegisteredCountry struct {
			Names struct {
				De   string `json:"de,omitempty"`
				En   string `json:"en,omitempty"`
				Es   string `json:"es,omitempty"`
				Fr   string `json:"fr,omitempty"`
				Ja   string `json:"ja,omitempty"`
				PtBR string `json:"pt-BR,omitempty"`
				Ru   string `json:"ru,omitempty"`
				ZhCN string `json:"zh-CN,omitempty"`
			} `json:"Names,omitempty"`
			IsoCode           string `json:"IsoCode,omitempty"`
			GeoNameID         int    `json:"GeoNameID,omitempty"`
			IsInEuropeanUnion bool   `json:"IsInEuropeanUnion,omitempty"`
		} `json:"RegisteredCountry,omitempty"`
		RepresentedCountry struct {
			Names             interface{} `json:"Names,omitempty"`
			Type              string      `json:"Type,omitempty"`
			IsoCode           string      `json:"IsoCode,omitempty"`
			GeoNameID         int         `json:"GeoNameID,omitempty"`
			IsInEuropeanUnion bool        `json:"IsInEuropeanUnion,omitempty"`
		} `json:"RepresentedCountry,omitempty"`
		Location struct {
			TimeZone       string  `json:"TimeZone,omitempty"`
			Latitude       float64 `json:"Latitude,omitempty"`
			Longitude      float64 `json:"Longitude,omitempty"`
			AccuracyRadius int     `json:"AccuracyRadius,omitempty"`
			MetroCode      int     `json:"MetroCode,omitempty"`
		} `json:"Location,omitempty"`
		Traits struct {
			IsAnonymousProxy    bool `json:"IsAnonymousProxy,omitempty"`
			IsSatelliteProvider bool `json:"IsSatelliteProvider,omitempty"`
		} `json:"Traits,omitempty"`
	} `json:"city,omitempty"`
}

// ----------------------------------------------------------------------------
//  Methods for Response
// ----------------------------------------------------------------------------

// String returns the struct pretty in JSON format.
func (r *Response) String() string {
	return util.FmtStructPretty(r)
}
