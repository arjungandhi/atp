package cli

import (
	"fmt"
	"strings"
	"time"

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
		recurCmd,
		remindCmd,
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

var recurCmd = &bonzai.Cmd{
	Name:    "recur",
	Aliases: []string{"r"},
	Summary: "manage recurring tasks",
	Commands: []*bonzai.Cmd{
		help.Cmd,
		recurEditCmd,
	},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		path, err := TodoDir()
		if err != nil {
			return err
		}

		today := time.Now()
		err = todo.AddRecurringTodosToDir(path, today)
		if err != nil {
			return fmt.Errorf("failed to generate recurring todos: %w", err)
		}

		fmt.Println("Generated recurring todos for today")
		return nil
	},
}

var recurEditCmd = &bonzai.Cmd{
	Name:     "edit",
	Aliases:  []string{"e"},
	Summary:  "edit recurring tasks file",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		path, err := TodoDir()
		if err != nil {
			return err
		}

		recurPath := todo.RecurringTasksPath(path)
		shell.OpenInEditor(recurPath)
		return nil
	},
}

var remindCmd = &bonzai.Cmd{
	Name:    "remind",
	Aliases: []string{"rem"},
	Summary: "manage reminder tasks",
	Commands: []*bonzai.Cmd{
		help.Cmd,
		remindAddCmd,
		remindEditCmd,
		remindListCmd,
	},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		path, err := TodoDir()
		if err != nil {
			return err
		}

		var processDate time.Time
		if len(args) > 0 {
			// Parse specified date
			processDate, err = time.Parse("2006-01-02", args[0])
			if err != nil {
				return fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", args[0])
			}
		} else {
			// Use today's date
			processDate = time.Now()
		}

		err = todo.ProcessReminders(path, processDate)
		if err != nil {
			return fmt.Errorf("failed to process reminders: %w", err)
		}

		fmt.Printf("Processed reminders for %s\n", processDate.Format("2006-01-02"))
		return nil
	},
}

var remindAddCmd = &bonzai.Cmd{
	Name:     "add",
	Aliases:  []string{"a"},
	Summary:  "add a reminder task",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// convert args to a string split by " "
		task_str := strings.Join(args, " ")
		var err error
		
		// if task_str is empty, prompt for input
		if task_str == "" {
			task_str, err = prompt.PromptString("Enter reminder task")
			if err != nil {
				return err
			}
		}

		// Prompt for reminder date
		dateStr, err := prompt.PromptString("Enter reminder date (YYYY-MM-DD)")
		if err != nil {
			return err
		}

		// Validate date format
		_, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", dateStr)
		}

		// Add remind label to task string if not already present
		if !strings.Contains(task_str, "remind:") {
			task_str += " remind:" + dateStr
		}

		// convert this string to a todo task
		reminder := todo.FromString(task_str)

		// Validate that remind label exists
		if _, exists := reminder.Labels["remind"]; !exists {
			return fmt.Errorf("reminder task must have a remind:YYYY-MM-DD label")
		}

		// Add the reminder task
		path, err := TodoDir()
		if err != nil {
			return err
		}

		err = todo.AddReminderTask(path, reminder)
		if err != nil {
			return err
		}

		// print confirmation message
		fmt.Printf("Added reminder task: %s\n", reminder.String())
		return nil
	},
}

var remindEditCmd = &bonzai.Cmd{
	Name:     "edit",
	Aliases:  []string{"e"},
	Summary:  "edit reminders file",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		path, err := TodoDir()
		if err != nil {
			return err
		}

		reminderPath := todo.ReminderTasksPath(path)
		shell.OpenInEditor(reminderPath)
		return nil
	},
}

var remindListCmd = &bonzai.Cmd{
	Name:     "list",
	Aliases:  []string{"l", "ls"},
	Summary:  "list all pending reminder tasks",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		path, err := TodoDir()
		if err != nil {
			return err
		}

		reminderPath := todo.ReminderTasksPath(path)
		reminders, err := todo.LoadReminderTasks(reminderPath)
		if err != nil {
			return err
		}

		if len(reminders) == 0 {
			fmt.Println("No pending reminder tasks")
			return nil
		}

		// Sort by reminder date
		todo.SortRemindersByDate(reminders)

		fmt.Printf("Pending reminder tasks (%d):\n", len(reminders))
		for _, reminder := range reminders {
			remindDate := reminder.Labels["remind"]
			fmt.Printf("  %s: %s\n", remindDate, reminder.String())
		}

		return nil
	},
}

