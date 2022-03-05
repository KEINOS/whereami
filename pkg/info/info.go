/*
Package info is a simple logger that displays detailed information when the "--verbose" flag is set.
*/
package info

import (
	"bytes"
	"math"
	"net"
	"strconv"
	"strings"
)

var logBuff bytes.Buffer

// Prefix is the prefix of each record in the log.
var Prefix = "[LOG]: "

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// Clear clears the current log.
func Clear() {
	logBuff = bytes.Buffer{}
}

// Get returns the current log.
func Get() string {
	return logBuff.String()
}

// Log writes the given logs to the log buffer.
//
// Note that if a "logs" is empty, or all blank, nothing is recorded.
func Log(logs ...string) (n int, err error) {
	s := strings.TrimSpace(strings.Join(logs, " "))
	if s == "" {
		return 0, nil // do nothing
	}

	return logBuff.Write([]byte(Prefix + s + "\n"))
}

// NormalizeIPv4 will trim the zero padded IP address. For example, "001.001.001.001"
// will be "1.1.1.1".
//
// Note that it is not a validator. If the given "ip" is invalid it will return as is.
func NormalizeIPv4(ip string) string {
	// golden input
	if net.ParseIP(ip).String() == ip {
		return ip
	}

	// Split quad-dotted IP address
	chunks := strings.Split(ip, ".")
	if len(chunks) != 4 {
		return ip
	}

	ipv4 := make([]byte, 4)

	for i, chunk := range chunks {
		b, err := strconv.Atoi(chunk)
		if err != nil {
			return ip
		}

		// Check for lower and upper bounds
		if b > 0 && b <= math.MaxInt8 {
			ipv4[i] = byte(b)
		}
	}

	return net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3]).String()
}
