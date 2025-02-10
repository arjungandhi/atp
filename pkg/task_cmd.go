package atp

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var TaskCmd = &Z.Cmd{
	Name:    "task",
	Summary: "manage tasks for all cmd",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
}

var taskListCmd = &Z.Cmd{
	Name:    "list",
	Summary: "list all tasks in priority order",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(cmd *Z.Cmd, args ...string) error {
		return nil
	},
}
