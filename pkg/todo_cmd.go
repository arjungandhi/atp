package atp

import (
	"fmt"
	"strings"

	"github.com/arjungandhi/atp/pkg/todo"
	"github.com/arjungandhi/go-utils/pkg/prompt"
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
		taskAddCmd,
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

var taskAddCmd = &bonzai.Cmd{
	Name:     "add",
	Aliases:  []string{"a"},
	Summary:  "add a task",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// convert args to a string split by " "
		task_str := strings.Join(args, " ")
		var err error
		// if task_str is empty, prompt for input
		if task_str == "" {
			// prompt for input
			task_str, err = prompt.PromptString("Enter Task")
			if err != nil {
				return err
			}
		}

		// convert this string to a todo task
		input_todo := todo.FromString(task_str)

		// add the todo to the list
		err = AddTodo(input_todo)
		if err != nil {
			return err
		}

		// print confirmation message
		fmt.Printf("Added task: %s\n", input_todo.String())
		return nil
	},
}
