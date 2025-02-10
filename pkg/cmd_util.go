package atp

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ------------------------------- Project Utils -------------------------------

func getRepoDir() (string, error) {
	repo_dir := os.Getenv("REPOS")
	if repo_dir == "" {
		return "", errors.New("REPOS env var not set or empty")
	}

	return repo_dir, nil
}

func getRepoProjects() ([]*Project, error) {
	repo_dir, err := getRepoDir()
	if err != nil {
		return nil, err
	}

	projects, err := GetProjects(repo_dir)
	if err != nil {
		return nil, err
	}

	return projects, err
}

// checks if user is within list of projects P
// if they are returns project else returs nil
func userInProject(projects []*Project) (*Project, error) {
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

// -------------------------------- Task Utils --------------------------------
// get the user set task_directory defaults to ~/.tasks
func getTaskDir() (string, error) {
	task_dir := os.Getenv("TASK_DIR")
	if task_dir == "" {
		task_dir = "~/.tasks"
	}

	task_dir, err := filepath.Abs(task_dir)
	if err != nil {
		return "", err
	}

	// make the dir if it does not exist
	err = os.MkdirAll(task_dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return task_dir, nil
}

// get all tasks
func getAllTasks() ([]*Task, error) {
	// get our taskdir
	task_dir, err := getTaskDir()
	if err != nil {
		return nil, err
	}
	// there are two files in our task_dir we care about
	// todo.txt, done.txt
	all_tasks := []*Task{}
	file_paths := []string{"todo.txt", "done.txt"}

	for _, file_path := range file_paths {
		file_tasks, err := LoadTaskFile(filepath.Join(task_dir, file_path))
		if err != nil {
			return nil, err
		}
		all_tasks = append(all_tasks, file_tasks...)
	}

	return all_tasks, nil
}

// write all tasks back to the files
func writeAllTasks(tasks []*Task) error {
	// get our taskdir
	task_dir, err := getTaskDir()
	if err != nil {
		return err
	}

	// sort tasks into complete and incomplete
	done := []*Task{}
	not_done := []*Task{}

	for _, task := range tasks {
		if task.Done {
			done = append(done, task)
		} else {
			not_done = append(not_done, task)
		}
	}

	err = WriteTaskFile(filepath.Join(task_dir, "done.txt"), done)
	if err != nil {
		return nil
	}

	err = WriteTaskFile(filepath.Join(task_dir, "not_done.txt"), done)
	if err != nil {
		return nil
	}

	return nil
}
