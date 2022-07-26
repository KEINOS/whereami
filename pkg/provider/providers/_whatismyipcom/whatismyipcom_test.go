package whatismyipcom_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KEINOS/whereami/pkg/info"
	"github.com/KEINOS/whereami/pkg/provider/providers/whatismyipcom"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

//nolint:lll // test data
var testDataGolden = `<div class="card-header mb-2">
<p class="h3">My Public IPv4 is: <br><a id="ipv4" href="https://www.whatismyip.com/123.123.123.123/" title="Detailed information about IP address 123.123.123.123">123.123.123.123</a><button title="Copy IPv4" onclick="copyIPv4()" class="ml-1 btn btn-primary btn-sm"><i class="fa fa-files-o" aria-hidden="true"></i></button>
<br>
<div id="ip-version-check"></div>
</p>
</div>
<div class="card-group">
<div class="card">
<div class="card-header">
<p class="h4">My IP Location Info <a href="/ip-address-geolocation-incorrect/" data-toggle="tooltip" target="_blank" data-placement="bottom" title="Is this incorrect?"><i class="fa fa-question-circle" aria-hidden="true"></i></a></p>
</div>
<div class="card-body">
<ul class="list-group text-left">
<li class="list-group-item">City: Tokyo</li>
<li class="list-group-item">State: Tokyo</li>
<li class="list-group-item">Country: Japan</li>
<li class="list-group-item">Postal Code: 123-0123</li>
<li class="list-group-item">Time Zone: +09:00</li>
</ul>
</div>
</div>
<div class="card">
<div class="card-header">
<p class="h4">My IP Hostname</p>
</div>
<div class="card-body">
 <ul class="list-group text-left">
<li class="list-group-item">ISP: DummyNet Corporation</li>
<li class="list-group-item">Host Name: 123x123x123x123.ap123.ftth.dummy.ne.jp</li>
<li class="list-group-item">ASN: <a href="../asn/99999/" title="ASN 17506">99999</a><a href="../asn/" class="ml-1" data-toggle="tooltip" target="_blank" rel="noopener" data-placement="bottom" title="An AS consists of blocks of IP addresses which have a distinctly defined policy for accessing external networks and are administered by a single organization but may be made up of several operators."><i class="fa fa-question-circle" aria-hidden="true"></i></a></li> </ul>
</div>
</div>
</div>
`

// ----------------------------------------------------------------------------
//  Examples
// ----------------------------------------------------------------------------

func ExampleClient_Name() {
	cli := whatismyipcom.New()

	fmt.Println(cli.Name())

	// Output: https://www.whatismyip.com/
}

// ----------------------------------------------------------------------------
//  Tests for Methods
// ----------------------------------------------------------------------------

func TestGetIP_golden(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with golden data
		if _, err := w.Write([]byte(testDataGolden)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer restore whatismyipcom.LogInfo.
	oldLogInfo := whatismyipcom.LogInfo
	defer func() {
		whatismyipcom.LogInfo = oldLogInfo
	}()

	// Instantiate whatismyipcom
	cli := whatismyipcom.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()
	require.NoError(t, err)

	// Assertion
	expect := "123.123.123.123"
	actual := ip.String()
	assert.Equal(t, expect, actual)

	outLog := info.Get() // Get current log

	assert.Contains(t, outLog, "123.123.123.123", "the log should contain the IP address detected")
	assert.Contains(t, outLog, dummySrv.URL, "the log should contain the provider URL")
}

func TestGetIP_fail_get_response(t *testing.T) {
	// Instantiate whatismyipcom
	cli := whatismyipcom.New()
	cli.SetURL("") // Override URL as empty

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err)
	require.Nil(t, ip, "returned response should be nil on error")

	assert.Contains(t, err.Error(), "failed to get IP address:")
	assert.Contains(t, err.Error(), "failed to GET HTTP request:")
}

func TestGetIP_fail_log(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with golden data
		if _, err := w.Write([]byte(testDataGolden)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer restore whatismyipcom.LogInfo.
	oldLogInfo := whatismyipcom.LogInfo
	defer func() {
		whatismyipcom.LogInfo = oldLogInfo
	}()

	whatismyipcom.LogInfo = func(logs ...string) (n int, err error) {
		return 0, errors.New("forced error in LogInfo in " + t.Name())
	}

	// Instantiate whatismyipcom
	cli := whatismyipcom.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err)
	require.Nil(t, ip, "returned IP should be nil on error")

	assert.Contains(t, err.Error(), "failed to log response")
	assert.Contains(t, err.Error(), "forced error in LogInfo in "+t.Name())
}

// ----------------------------------------------------------------------------
//  Tests for Functions
// ----------------------------------------------------------------------------

func TestGetResponse_fail_parse_url(t *testing.T) {
	malformedURL := string(byte(0x7f))

	res, err := whatismyipcom.GetResponse(malformedURL)

	require.Error(t, err, "URL with control code should return an error")
	require.Nil(t, res, "returned response should be nil on error")

	assert.Contains(t, err.Error(), "failed to parse URL")
	assert.Contains(t, err.Error(), "invalid control character in URL")
}

func TestGetResponse_fail_get_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	// Test
	out := capturer.CaptureOutput(func() {
		res, err := whatismyipcom.GetResponse(dummySrv.URL)

		require.Error(t, err, "status code other than 200 should return an error")
		require.Nil(t, res, "returned response should be nil on error")

		assert.Contains(t, err.Error(), "fail to GET response from:")
		assert.Contains(t, err.Error(), "400 Bad Request")
		assert.Contains(t, err.Error(), "invalid request")
	})

	assert.Empty(t, out, "output should be empty on error")
}

func TestGetResponse_fail_read_response_body(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	// Backup and defer restore io.ReadAll.
	oldIOReadAll := whatismyipcom.IOReadAll
	defer func() {
		whatismyipcom.IOReadAll = oldIOReadAll
	}()

	// Mock io.ReadAll.
	whatismyipcom.IOReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("forced error in io.ReadAll")
	}

	// Test
	out := capturer.CaptureOutput(func() {
		res, err := whatismyipcom.GetResponse(dummySrv.URL)

		require.Error(t, err, "status code other than 200 should return an error")
		require.Nil(t, res, "returned response should be nil on error")

		assert.Contains(t, err.Error(), "fail to read response body: forced error in io.ReadAll")
	})

	assert.Empty(t, out, "output should be empty on error")
}

func TestGetResponse_fail_create_goquery_document(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with golden data
		if _, err := w.Write([]byte(testDataGolden)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer restore NewDocument (goquery.NewDocumentFromReader).
	oldNewDocument := whatismyipcom.NewDocument
	defer func() {
		whatismyipcom.NewDocument = oldNewDocument
	}()

	// Mock NewDocument (goquery.NewDocumentFromReader).
	whatismyipcom.NewDocument = func(r io.Reader) (*goquery.Document, error) {
		return nil, errors.New("forced error in NewDocument")
	}

	// Test
	out := capturer.CaptureOutput(func() {
		res, err := whatismyipcom.GetResponse(dummySrv.URL)

		require.Error(t, err, "status code other than 200 should return an error")
		require.Nil(t, res, "returned response should be nil on error")

		assert.Contains(t, err.Error(), "failed to construct goquery document: forced error in NewDocument")
	})

	assert.Empty(t, out, "output should be empty on error")
}

func TestIsIPv4(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"111.111.111.111", true},
		{"255.255.255.255", true},
		{"256.256.256.256", false},
		{"111.111.111.111.111", false},
		{"foo111.111.111.111bar", false},
	} {
		actual := whatismyipcom.IsIPv4(test.input)

		if test.expect {
			require.True(t, actual, "Input IP: %v", test.input)

			continue
		}

		require.False(t, actual, "Input IP: %v", test.input)
	}
}
