package inetipinfo

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_issue9(t *testing.T) {
	// Get golden JSON of server response
	data, err := os.ReadFile("testdata/response.json")
	require.NoError(t, err)

	oldIOReadAll := IOReadAll
	defer func() {
		IOReadAll = oldIOReadAll
	}()

	// Mock server sesponse
	IOReadAll = func(r io.Reader) ([]byte, error) {
		return data, nil
	}

	client := New()

	expect := "123.123.123.123"
	actual, err := client.GetIP()

	require.NoError(t, err, "well-formed JSON should not return error")
	require.Equal(t, expect, actual.String(), "well-formed JSON response should be parsed")
}
