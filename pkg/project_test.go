package atp

import (
	"path/filepath"
	"testing"
)

func TestNewProject(t *testing.T) {
	// create a new project
	p := NewProject("test", "project", "/tmp", "https://github.cm/test/project")
	if p.Owner != "test" {
		t.Errorf("Expected Owner to be 'test', got %s", p.Owner)
	}
	if p.Name != "project" {
		t.Errorf("Expected Name to be 'project', got %s", p.Name)
	}
	if p.Dir != "/tmp" {
		t.Errorf("Expected Dir to be '/tmp', got %s", p.Dir)
	}
	if p.Url != "https://github.cm/test/project" {
		t.Errorf("Expected Url to be 'https://github.cm/test/project', got %s", p.Url)
	}
}

func TestString(t *testing.T) {
	// create a new project
	p := NewProject("test", "project", "/tmp", "https://github.cm/test/project")
	if p.String() != "test/project" {
		t.Errorf("Expected 'test/project', got %s", p.String())
	}
}

func TestGetProjects(t *testing.T) {
	// get the list of projects
	projects, err := GetProjects("./testdata")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if len(projects) != 4 {
		t.Errorf("Expected 4 projects, got %d", len(projects))
	}

	// check the first project
	expectedOwners := map[string]bool{"arjungandhi": true, "billybob": true}
	expectedNames := map[string]bool{"project1": true, "project2": true}

	actualOwners := map[string]bool{}
	actualNames := map[string]bool{}

	for _, p := range projects {
		actualOwners[p.Owner] = true
		actualNames[p.Name] = true
	}

	for owner, _ := range expectedOwners {
		if _, ok := actualOwners[owner]; !ok {
			t.Errorf("Expected owner %s not found", owner)
		}
	}

	for name, _ := range expectedNames {
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

func TestGetProjectDirs(t *testing.T) {
	// relative path
	project_dirs, err := getProjectDirs("./testdata")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if len(project_dirs) != 4 {
		t.Errorf("Expected 4 project dirs, got %d", len(project_dirs))
	}

	// absolute path
	abs_path, err := filepath.Abs("./testdata")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	project_dirs, err = getProjectDirs(abs_path)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(project_dirs) != 4 {
		t.Errorf("Expected 4 project dirs, got %d", len(project_dirs))
	}
}
