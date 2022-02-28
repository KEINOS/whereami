package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func Test_main(t *testing.T) {
	out := capturer.CaptureStdout(func() {
		main()
	})

	assert.Equal(t, out, "Hello, Gopher!\n")
}
