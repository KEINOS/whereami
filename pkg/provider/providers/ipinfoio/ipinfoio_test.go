package ipinfoio_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KEINOS/whereami/pkg/provider/providers/ipinfoio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

func TestGetIP_golden(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"ip": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := ipinfoio.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

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
		if _, err := w.Write([]byte(`{"ip": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := ipinfoio.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Backup and defer restore inetipinfo.LogInfo.
	oldLogInfo := ipinfoio.LogInfo
	defer func() {
		ipinfoio.LogInfo = oldLogInfo
	}()

	// Modck LogInfo to force fail logging.
	ipinfoio.LogInfo = func(logs ...string) (n int, err error) {
		return 0, errors.New("forced fail to log")
	}

	// Test
	ip, err := cli.GetIP()
	require.Error(t, err)
	require.Nil(t, ip, "returned IP should be nil on error")

	assert.Contains(t, err.Error(), "failed to log response:")
	assert.Contains(t, err.Error(), "forced fail to log")
}

func TestGet_error_no_URL(t *testing.T) {
	cli := ipinfoio.New()
	cli.SetURL("") // Set empty URL

	// Test
	out := capturer.CaptureOutput(func() {
		ip, err := cli.GetIP()
		require.Error(t, err, "empty URL should return an error")

		assert.Nil(t, ip, "the returned IP should be nil on error")
		assert.Contains(t, err.Error(), "failed to GET HTTP request:")
	})

	assert.Empty(t, out, "it should not print anything on error")
}

func TestGet_error_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	cli := ipinfoio.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server
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

	assert.Empty(t, out, "it should not print anything on error")
}

func TestGet_error_read_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"ip": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer recover
	oldIOReadAll := ipinfoio.IOReadAll
	defer func() {
		ipinfoio.IOReadAll = oldIOReadAll
	}()

	// Force fail read response body
	ipinfoio.IOReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("forced error to read body")
	}

	cli := ipinfoio.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err, "it should return an error on read body failure")
	require.Nil(t, ip, "the IP should be nil on error")
	assert.Contains(t, err.Error(), "fail to read response body")
	assert.Contains(t, err.Error(), "forced error to read body")
}

func TestName(t *testing.T) {
	cli := ipinfoio.New()

	expect := cli.EndpointURL
	actual := cli.Name()

	assert.Equal(t, expect, actual, "currently the provider name should be the endpoint URL")
}
