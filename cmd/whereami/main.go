package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/KEINOS/whereami/pkg/provider"
)

const sleepTime = 1

var (
	// Max number of providers to use to fetch the public IP.
	maxNumUseDefault = 3
	// List of public IP address detector service providers.
	listProvider []provider.Provider
)

var (
	// This infoLog is a copy of info.Log to ease mock its behavior during test.
	infoLog = info.Log

	// InfoLog is a function that wraps info.Log, but panics immediately if it
	// fails to log.
	InfoLog = func(logs ...string) {
		if _, err := infoLog(logs...); err != nil {
			panic(err)
		}
	}
)

/* Flag variables */

// Variable of --verbose option flag.
var isVerbose bool

// ----------------------------------------------------------------------------
//  Main
// ----------------------------------------------------------------------------

func init() {
	// Set global/public IP address detection service providers
	listProvider = provider.GetAll()
	// Define flag options
	flag.BoolVar(&isVerbose, "verbose", false, "prints detailed information if any. such as IPv6 and etc.")
}

func main() {
	flag.Parse()          // Parse the flag options
	info.Clear()          // Ensure to clear the log before run
	util.ExitOnErr(Run()) // Print the current global/public IP address
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// Returns providers in random order.
func getRandProviders() []provider.Provider {
	l := listProvider

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })

	return l
}

// Calls GetIP method from the given provider and returns the detected IP address.
func request(prov provider.Provider) (net.IP, error) {
	ip, err := prov.GetIP()

	switch {
	case err != nil:
		errMsg := fmt.Sprintf("provider %v returned an error: %v", prov, err.Error())

		return nil, errors.New(errMsg)
	case ip == nil:
		errMsg := fmt.Sprintf("provider %v returned an empty IP address", prov.Name())

		return nil, errors.New(errMsg)
	default:
		return ip, nil
	}
}

// Returns the IPv4 address if all maxNumUse providers returns the same IP.
func getIPPublic(maxNumUse int) (string, error) {
	if maxNumUse == 0 {
		return "", errors.New("error: zero provider. you need at least one provider")
	}

	providers := getRandProviders()
	if lenProv := len(providers); maxNumUse > lenProv {
		maxNumUse = lenProv
	}

	foundIP := make(map[string]int)

	for _, prov := range providers {
		ip, err := request(prov)
		if err != nil {
			InfoLog(fmt.Sprintf("%v: %v", prov.Name(), err.Error()))

			continue
		}

		key := ip.String()

		InfoLog(fmt.Sprintf(
			"Provider %v returned the global/public IP as: %v",
			prov.Name(),
			key,
		))

		foundIP[key]++

		if foundIP[key] == maxNumUse {
			return key, nil // IP Found!
		}
	}

	return "", errors.New("all returned IP addresses are different from each other")
}

// Run is the actual function of the app.
func Run() error {
	ip, err := getIPPublic(maxNumUseDefault)
	if err != nil {
		return err
	}

	fmt.Printf("%v", ip)

	// Print verbose information
	if isVerbose {
		fmt.Printf("\n%v", info.Get())
	}

	// Force sleep to avoide large number of requests.
	time.Sleep(sleepTime * time.Second)

	return nil
}
