package provider_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/KEINOS/whereami/pkg/provider"
)

func ExampleGetAll() {
	listProviders := provider.GetAll()

	// Use the 1st provider
	providerA := listProviders[0]

	// To avoid unnecessary API requests during running the example, the URL is
	// temporarily set to a dummy server. This server returns "123.123.123.123".
	dummyURL, closer := getDummyServerURL()
	defer closer()

	providerA.SetURL(dummyURL) // Override the default API endpoint URL.

	// Get the current global/public IP address
	ipAddress, err := providerA.GetIP()
	if err != nil {
		log.Println(err)

		return
	}

	fmt.Println(ipAddress.String())

	// Output: 123.123.123.123
}

//nolint:nonamedreturns // Allow named returns for readability
func getDummyServerURL() (dummyURL string, deferFn func()) {
	dummySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(`{"ip": "123.123.123.123"}`)); err != nil {
			log.Fatalf("dummy server creation failed during test. Error: %v", err)
		}
	}))

	return dummySrv.URL, dummySrv.Close
}
