package provider

import (
	"net"

	"github.com/KEINOS/whereami/pkg/provider/providers/inetcluecom"
	"github.com/KEINOS/whereami/pkg/provider/providers/inetipinfo"
	"github.com/KEINOS/whereami/pkg/provider/providers/ipinfoio"
)

// Provider is the interface which each provider must implement.
type Provider interface {
	// GetIP returns the global/public IP address of the current machine.
	GetIP() (net.IP, error)
	// URL overrides the default value of the API endpoint URL.
	URL(url string)
	// Name returns the current providers URL.
	Name() string
}

// GetAll returns all providers.
//
// Note that if you implement a new provider, you must add it in this function.
func GetAll() []Provider {
	return []Provider{
		ipinfoio.New(),    // ipinfo.io
		inetipinfo.New(),  // inet-ip.info
		inetcluecom.New(), // inetclue.com
	}
}
