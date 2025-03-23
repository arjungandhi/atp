package atp

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

// rootCmd is the main command for the list command line tool
// its just holds all the other useful commands
var Cmd = &Z.Cmd{
	Name:    "atp",
	Summary: "atp is a command line tool for managing tasks",
	Commands: []*Z.Cmd{
		help.Cmd,
		ProjectCmd,
		TodoCmd,
	},
}
