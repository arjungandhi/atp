// repo is a package designed for working with repos on my local machine, it provides access and info about specific repos
package repo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Repo struct {
	Owner string
	Name  string
	Dir   string
	Url   string
}

// Make a new Repo
func NewRepo(owner string, name string, dir string, url string) *Repo {
	return &Repo{
		Owner: owner,
		Name:  name,
		Dir:   dir,
		Url:   url,
	}
}

// --------------------------- Repo Struct Funcs ---------------------------

// Gets the fullpath for the design doc of the repo.
// by default this is the design_doc.md file for in the root of the repo.
func (repo *Repo) getRepoDoc() string {
	// TODO(arjun): future consider adding away to override this
	doc_dir := filepath.Join(repo.Dir, "design_doc.md")
	return doc_dir
}

func (repo *Repo) String() string {
	return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
}

// ------------------------------- Static Funcs -------------------------------

// gets the list of all repos from our local directory
func GetRepos(repo_dir string) ([]*Repo, error) {
	// get the list of all directories in the repos dir
	repo_dirs, err := getRepoDirs(repo_dir)
	if err != nil {
		return nil, err
	}

	// create a list of repos
	var repos []*Repo
	for _, repo_dir := range repo_dirs {
		// split the repo_dir into owner and name
		parts := strings.Split(repo_dir, string(os.PathSeparator))
		owner := parts[len(parts)-2]
		name := parts[len(parts)-1]

		// generate the url
		// TODO: (arjun) hard coding this to assume github.com for now
		url := fmt.Sprintf("https://github.com/%s/%s", owner, name)

		// create a new repo
		p := NewRepo(owner, name, repo_dir, url)
		repos = append(repos, p)
	}

	return repos, nil

}

func getRepoDirs(repo_dir string) ([]string, error) {
	// resolve to an absolute path
	abs_repo_dir, err := filepath.Abs(repo_dir)
	if err != nil {
		return nil, err
	}

	// check the depth of the input path
	inital_depth := strings.Count(abs_repo_dir, string(os.PathSeparator))

	// get the list of all directories in the repos dir
	var repo_dirs []string
	err = filepath.WalkDir(abs_repo_dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// check that there are 3 levels exlude the root dir
			depth := strings.Count(path, string(os.PathSeparator)) - inital_depth

			if depth == 3 {
				repo_dirs = append(repo_dirs, path)
				return filepath.SkipDir
			}

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return repo_dirs, nil

}
