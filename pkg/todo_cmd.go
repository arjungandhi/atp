package atp

import (
	"github.com/arjungandhi/go-utils/pkg/shell"
	bonzai "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var TodoCmd = &bonzai.Cmd{
	Name:    "todo",
	Aliases: []string{"t"},
	Summary: "manage todos",
	Commands: []*bonzai.Cmd{
		help.Cmd,
		taskEditCmd,
	},
}

var taskEditCmd = &bonzai.Cmd{
	Name:     "edit",
	Aliases:  []string{"e"},
	Summary:  "edit the tasks",
	Commands: []*bonzai.Cmd{help.Cmd, taskEditAllCmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get the todo tasks, path
		path, err := getTodoPath()
		if err != nil {
			return err
		}

		// Open the tasks file in the editor
		shell.OpenInEditor(path)

		return nil
	},
}

var taskEditAllCmd = &bonzai.Cmd{
	Name:     "all",
	Aliases:  []string{"a"},
	Summary:  "edit all tasks",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get tasks path and done path
		path, err := getTodoPath()

		if err != nil {
			return err
		}

		done_path, err := getDoneTodoPath()
		if err != nil {
			return err
		}

		// Open the tasks file in the editor
		shell.OpenInEditor(path, done_path)

		return nil
	},
}
