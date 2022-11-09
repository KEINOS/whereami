package main

import (
	"log"
	"net"
	"os"
	"testing"

	"github.com/KEINOS/go-utiles/util"
	"github.com/KEINOS/whereami/pkg/info"
	"github.com/KEINOS/whereami/pkg/provider"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

// ----------------------------------------------------------------------------
//  main()
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to mocking global variables
func Test_main_golden(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	dummyIP := "127.0.0.1"

	dFn := func() (net.IP, error) {
		return net.ParseIP(dummyIP), nil
	}

	// Mock listProvider with dummy providers
	// This value will be recovered by restoreFn.
	listProvider = []provider.Provider{
		&DummyStruct{ID: 0, DummyFunc: dFn},
	}

	out := capturer.CaptureStdout(func() {
		main()
	})

	assert.Equal(t, out, dummyIP)
}

//nolint:paralleltest // do not parallelize due to mocking global variables
func Test_main_golden_verbose(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	dummyIP := "127.0.0.1"

	dFn := func() (net.IP, error) {
		return net.ParseIP(dummyIP), nil
	}

	// Mock listProvider and os.Args with dummy values.
	// This values will be recovered by restoreFn.
	listProvider = []provider.Provider{
		&DummyStruct{ID: 0, DummyFunc: dFn},
	}
	os.Args = []string{
		t.Name(),    // dummy app name
		"--verbose", // verbose flag
	}

	out := capturer.CaptureStdout(func() {
		main()
	})

	require.Contains(t, out, "127.0.0.1")
	require.Contains(t, out, "[LOG]:")
	require.Contains(t, out, "Provider http://dummy.com/ returned the global/public IP as: 127.0.0.1")
}

//nolint:paralleltest // do not parallelize due to mocking global variables
func Test_main_no_provider_set(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	var capturedCode int

	// Mock os.Exit. This will be recovered by restoreFn.
	util.OsExit = func(code int) {
		capturedCode = code
	}

	// Mock max provider number to let main() fail.
	// This will be recovered by restoreFn as well.
	maxNumUseDefault = 0

	out := capturer.CaptureStderr(func() {
		main()
	})

	expectCode := 1
	actualCode := capturedCode
	require.Equal(t, expectCode, actualCode, "it should end with status 1 on error")

	assert.Contains(t, out, "zero provider.")
	assert.Contains(t, out, "you need at least one provider")
}

// ----------------------------------------------------------------------------
//  getRandProviders()
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to mocking global variables
func Test_getRandProviders(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	dFn := func() (net.IP, error) {
		return net.ParseIP("127.0.0.0"), nil
	}

	// Mock listProvider with dummy providers
	// This value will be recovered by restoreFn.
	listProvider = []provider.Provider{
		&DummyStruct{ID: 0, DummyFunc: dFn},
		&DummyStruct{ID: 1, DummyFunc: dFn},
		&DummyStruct{ID: 2, DummyFunc: dFn},
		&DummyStruct{ID: 3, DummyFunc: dFn},
		&DummyStruct{ID: 4, DummyFunc: dFn},
	}

	randProviders := getRandProviders()
	resultOK := false

	for index, obj := range randProviders {
		tmpObj, ok := obj.(*DummyStruct)

		require.True(t, ok, "it should be able to convert to *DummyStruct")

		id := tmpObj.ID
		if index != id {
			resultOK = true
		}
	}

	require.True(t, resultOK, "the returned slice should be shuffled")
}

// ----------------------------------------------------------------------------
//  InfoLog()
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to mocking global variables
func TestInfoLog(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	infoLog = func(logs ...string) (int, error) {
		return 0, errors.New("forced error")
	}

	assert.PanicsWithError(t, "forced error", func() {
		InfoLog("foo")
	}, "on panic it should contain the error message")
}

// ----------------------------------------------------------------------------
//  Run()
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to mocking global variables
func TestRun_error_from_GetIP_method(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	// Force Get method be an error
	dFn := func() (net.IP, error) {
		return nil, errors.New("forced error in GetIP method")
	}

	// Mock listProvider with dummy providers
	// This value will be recovered by restoreFn.
	listProvider = []provider.Provider{
		&DummyStruct{ID: 0, DummyFunc: dFn},
	}

	err := Run()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "all returned IP addresses are different from each other")

	// Get Error log
	errLog := info.Get()
	assert.Contains(t, errLog, "forced error in GetIP method")
}

//nolint:paralleltest // do not parallelize due to mocking global variables
func TestRun_error_if_Get_returns_nil(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	// Force Get method return nil with no error
	dFn := func() (net.IP, error) {
		return nil, nil
	}

	// Mock listProvider with dummy providers
	// This value will be recovered by restoreFn.
	listProvider = []provider.Provider{
		&DummyStruct{ID: 0, DummyFunc: dFn},
	}

	err := Run()
	require.Error(t, err)

	logs := info.Get()
	assert.Contains(t, logs, "returned an empty IP address")
}

//nolint:paralleltest // do not parallelize due to mocking global variables
func TestRun_error_if_providers_returns_different_IPs(t *testing.T) {
	restoreFn := backupAndRestore()
	defer restoreFn()

	// Mock listProvider with dummy providers
	// This value will be recovered by restoreFn.
	listProvider = []provider.Provider{
		&DummyStruct{ID: 0, DummyFunc: func() (net.IP, error) {
			return net.ParseIP("127.0.0.1"), nil
		}},
		&DummyStruct{ID: 0, DummyFunc: func() (net.IP, error) {
			return net.ParseIP("169.254.1.1"), nil
		}},
	}

	err := Run()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "all returned IP addresses are different from each other")

	logs := info.Get()

	assert.Contains(t, logs, "127.0.0.1")
	assert.Contains(t, logs, "169.254.1.1")
}

// ============================================================================
//  Helper Functions
// ============================================================================

// Backup and restore global variables.
//
// Do not parallelize tests using this function. Due to temporary mocking global
// variables.
//
//nolint:nonamedreturns // Allow named returns for readablity
func backupAndRestore() (deferFunc func()) {
	oldInfoLog := infoLog
	oldOsArgs := os.Args
	oldMaxNumUseDefault := maxNumUseDefault
	oldListProvider := listProvider
	oldOsExt := util.OsExit
	oldLog := info.Get()

	return func() {
		infoLog = oldInfoLog
		os.Args = oldOsArgs
		listProvider = oldListProvider
		maxNumUseDefault = oldMaxNumUseDefault
		util.OsExit = oldOsExt

		// Clear the current log and restore the old log
		info.Clear()

		if _, err := info.Log(oldLog); err != nil {
			log.Fatalf("failed to restore the backupped log during test: %v", err)
		}
	}
}

// ----------------------------------------------------------------------------
//  Type: DummyStruct
// ----------------------------------------------------------------------------

// DummyStruct is a dummy provider which implements the provider.Provider requirements.
type DummyStruct struct {
	DummyFunc func() (net.IP, error)
	ID        int
}

// GetIP is an implementation of provider.Provider interface.
func (d DummyStruct) GetIP() (net.IP, error) {
	return d.DummyFunc()
}

// Name is an implementation of provider.Provider interface.
func (d DummyStruct) Name() string {
	return "http://dummy.com/"
}

// SetURL is an implementation of provider.Provider interface.
func (d DummyStruct) SetURL(url string) {}
