package inetcluecom_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KEINOS/whereami/pkg/provider/providers/inetcluecom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

//nolint:lll // test data
var responseData = `
<td colspan="1" valign="top"><strong>Details About Your IP Address:</strong> <a href="/my_ip">123.123.123.123</a> <a href="/ip_location?ip=123.123.123.123"><img src="/images/earth_16x16.png" alt="IP Location" width="16" height="16" title="Get geo location of 123.123.123.123"></a> <a href="/ip_icon?ip=123.123.123.123"><img src="/ipicon/123.123.123.123" alt="IP Icon for 123.123.123.123" width="16" height="16" title="IP Icon for 123.123.123.123"></a></td>
<input type="text" name="ip" id="ip" value="123.123.123.123">
`

func TestGetIP_golden(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(responseData)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := inetcluecom.New()
	cli.URL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()
	require.NoError(t, err)

	// Assertion
	expect := "123.123.123.123"
	actual := ip.String()

	assert.Equal(t, expect, actual)
}

func TestGetIP_error_fail_logging(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"origin": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := inetcluecom.New()
	cli.URL(dummySrv.URL) // Override URL to dummy server

	// Backup and defer restore inetclue.com.LogInfo.
	oldLogInfo := inetcluecom.LogInfo
	defer func() {
		inetcluecom.LogInfo = oldLogInfo
	}()

	// Modck LogInfo to force fail logging.
	inetcluecom.LogInfo = func(logs ...string) (n int, err error) {
		return 0, errors.New("forced fail to log")
	}

	// Test
	ip, err := cli.GetIP()
	require.Error(t, err)
	require.Nil(t, ip, "returned IP should be nil on error")

	assert.Contains(t, err.Error(), "failed to log response:")
	assert.Contains(t, err.Error(), "forced fail to log")
}

func TestGetIP_error_no_URL(t *testing.T) {
	cli := inetcluecom.New()
	cli.URL("") // Set empty URL

	// Test
	out := capturer.CaptureOutput(func() {
		ip, err := cli.GetIP()
		require.Error(t, err, "empty URL should return an error")

		assert.Nil(t, ip, "the returned IP should be nil on error")
		assert.Contains(t, err.Error(), "failed to GET HTTP request:")
	})

	assert.Empty(t, out)
}

func TestGetIP_error_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	cli := inetcluecom.New()
	cli.URL(dummySrv.URL) // Override URL to dummy server
	t.Log("Dummy server:", dummySrv.URL)

	// Test
	out := capturer.CaptureOutput(func() {
		ip, err := cli.GetIP()
		require.Error(t, err, "status code other than 200 should return an error")

		assert.Nil(t, ip, "the returned IP should be nil on error")
		assert.Contains(t, err.Error(), "fail to GET response from:")
		assert.Contains(t, err.Error(), "400 Bad Request")
		assert.Contains(t, err.Error(), "invalid request")
	})

	assert.Empty(t, out)
}

func TestGetIP_error_read_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"origin": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer recover
	oldIOReadAll := inetcluecom.IOReadAll
	defer func() {
		inetcluecom.IOReadAll = oldIOReadAll
	}()

	// Force fail read response body
	inetcluecom.IOReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("forced error to read body")
	}

	cli := inetcluecom.New()
	cli.URL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err, "it should return an error on read body failure")
	require.Nil(t, ip, "the IP should be nil on error")
	assert.Contains(t, err.Error(), "fail to read response body")
	assert.Contains(t, err.Error(), "forced error to read body")
}

func TestName(t *testing.T) {
	cli := inetcluecom.New()

	expect := cli.EndpointURL
	actual := cli.Name()

	assert.Equal(t, expect, actual, "currently the provider name should be the endpoint URL")
}

func TestScrapeIPv4(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		// regular
		{input: "123.123.123.123", expect: "123.123.123.123"},
		{input: "001.001.001.001", expect: "1.1.1.1"},
		{input: "001.010.100.101", expect: "1.10.100.101"},
		{input: "234.234.234.234 123.123.123.123", expect: "234.234.234.234"},
		// irregular
		{input: "123.123.123.123/24", expect: "123.123.123.123"},
		{input: "123.123.123.123.254", expect: "123.123.123.123"},
		{input: "123.123.123.1233", expect: "123.123.123.123"},
		{input: "foo123.123.123.123bar", expect: "123.123.123.123"},
	} {
		expect := test.expect
		actual := inetcluecom.ScrapeIPv4([]byte(test.input))

		require.Equal(t, expect, actual, "input: %v", test.input)
	}
}
