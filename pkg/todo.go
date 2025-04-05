package atp

import (
	"fmt"
	"github.com/arjungandhi/atp/pkg/todo"
	"os"
	"path/filepath"
)

// get path to todo file
func getTodoPath() (string, error) {
	// get ATP dir
	atp_dir, err := getAtpDir()
	if err != nil {
		return "", fmt.Errorf(
			"Could not get atp dir: %w", err,
		)
	}

	// append tasks
	path := filepath.Join(atp_dir, "todo")
	// ensure the file exists, make it if it does not
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to ensure existance of file: %w", err,
		)
	}

	// append todo.txt
	path = filepath.Join(path, "todo.txt")

	// ensure the file exists, make it if it does not
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf(
				"Unable to create file: %w", err,
			)
		}
		defer file.Close()
	}

	return path, nil
}

// load path to done todo file
func getDoneTodoPath() (string, error) {
	// get ATP dir
	atp_dir, err := getAtpDir()
	if err != nil {
		return "", fmt.Errorf(
			"Could not get atp dir: %w", err,
		)
	}

	// append tasks
	path := filepath.Join(atp_dir, "todo")
	// ensure the file exists, make it if it does not
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to ensure existance of file: %w", err,
		)
	}

	// append done.txt
	path = filepath.Join(path, "done.txt")

	// ensure the file exists, make it if it does not
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf(
				"Unable to create file: %w", err,
			)
		}
		defer file.Close()
	}

	return path, nil
}

// get todos
func GetTodos() ([]*todo.Todo, error) {
	// get the todo file path
	path, err := getTodoPath()

	if err != nil {
		return nil, fmt.Errorf(
			"Could not get todo path: %w", err,
		)
	}

	// load the todo file
	todos, err := todo.LoadTodoFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not load todo file: %w", err,
		)
	}

	return todos, nil
}

// get all todos
func GetAllTodos() ([]*todo.Todo, error) {
	// get the todo file path
	path, err := getTodoPath()

	if err != nil {
		return nil, fmt.Errorf(
			"Could not get todo path: %w", err,
		)
	}

	// load the todo file
	todos, err := todo.LoadTodoFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not load todo file: %w", err,
		)
	}

	// get the done todo file path
	done_path, err := getDoneTodoPath()

	if err != nil {
		return nil, fmt.Errorf(
			"Could not get done todo path: %w", err,
		)
	}

	// load the done todo file
	done_todos, err := todo.LoadTodoFile(done_path)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not load done todo file: %w", err,
		)
	}

	return append(todos, done_todos...), nil
}

// write todos
func WriteTodos(todos []*todo.Todo) error {
	path, err := getTodoPath()
	if err != nil {
		return fmt.Errorf("Unable to write todos to file: %w", err)
	}

	return todo.WriteTodoFile(path, todos)
}

// write all todos
func WriteAllTodos(todos []*todo.Todo) error {
	active_path, err := getTodoPath()
	if err != nil {
		return fmt.Errorf("Unable to write todos to file: %w", err)
	}

	done_path, err := getDoneTodoPath()
	if err != nil {
		return fmt.Errorf("Unable to write todos to file: %w", err)
	}

	// separate the todos into done and not done
	done_todos := []*todo.Todo{}
	not_done_todos := []*todo.Todo{}
	for _, t := range todos {
		if t.Done {
			done_todos = append(done_todos, t)
		} else {
			not_done_todos = append(not_done_todos, t)
		}
	}

	// write the todos to the file
	err = todo.WriteTodoFile(active_path, not_done_todos)

	if err != nil {
		return fmt.Errorf("Unable to write todos to file: %w", err)
	}

	// write the done todos to the file
	err = todo.WriteTodoFile(done_path, done_todos)

	if err != nil {
		return fmt.Errorf("Unable to write todos to file: %w", err)
	}

	return nil
}

// Add todo item
func AddTodo(todo *todo.Todo) error {
	// get all todo items
	todos, err := GetAllTodos()
	if err != nil {
		return fmt.Errorf("Unable to get todos: %w", err)
	}

	// append the new todo item
	todos = append(todos, todo)
	// write the todo items to the file
	err = WriteAllTodos(todos)
	if err != nil {
		return fmt.Errorf("Unable to write todos: %w", err)
	}
	return nil
}
