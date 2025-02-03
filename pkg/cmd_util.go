package atp

import (
	"errors"
	"os"
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
