package project

import (
	"os"
	"testing"

	"github.com/arjungandhi/atp/repo"
	"github.com/arjungandhi/atp/todo"
)

func TestNewProject(t *testing.T) {
	p := NewProject()

	if p.Name != "" {
		t.Errorf("Expected empty Name, got %s", p.Name)
	}
	if p.Phase != "" {
		t.Errorf("Expected empty Phase, got %s", p.Phase)
	}
	if p.Active != false {
		t.Errorf("Expected Active to be false, got %v", p.Active)
	}
	if p.Done != false {
		t.Errorf("Expected Done to be false, got %v", p.Done)
	}
	if p.Repo != nil {
		t.Errorf("Expected Repo to be nil, got %v", p.Repo)
	}
	if p.todo_data != nil {
		t.Errorf("Expected todo_data to be nil, got %v", p.todo_data)
	}
}

func TestToTodo(t *testing.T) {
	repo := &repo.Repo{
		Owner: "testOwner",
		Name:  "testRepo",
	}

	p := &Project{
		Name:   "Test Project",
		Phase:  "Phase 1",
		Active: true,
		Done:   true,
		Repo:   repo,
		todo_data: &todo.Todo{
			Description: "Test Project Todo",
			Labels:      make(map[string]string),
		},
	}

	todoResult := p.ToTodo()

	if todoResult.Description != p.Name {
		t.Errorf("Expected Description %s, got %s", p.Name, todoResult.Description)
	}
	if todoResult.Labels["phase"] != p.Phase {
		t.Errorf("Expected phase label %s, got %s", p.Phase, todoResult.Labels["phase"])
	}
	if todoResult.Priority != "A" {
		t.Errorf("Expected priority A, got %s", todoResult.Priority)
	}
	if todoResult.Done != p.Done {
		t.Errorf("Expected Done to be %v, got %v", p.Done, todoResult.Done)
	}
	if todoResult.Labels["repo"] != "testOwner/testRepo" {
		t.Errorf("Expected repo label testOwner/testRepo, got %s", todoResult.Labels["repo"])
	}
}

func TestFromTodo(t *testing.T) {
	repos := []*repo.Repo{
		{
			Owner: "testOwner",
			Name:  "testRepo",
		},
	}

	todoItem := &todo.Todo{
		Description: "Test Project",
		Priority:    "A",
		Done:        true,
		Labels: map[string]string{
			"phase": "Phase 1",
			"repo":  "testOwner/testRepo",
		},
	}

	project, err := FromTodo(todoItem, repos)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if project.Name != todoItem.Description {
		t.Errorf("Expected Name %s, got %s", todoItem.Description, project.Name)
	}
	if project.Phase != "Phase 1" {
		t.Errorf("Expected Phase 'Phase 1', got %s", project.Phase)
	}
	if !project.Active {
		t.Errorf("Expected Active to be true, got false")
	}
	if project.Done != todoItem.Done {
		t.Errorf("Expected Done to be %v, got %v", todoItem.Done, project.Done)
	}
	if project.Repo == nil || project.Repo.Owner != "testOwner" || project.Repo.Name != "testRepo" {
		t.Errorf("Expected Repo with Owner testOwner and Name testRepo, got %v", project.Repo)
	}
}

func TestLoadProjectFile(t *testing.T) {
	// Assuming there's a mock path to the file
	temp_file, err := os.CreateTemp("", "test_project_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	// Close the file after creating
	temp_file.Close()
	// defer a deleteion
	defer os.Remove(temp_file.Name())

	repos := []*repo.Repo{
		{
			Owner: "testOwner",
			Name:  "testRepo",
		},
	}

	// Create a mock todo file
	err = todo.WriteTodoFile(temp_file.Name(), []*todo.Todo{
		{
			Description: "Test Project 1",
			Priority:    "A",
			Done:        false,
			Labels: map[string]string{
				"phase": "1",
				"repo":  "testOwner/testRepo",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to write mock todo file: %v", err)
	}

	// Load the projects
	projects, err := LoadProjectFile(temp_file.Name(), repos)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(projects) == 0 {
		t.Errorf("Expected at least one project, got %d", len(projects))
	}

	if projects[0].Name != "Test Project 1" {
		t.Errorf("Expected project name 'Test Project 1', got %s", projects[0].Name)
	}
}

func TestWriteProjectFile(t *testing.T) {
	// Prepare mock data
	projects := []*Project{
		{
			Name:   "Test Project",
			Phase:  "1",
			Active: true,
			Done:   false,
			Repo:   nil,
			todo_data: &todo.Todo{
				Description: "Test Project",
				Labels:      make(map[string]string),
			},
		},
	}

	temp_file, err := os.CreateTemp("", "test_project_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	// Close the file after creating
	temp_file.Close()
	// defer a deleteion
	defer os.Remove(temp_file.Name())

	err = WriteProjectFile(temp_file.Name(), projects)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Load the written file to check its content
	todos, err := todo.LoadTodoFile(temp_file.Name())
	if err != nil {
		t.Errorf("Failed to load written todo file: %v", err)
	}

	if len(todos) == 0 {
		t.Errorf("Expected at least one todo, got %d", len(todos))
	}

	if todos[0].Description != "Test Project" {
		t.Errorf("Expected todo description 'Test Project', got %s", todos[0].Description)
	}
}
