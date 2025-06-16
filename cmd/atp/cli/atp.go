package cli

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

// rootCmd is the main command for the list command line tool
// its just holds all the other useful commands
var Cmd = &Z.Cmd{
	Name:    "atp",
	Summary: "automated task planner - manage projects and todos using todo.txt format",
	Commands: []*Z.Cmd{
		help.Cmd,
		ProjectCmd,
		TodoCmd,
	},
}
