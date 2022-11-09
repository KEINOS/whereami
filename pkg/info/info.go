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

	"github.com/pkg/errors"
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
func Log(logs ...string) (int, error) {
	s := strings.TrimSpace(strings.Join(logs, " "))
	if s == "" {
		return 0, nil // do nothing
	}

	lenData, err := logBuff.Write([]byte(Prefix + s + "\n"))

	return lenData, errors.Wrap(err, "failed to write to log buffer")
}

// NormalizeIPv4 will trim the zero padded IP address. For example, "001.001.001.001"
// will be "1.1.1.1".
//
// **Note** that it is not a validator. If the given ipAddress is invalid it will return as is.
func NormalizeIPv4(ipAddress string) string {
	// golden input
	if net.ParseIP(ipAddress).String() == ipAddress {
		return ipAddress
	}

	// Number of numerical labels such as 1.1.1.1 has 4 labels.
	const numChunks = 4

	// Split quad-dotted IP address
	chunks := strings.Split(ipAddress, ".")
	if len(chunks) != numChunks {
		return ipAddress
	}

	ipv4 := make([]byte, numChunks)

	for index, chunk := range chunks {
		chunkInt, err := strconv.Atoi(chunk)
		if err != nil {
			return ipAddress
		}

		// Check for lower and upper bounds
		if chunkInt > 0 && chunkInt <= math.MaxInt8 {
			ipv4[index] = byte(chunkInt)
		}
	}

	return net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3]).String()
}
