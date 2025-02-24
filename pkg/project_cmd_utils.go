package atp

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// Get project file path
func getProjectPath() (string, error) {
	// get ATP dir
	atp_dir, err := getAtpDir()
	if err != nil {
		return "", fmt.Errorf(
			"Could not get atp dir: %w", err,
		)
	}

	// append projects
	path := filepath.Join(atp_dir, "project")
	// ensure the file exists, make it if it does not
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to ensure existance of file: %w", err,
		)
	}

	// append projects.txt
	path = filepath.Join(path, "projects.txt")

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

// Get done project file path
func getDoneProjectPath() (string, error) {
	// get ATP dir
	atp_dir, err := getAtpDir()
	if err != nil {
		return "", fmt.Errorf(
			"Could not get atp dir: %w", err,
		)
	}

	// append projects
	path := filepath.Join(atp_dir, "project")
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

// get projects
func getProjects() ([]*Project, error) {
	path, err := getProjectPath()

	if err != nil {
		return nil, fmt.Errorf("Unable to load project file into projects: %w", err)
	}

	repos, err := getRepos()
	if err != nil {
		return nil, fmt.Errorf("Unable to repos: %w", err)
	}

	return LoadProjectFile(path, repos)
}

// get done projects
func getDoneProjects() ([]*Project, error) {
	path, err := getDoneProjectPath()

	if err != nil {
		return nil, fmt.Errorf("Unable to load done project file into projects: %w", err)
	}

	repos, err := getRepos()
	if err != nil {
		return nil, fmt.Errorf("Unable to repos: %w", err)
	}

	return LoadProjectFile(path, repos)
}

// get all projects
func getAllProjects() ([]*Project, error) {
	projects, err := getProjects()
	if err != nil {
		return nil, err
	}

	done_projects, err := getDoneProjects()
	if err != nil {
		return nil, err
	}

	return slices.Concat(projects, done_projects), nil

}

// get just the projects that are active
func getActiveProjects() ([]*Project, error) {
	projects, err := getProjects()
	if err != nil {
		return nil, err
	}

	active_projects := []*Project{}
	for _, p := range projects {
		if p.Active {
			active_projects = append(active_projects, p)
		}
	}

	return active_projects, nil
}

// write all projects to the file
func WriteAllProjects(projects []*Project) error {
	path, err := getProjectPath()
	if err != nil {
		return fmt.Errorf("Unable to write projects to file: %w", err)
	}

	done_path, err := getDoneProjectPath()
	if err != nil {
		return fmt.Errorf("Unable to write projects to file: %w", err)
	}

	// separate the projects into done and not done
	done_projects := []*Project{}
	not_done_projects := []*Project{}
	for _, p := range projects {
		if p.Done {
			done_projects = append(done_projects, p)
		} else {
			not_done_projects = append(not_done_projects, p)
		}
	}

	// write the projects to the file
	err = WriteProject(path, not_done_projects)
	if err != nil {
		return fmt.Errorf("Unable to write projects to file: %w", err)
	}

	// write the done projects to the file
	err = WriteProject(done_path, done_projects)
	if err != nil {
		return fmt.Errorf("Unable to write projects to file: %w", err)
	}

	return nil
}
