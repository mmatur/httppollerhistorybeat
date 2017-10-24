package cmd

import (
	"github.com/mmatur/httppollerhistorybeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "httppollerhistorybeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "1.0.0", beater.New)
