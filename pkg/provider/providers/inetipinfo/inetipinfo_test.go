package inetipinfo_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KEINOS/whereami/pkg/provider/providers/inetipinfo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

//nolint:paralleltest // do not parallelize due to the race condition
func TestGetIP_golden(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"ipAddress": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := inetipinfo.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()
	require.NoError(t, err)

	// Assertion
	expect := "123.123.123.123"
	actual := ip.String()

	assert.Equal(t, expect, actual)
}

//nolint:paralleltest // do not parallelize due to race condition
func TestGetIP_error_bad_json(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte("[ipAddress]\n123.123.123.123\n")); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := inetipinfo.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err, "malformed JSON should return an error")
	require.Contains(t, err.Error(), "fail to parse JSON response")
	require.Nil(t, ip, "the returned IP should be nil on error")
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGetIP_error_fail_logging(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"ipAddress": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	cli := inetipinfo.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Backup and defer restore inetipinfo.LogInfo.
	oldLogInfo := inetipinfo.LogInfo
	defer func() {
		inetipinfo.LogInfo = oldLogInfo
	}()

	// Modck LogInfo to force fail logging.
	inetipinfo.LogInfo = func(logs ...string) (int, error) {
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
	t.Parallel()

	cli := inetipinfo.New()
	cli.SetURL("") // Set empty URL

	// Test
	out := capturer.CaptureOutput(func() {
		ip, err := cli.GetIP()
		require.Error(t, err, "empty URL should return an error")

		assert.Nil(t, ip, "the returned IP should be nil on error")
		assert.Contains(t, err.Error(), "failed to GET HTTP request:")
	})

	assert.Empty(t, out)
}

//nolint:paralleltest // do not parallelize due to race condition
func TestGetIP_error_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	cli := inetipinfo.New()
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

	assert.Empty(t, out)
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGetIP_error_read_response(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"ip": "123.123.123.123"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer recover
	oldIOReadAll := inetipinfo.IOReadAll
	defer func() {
		inetipinfo.IOReadAll = oldIOReadAll
	}()

	// Force fail read response body
	inetipinfo.IOReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("forced error to read body")
	}

	cli := inetipinfo.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err, "it should return an error on read body failure")
	require.Nil(t, ip, "the IP should be nil on error")
	assert.Contains(t, err.Error(), "fail to read response body")
	assert.Contains(t, err.Error(), "forced error to read body")
}

func TestName(t *testing.T) {
	t.Parallel()

	cli := inetipinfo.New()

	expect := cli.EndpointURL
	actual := cli.Name()

	assert.Equal(t, expect, actual, "currently the provider name should be the endpoint URL")
}
