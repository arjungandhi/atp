package atp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Owner string
	Name  string
	Dir   string
	Url   string
}

// Make a new Project
func NewProject(owner string, name string, dir string, url string) *Project {
	return &Project{
		Owner: owner,
		Name:  name,
		Dir:   dir,
		Url:   url,
	}
}

// --------------------------- Project Struct Funcs ---------------------------

// Gets the fullpath for the design doc of the project.
// by default this is the design_doc.md file for in the root of the repo.
func (project *Project) getProjectDoc() string {
	// TODO(arjun): future consider adding away to override this
	doc_dir := filepath.Join(project.Dir, "design_doc.md")
	return doc_dir
}

func (project *Project) String() string {
	return fmt.Sprintf("%s/%s", project.Owner, project.Name)
}

// ------------------------------- Static Funcs -------------------------------

// gets the list of all projects from our local directory
func GetProjects(repo_dir string) ([]*Project, error) {
	// get the list of all directories in the repos dir
	project_dirs, err := getProjectDirs(repo_dir)
	if err != nil {
		return nil, err
	}

	// create a list of projects
	var projects []*Project
	for _, project_dir := range project_dirs {
		// split the project_dir into owner and name
		parts := strings.Split(project_dir, string(os.PathSeparator))
		owner := parts[len(parts)-2]
		name := parts[len(parts)-1]

		// generate the url
		// TODO: (arjun) hard coding this to assume github.com for now
		url := fmt.Sprintf("https://github.com/%s/%s", owner, name)

		// create a new project
		p := NewProject(owner, name, project_dir, url)
		projects = append(projects, p)
	}

	return projects, nil

}

func getProjectDirs(repo_dir string) ([]string, error) {
	// resolve to an absolute path
	abs_repo_dir, err := filepath.Abs(repo_dir)
	if err != nil {
		return nil, err
	}

	// check the depth of the input path
	inital_depth := strings.Count(abs_repo_dir, string(os.PathSeparator))

	// get the list of all directories in the repos dir
	var project_dirs []string
	err = filepath.WalkDir(abs_repo_dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// check that there are 3 levels exlude the root dir
			depth := strings.Count(path, string(os.PathSeparator)) - inital_depth

			if depth == 3 {
				project_dirs = append(project_dirs, path)
				return filepath.SkipDir
			}

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return project_dirs, nil

}
