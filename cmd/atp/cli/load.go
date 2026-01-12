package cli

import (
	"errors"
	"fmt"
	"github.com/arjungandhi/atp/project"
	"github.com/arjungandhi/atp/repo"
	"github.com/arjungandhi/atp/todo"
	"os"
	"path/filepath"
	"strings"
)

// get the user specified ATP directory
func AtpDir() (string, error) {
	atp_dir := os.Getenv("ATP_DIR")
	if atp_dir == "" {
		atp_dir = "~/.atp"
	}

	// Expand ~ to home directory
	if strings.HasPrefix(atp_dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		atp_dir = filepath.Join(home, atp_dir[2:])
	} else if atp_dir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		atp_dir = home
	}

	atp_dir, err := filepath.Abs(atp_dir)
	if err != nil {
		return "", err
	}

	// make the dir if it does not exist
	err = os.MkdirAll(atp_dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return atp_dir, nil
}

// ------------------------------- Repo Utils -------------------------------

func RepoDir() (string, error) {
	repo_dir := os.Getenv("REPOS")
	if repo_dir == "" {
		return "", errors.New("REPOS env var not set or empty")
	}

	return repo_dir, nil
}

func GetRepos() ([]*repo.Repo, error) {
	repo_dir, err := RepoDir()
	if err != nil {
		return nil, err
	}

	projects, err := repo.GetRepos(repo_dir)
	if err != nil {
		return nil, err
	}

	return projects, err
}

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

// -------------------------------- TODO utils --------------------------------
func TodoDir() (string, error) {
	atp_dir, err := AtpDir()
	if err != nil {
		return "", err
	}
	todo_dir := filepath.Join(atp_dir, "todo")
	if _, err := os.Stat(todo_dir); os.IsNotExist(err) {
		err = os.MkdirAll(todo_dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return todo_dir, nil
}

func GetTodos() ([]*todo.Todo, error) {
	todo_dir, err := TodoDir()
	if err != nil {
		return nil, fmt.Errorf("Unable to load todo file into todos: %w", err)
	}

	todos, err := todo.LoadTodoDir(todo_dir)
	if err != nil {
		return nil, fmt.Errorf("Unable to load todo file into todos: %w", err)
	}

	return todos, nil
}

func WriteTodos(todos []*todo.Todo) error {
	todo_dir, err := TodoDir()
	if err != nil {
		return fmt.Errorf("Unable to load todo file into todos: %w", err)
	}

	err = todo.WriteTodoDir(todo_dir, todos)
	if err != nil {
		return fmt.Errorf("Unable to write todos to file: %w", err)
	}

	return nil
}

// ------------------------------- Project Utils -------------------------------

func ProjectDir() (string, error) {
	atp_dir, err := AtpDir()
	if err != nil {
		return "", err
	}
	project_dir := filepath.Join(atp_dir, "project")
	if _, err := os.Stat(project_dir); os.IsNotExist(err) {
		err = os.MkdirAll(project_dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return project_dir, nil
}

// get all projects
func GetProjects() ([]*project.Project, error) {
	project_dir, err := ProjectDir()

	if err != nil {
		return nil, fmt.Errorf("Unable to load project file into projects: %w", err)
	}

	repos, err := GetRepos()
	if err != nil {
		return nil, fmt.Errorf("Unable to repos: %w", err)
	}

	projects, err := project.LoadProjectsDir(project_dir, repos)
	if err != nil {
		return nil, fmt.Errorf("Unable to load project file into projects: %w", err)
	}

	return projects, nil
}

// Write the projects to the file
func WriteProjects(projects []*project.Project) error {
	project_dir, err := ProjectDir()
	if err != nil {
		return fmt.Errorf("Unable to load project file into projects: %w", err)
	}

	err = project.WriteProjectsDir(project_dir, projects)
	if err != nil {
		return fmt.Errorf("Unable to write projects to file: %w", err)
	}

	return nil
}
