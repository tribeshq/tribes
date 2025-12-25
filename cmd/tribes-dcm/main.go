package main

import (
	"os"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/cmd/tribes-dcm/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
