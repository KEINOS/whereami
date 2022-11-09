package toolpageorg_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KEINOS/whereami/pkg/info"
	"github.com/KEINOS/whereami/pkg/provider/providers/toolpageorg"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

var testDataGolden = `<table class="outputTable">
<tr>
	<td class="outputTableKey">IP Address:</td>
	<td class="outputTableValue">123.123.123.123</td>
</tr>
				<tr>
	<td class="outputTableKey">Host Name:</td>
	<td class="outputTableValue">123x123x123x123.ap123.ftth.mynet.ne.jp</td>
</tr>
				<tr>
	<td class="outputTableKey">IP Version:</td>
	<td class="outputTableValue">IPv4</td>
</tr>
				<tr>
	<td class="outputTableKey">Remote Port:</td>
	<td class="outputTableValue">5963</td>
</tr>

</table>`

// ----------------------------------------------------------------------------
//  Examples
// ----------------------------------------------------------------------------

func ExampleClient_Name() {
	cli := toolpageorg.New()

	fmt.Println(cli.Name())

	// Output: https://en.toolpage.org/tool/ip-address
}

// ----------------------------------------------------------------------------
//  Tests for Methods
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGetIP_golden(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with golden data
		if _, err := w.Write([]byte(testDataGolden)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer restore inetipinfo.LogInfo.
	oldLogInfo := toolpageorg.LogInfo
	defer func() {
		toolpageorg.LogInfo = oldLogInfo
	}()

	// Instantiate toolpageorg
	cli := toolpageorg.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()
	require.NoError(t, err)

	// Assertion
	expect := "123.123.123.123"
	actual := ip.String()
	assert.Equal(t, expect, actual)

	outLog := info.Get() // Get current log

	assert.Contains(t, outLog, "123x123x123x123.ap123.ftth.mynet.ne.jp")
	assert.Contains(t, outLog, "IPv4")
	assert.Contains(t, outLog, "5963")
}

func TestGetIP_fail_get_response(t *testing.T) {
	t.Parallel()

	// Instantiate toolpageorg
	cli := toolpageorg.New()
	cli.SetURL("") // Override URL as empty

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err)
	require.Nil(t, ip, "returned response should be nil on error")

	assert.Contains(t, err.Error(), "failed to get IP address:")
	assert.Contains(t, err.Error(), "failed to GET HTTP request:")
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGetIP_fail_log(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with golden data
		if _, err := w.Write([]byte(testDataGolden)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer restore inetipinfo.LogInfo.
	oldLogInfo := toolpageorg.LogInfo
	defer func() {
		toolpageorg.LogInfo = oldLogInfo
	}()

	toolpageorg.LogInfo = func(logs ...string) (int, error) {
		return 0, errors.New("forced error in LogInfo")
	}

	// Instantiate toolpageorg
	cli := toolpageorg.New()
	cli.SetURL(dummySrv.URL) // Override URL to dummy server

	// Test
	ip, err := cli.GetIP()

	require.Error(t, err)
	require.Nil(t, ip, "returned IP should be nil on error")

	assert.Contains(t, err.Error(), "failed to log response: forced error in LogInfo")
}

// ----------------------------------------------------------------------------
//  Tests for Functions
// ----------------------------------------------------------------------------

func TestGetResponse_fail_parse_url(t *testing.T) {
	t.Parallel()

	malformedURL := string(byte(0x7f))

	res, err := toolpageorg.GetResponse(malformedURL)

	require.Error(t, err, "URL with control code should return an error")
	require.Nil(t, res, "returned response should be nil on error")

	assert.Contains(t, err.Error(), "failed to parse URL")
	assert.Contains(t, err.Error(), "invalid control character in URL")
}

func TestGetResponse_fail_get_response(t *testing.T) {
	t.Parallel()

	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	// Test
	out := capturer.CaptureOutput(func() {
		res, err := toolpageorg.GetResponse(dummySrv.URL)

		require.Error(t, err, "status code other than 200 should return an error")
		require.Nil(t, res, "returned response should be nil on error")

		assert.Contains(t, err.Error(), "fail to GET response from:")
		assert.Contains(t, err.Error(), "400 Bad Request")
		assert.Contains(t, err.Error(), "invalid request")
	})

	assert.Empty(t, out, "output should be empty on error")
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGetResponse_fail_read_response_body(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
	}))
	defer dummySrv.Close()

	// Backup and defer restore io.ReadAll.
	oldIOReadAll := toolpageorg.IOReadAll
	defer func() {
		toolpageorg.IOReadAll = oldIOReadAll
	}()

	// Mock io.ReadAll.
	toolpageorg.IOReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("forced error in io.ReadAll")
	}

	// Test
	out := capturer.CaptureOutput(func() {
		res, err := toolpageorg.GetResponse(dummySrv.URL)

		require.Error(t, err, "status code other than 200 should return an error")
		require.Nil(t, res, "returned response should be nil on error")

		assert.Contains(t, err.Error(), "fail to read response body: forced error in io.ReadAll")
	})

	assert.Empty(t, out, "output should be empty on error")
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGetResponse_fail_create_goquery_document(t *testing.T) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Response with golden data
		if _, err := w.Write([]byte(testDataGolden)); err != nil {
			t.Fatal(err)
		}
	}))
	defer dummySrv.Close()

	// Backup and defer restore NewDocument (goquery.NewDocumentFromReader).
	oldNewDocument := toolpageorg.NewDocument
	defer func() {
		toolpageorg.NewDocument = oldNewDocument
	}()

	// Mock NewDocument (goquery.NewDocumentFromReader).
	toolpageorg.NewDocument = func(r io.Reader) (*goquery.Document, error) {
		return nil, errors.New("forced error in NewDocument")
	}

	// Test
	out := capturer.CaptureOutput(func() {
		res, err := toolpageorg.GetResponse(dummySrv.URL)

		require.Error(t, err, "status code other than 200 should return an error")
		require.Nil(t, res, "returned response should be nil on error")

		assert.Contains(t, err.Error(), "failed to construct goquery document: forced error in NewDocument")
	})

	assert.Empty(t, out, "output should be empty on error")
}
