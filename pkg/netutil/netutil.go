/*
Package netutil provides network-related utility functions, complementing the more common ones in the net package.
*/
package netutil

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// httpNewRequestWithContext is a copy of http.NewRequestWithContext to ease testing.
var httpNewRequestWithContext = http.NewRequestWithContext

func HTTPGet(url string) (*http.Response, error) {
	body := strings.NewReader("")

	request, err := httpNewRequestWithContext(context.Background(), http.MethodGet, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := http.DefaultClient.Do(request)

	return resp, errors.Wrap(err, "failed to do HTTP request")
}
