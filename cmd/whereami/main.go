package main

import (
	"fmt"

	"github.com/KEINOS/go-utiles/util"
)

func main() {
	util.ExitOnErr(Run())
}

// Run is the actual function of the app.
func Run() error {
	fmt.Println("Hello, Gopher!")

	return nil
}
