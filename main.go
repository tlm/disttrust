package main

import (
	"os"

	"github.com/tlmiller/disttrust/cmd"
)

func main() {
	os.Exit(cmd.Run())
}
