package cli

import (
	"fmt"
	"strings"

	"github.com/arjungandhi/atp/todo"
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
		path, err := TodoDir()
		if err != nil {
			return err
		}

		todo_path := todo.ActiveTodoPath(path)

		// Open the tasks file in the editor
		shell.OpenInEditor(todo_path)

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
		path, err := TodoDir()

		if err != nil {
			return err
		}

		active_path := todo.ActiveTodoPath(path)
		done_path := todo.DoneTodoPath(path)

		// Open the tasks file in the editor
		shell.OpenInEditor(active_path, done_path)

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

		todos, err := GetTodos()
		if err != nil {
			return err
		}

		// add the task to the list
		todos = append(todos, input_todo)

		// write the todos to the file
		WriteTodos(todos)

		// print confirmation message
		fmt.Printf("Added task: %s\n", input_todo.String())
		return nil
	},
}
