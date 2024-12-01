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
	},
}

// list all tasks to do
var listCmd = &Z.Cmd{
	Name:    "list",
	Summary: "list all tasks in completion order",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		// load atp dir
		// get all  tasks
		// load the active task from the active task storage place
		// one by one print the tasks in order
		// if this task is the active taks print it in a new color
		return nil
	},
}

var addCmd = &Z.Cmd{
	Name:    "add",
	Summary: "add a new task",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		return nil
	},
}

// start working on the current task
var startCmd = &Z.Cmd{
	Name:    "start",
	Summary: "start working on the current task",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		// load atp dir
		// check the active tasks storage location if theres already a file there raise an error
		// optionally take a task name to start
		// if no arg is passed we start the next task
		// if a speciric task is selected use FZF to find it
		// crate a new file in the active task storage place
		return nil
	},
}

// stop working on the current task
var stopCmd = &Z.Cmd{
	Name:    "stop",
	Summary: "stop working on the current task",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		// load atp
		// load/check if theres a task in the active task storage dir
		// do a defer to delete the existing task file
		// find the task in the list of all atp tasks
		// if we don't have task raise an error cleanly
		// add the duration time from file creation to now to the task duration field
		// dump atp
		return nil
	},
}

// complete
var completeCmd = &Z.Cmd{
	Name:    "complete",
	Summary: "complete a task",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		// load atp dir
		// check the active tasks storage location if theres already a file there raise an error
		// optionally take a task name to start
		// if no arg is passed we complete the next task
		// if a speciric task is selected use FZF to find it
		// mark the taskk as complete
		// dump atp
		return nil
	},
}
