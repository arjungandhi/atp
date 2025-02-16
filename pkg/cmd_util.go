package atp

import (
	"errors"
	"github.com/arjungandhi/atp/pkg/repo"
	"github.com/arjungandhi/atp/pkg/todo"
	"os"
	"path/filepath"
	"strings"
)

// general utils
// get the user specified ATP directory
func getAtpDir() (string, error) {
	todo_dir := os.Getenv("ATP_DIR")
	if todo_dir == "" {
		todo_dir = "~/.atp"
	}

	todo_dir, err := filepath.Abs(todo_dir)
	if err != nil {
		return "", err
	}

	// make the dir if it does not exist
	err = os.MkdirAll(todo_dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return todo_dir, nil
}

// ------------------------------- Repo Utils -------------------------------

func getRepoDir() (string, error) {
	repo_dir := os.Getenv("REPOS")
	if repo_dir == "" {
		return "", errors.New("REPOS env var not set or empty")
	}

	return repo_dir, nil
}

func getRepos() ([]*repo.Repo, error) {
	repo_dir, err := getRepoDir()
	if err != nil {
		return nil, err
	}

	projects, err := repo.GetRepos(repo_dir)
	if err != nil {
		return nil, err
	}

	return projects, err
}

// checks if user is within list of projects P
// if they are returns project else returs nil
func userInRepo(projects []*repo.Repo) (*repo.Repo, error) {
	// get the current user dir
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// for each project are we currently in it?
	for _, project := range projects {
		if strings.Contains(cwd, project.Dir) {
			return project, nil
		}
	}

	return nil, nil
}

// -------------------------------- Todo Utils --------------------------------
// get all todos
func getAllTodos() ([]*todo.Todo, error) {
	// get our tododir
	todo_dir, err := getAtpDir()
	if err != nil {
		return nil, err
	}
	// there are two files in our todo_dir we care about
	// todo.txt, done.txt
	all_todos := []*todo.Todo{}
	file_paths := []string{"todo.txt", "done.txt"}

	for _, file_path := range file_paths {
		file_todos, err := todo.LoadTodoFile(filepath.Join(todo_dir, file_path))
		if err != nil {
			return nil, err
		}
		all_todos = append(all_todos, file_todos...)
	}

	return all_todos, nil
}

// write all todos back to the files
func writeAllTodos(todos []*todo.Todo) error {
	// get our tododir
	todo_dir, err := getAtpDir()
	if err != nil {
		return err
	}

	// sort todos into complete and incomplete
	done := []*todo.Todo{}
	not_done := []*todo.Todo{}

	for _, todo := range todos {
		if todo.Done {
			done = append(done, todo)
		} else {
			not_done = append(not_done, todo)
		}
	}

	err = todo.WriteTodoFile(filepath.Join(todo_dir, "done.txt"), done)
	if err != nil {
		return nil
	}

	err = todo.WriteTodoFile(filepath.Join(todo_dir, "not_done.txt"), done)
	if err != nil {
		return nil
	}

	return nil
}
