package provider

import (
	"net"

	"github.com/KEINOS/whereami/pkg/provider/providers/inetcluecom"
	"github.com/KEINOS/whereami/pkg/provider/providers/inetipinfo"
	"github.com/KEINOS/whereami/pkg/provider/providers/ipifyorg"
	"github.com/KEINOS/whereami/pkg/provider/providers/ipinfoio"
	"github.com/KEINOS/whereami/pkg/provider/providers/toolpageorg"
)

// Provider is the interface which each provider package must implement.
type Provider interface {
	// GetIP returns the global/public IP address of the current machine.
	GetIP() (net.IP, error)
	// SetURL overrides the default value of the API endpoint URL.
	SetURL(url string)
	// Name returns the current providers URL.
	Name() string
}

// GetAll returns all providers.
//
// Note that if you implement a new provider, you must add it in this function.
func GetAll() []Provider {
	return []Provider{
		ipinfoio.New(),    // IPInfo.io
		inetipinfo.New(),  // Inet-IP.info
		inetcluecom.New(), // InetClue.com
		toolpageorg.New(), // ToolPage.org
		ipifyorg.New(),    // ipify.org
		// WhatIsMyIP.com is Currently disabled. See: https://github.com/KEINOS/whereami/issues/2
		// whatismyipcom.New(),
	}
}
