package main

import (
	"os"

	"github.com/mmatur/httppollerhistorybeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
