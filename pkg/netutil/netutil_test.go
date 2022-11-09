package netutil

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestHTTPGet(t *testing.T) {
	t.Parallel()

	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte("hello")); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	resp, err := HTTPGet(dummySrv.URL)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestHTTPGet_failed_to_create_request(t *testing.T) {
	oldHTTPNewRequestWithContext := httpNewRequestWithContext
	defer func() {
		httpNewRequestWithContext = oldHTTPNewRequestWithContext
	}()

	// Mock http.NewRequestWithContext to return an error
	//nolint:lll // long line is OK for test
	httpNewRequestWithContext = func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("forced error")
	}

	resp, err := HTTPGet("http://localhost/")
	require.Error(t, err, "empty URL should fail")

	defer resp.Body.Close()

	require.Contains(t, err.Error(), "failed to create HTTP request", "it should contain the error reason")
	require.Nil(t, resp, "returned response should be nil on error")
}

func TestHTTPGet_failed_to_do_request(t *testing.T) {
	t.Parallel()

	resp, err := HTTPGet("")
	require.Error(t, err, "empty URL should fail")

	defer resp.Body.Close()

	require.Contains(t, err.Error(), "failed to do HTTP request", "it should contain the error reason")
	require.Nil(t, resp, "returned response should be nil on error")
}
