package info

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleLog() {
	// To avoid conflicts during testing this example, back up the current buffer
	// and restore it later. This is usually not necessary.
	oldLogBuff := logBuff
	defer func() {
		logBuff = oldLogBuff
	}()

	/* Example usage */

	// Log
	if _, err := Log("This", "is", "the", "log1"); err != nil {
		log.Println(err)

		return
	}

	// Logging
	if _, err := Log("This", "is", "the", "log2"); err != nil {
		log.Println(err)

		return
	}

	// Get the current log.
	fmt.Println("Log output:")
	fmt.Println(Get())

	// Clear the log.
	Clear()

	// After calling `Clear()` the log becomes empty.
	fmt.Printf("Is log after 'Clear()' call empty: ")

	if Get() == "" {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}

	// Output:
	// Log output:
	// [LOG]: This is the log1
	// [LOG]: This is the log2
	//
	// Is log after 'Clear()' call empty: yes
}

//nolint:paralleltest // do not parallelize due to mocking global function variables
func TestGet_whitespaces(t *testing.T) {
	// Backup and defer restore.
	oldLogBuff := logBuff
	defer func() {
		logBuff = oldLogBuff
	}()

	Clear()

	n, err := Log(" ", " ", " ")

	require.NoError(t, err)
	require.Zero(t, n, "empty or white-spaces log should not be logged")

	result := Get()
	assert.Empty(t, result, "empty or white-spaces log should not be logged")
}

func TestNormalizeIPv4(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		input  string
		expect string
	}{
		// golden
		{input: "123.123.123.123", expect: "123.123.123.123"},
		// normalize
		{input: "001.001.001.001", expect: "1.1.1.1"},
		{input: "001.010.100.101", expect: "1.10.100.101"},
		// invalid
		{input: "123.123.123", expect: "123.123.123"},
		{input: "123.123.123.123.123", expect: "123.123.123.123.123"},
		{input: "foo123.123.123.123bar", expect: "foo123.123.123.123bar"},
	} {
		expect := test.expect
		actual := NormalizeIPv4(test.input)

		require.Equal(t, expect, actual, "input: %v", test.input)
	}
}
