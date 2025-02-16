package repo

import (
	"path/filepath"
	"testing"
)

func TestNewRepo(t *testing.T) {
	// create a new repo
	p := NewRepo("test", "repo", "/tmp", "https://github.cm/test/repo")
	if p.Owner != "test" {
		t.Errorf("Expected Owner to be 'test', got %s", p.Owner)
	}
	if p.Name != "repo" {
		t.Errorf("Expected Name to be 'repo', got %s", p.Name)
	}
	if p.Dir != "/tmp" {
		t.Errorf("Expected Dir to be '/tmp', got %s", p.Dir)
	}
	if p.Url != "https://github.cm/test/repo" {
		t.Errorf("Expected Url to be 'https://github.cm/test/repo', got %s", p.Url)
	}
}

func TestString(t *testing.T) {
	// create a new repo
	p := NewRepo("test", "repo", "/tmp", "https://github.cm/test/repo")
	if p.String() != "test/repo" {
		t.Errorf("Expected 'test/repo', got %s", p.String())
	}
}

func TestGetRepos(t *testing.T) {
	// get the list of repos
	repos, err := GetRepos("./testdata")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if len(repos) != 4 {
		t.Errorf("Expected 4 repos, got %d", len(repos))
	}

	// check the first repo
	expectedOwners := map[string]bool{"arjungandhi": true, "billybob": true}
	expectedNames := map[string]bool{"repo1": true, "repo2": true}

	actualOwners := map[string]bool{}
	actualNames := map[string]bool{}

	for _, p := range repos {
		actualOwners[p.Owner] = true
		actualNames[p.Name] = true
	}

	for owner := range expectedOwners {
		if _, ok := actualOwners[owner]; !ok {
			t.Errorf("Expected owner %s not found", owner)
		}
	}

	for name := range expectedNames {
		if _, ok := actualNames[name]; !ok {
			t.Errorf("Expected name %s not found", name)
		}
	}

	// assert the lengths of the maps are same
	if len(expectedOwners) != len(actualOwners) {
		t.Errorf("Expected %d owners, got %d", len(expectedOwners), len(actualOwners))
	}
	if len(expectedNames) != len(actualNames) {
		t.Errorf("Expected %d names, got %d", len(expectedNames), len(actualNames))
	}

}

func TestGetRepoDirs(t *testing.T) {
	// relative path
	repo_dirs, err := getRepoDirs("./testdata")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if len(repo_dirs) != 4 {
		t.Errorf("Expected 4 repo dirs, got %d", len(repo_dirs))
	}

	// absolute path
	abs_path, err := filepath.Abs("./testdata")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	repo_dirs, err = getRepoDirs(abs_path)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(repo_dirs) != 4 {
		t.Errorf("Expected 4 repo dirs, got %d", len(repo_dirs))
	}
}
