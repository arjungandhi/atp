package project

import (
	"fmt"
	"slices"
	"strings"

	"github.com/arjungandhi/atp/repo"
	"github.com/arjungandhi/atp/todo"
)

type Project struct {
	Name      string
	Phase     string
	Active    bool
	Done      bool
	Repo      *repo.Repo
	todo_data *todo.Todo
}

func NewProject() *Project {
	return &Project{
		Name:   "",
		Phase:  "",
		Active: false,
		Done:   false,
		Repo:   nil,
	}
}

// converts a project struct to a todo struct
func (p *Project) ToTodo() *todo.Todo {
	t := p.todo_data
	t.Description = p.Name
	if p.Phase != "" {
		t.Labels["phase"] = p.Phase
	}

	if p.Active {
		t.Priority = "A"
	} else {
		t.Priority = ""
	}

	t.Done = p.Done

	if p.Repo != nil {
		t.Labels["repo"] = fmt.Sprintf("%s/%s", p.Repo.Owner, p.Repo.Name)
	}

	return t
}

// convert from a todo struct -> a project
func FromTodo(t *todo.Todo, repos []*repo.Repo) (*Project, error) {
	p := NewProject()

	p.Name = t.Description
	val, ok := t.Labels["phase"]
	if ok {
		p.Phase = val
	}

	if t.Priority == "A" {
		p.Active = true
	}

	p.Done = t.Done

	// find a repo that has the same owner and name in the repos
	val, ok = t.Labels["repo"]
	if ok {
		vals := strings.Split(t.Labels["repo"], "/")
		if len(vals) != 2 {
			return nil, fmt.Errorf("Invalid repo label %s", t.Labels["repo"])
		}
		owner, name := vals[0], vals[1]

		// find the repo which matches owner/name
		for _, repo := range repos {
			if repo.Owner == owner && repo.Name == name {
				p.Repo = repo
				break
			}
		}
	}

	// add extra todo data
	p.todo_data = t

	return p, nil
}

// convert a project to a string
func (p *Project) String() string {
	// convert the project to a todo
	return p.Name
}

// convert a project to a string with more details
func (p *Project) TodoString() string {
	// convert the project to a todo
	return p.ToTodo().String()
}

// Load a Project Dir
func LoadProjectsDir(dir_path string, repos []*repo.Repo) ([]*Project, error) {
	// load the project file and the done file if it exists
	active_path := ActiveFilePath(dir_path)
	done_path := DoneFilePath(dir_path)
	projects := []*Project{}
	// load the active projects
	active_projects, err := LoadProjectFile(active_path, repos)
	if err != nil {
		return nil, fmt.Errorf("Failed to load active projects %s, %w", active_path, err)
	}

	// load the done projects
	done_projects, err := LoadProjectFile(done_path, repos)
	if err != nil {
		return nil, fmt.Errorf("Failed to load done projects %s, %w", done_path, err)
	}

	// combine the projects
	projects = append(active_projects, done_projects...)

	return projects, nil
}

// Load a Project file
func LoadProjectFile(path string, repos []*repo.Repo) ([]*Project, error) {
	// load the file
	todos, err := todo.LoadTodoFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load project file %s, %w", path, err)
	}

	projects := []*Project{}
	for _, t := range todos {
		p, err := FromTodo(t, repos)
		if err != nil {
			return nil, fmt.Errorf("Failed to load todo -> project %s, %w", t.Description, err)
		}

		projects = append(projects, p)

	}

	return projects, nil
}

// Write projects to a projecte dir
func WriteProjectsDir(dir_path string, projects []*Project) error {
	// sort the projects
	SortProjects(projects)

	// split the projects into active and done
	active_path := ActiveFilePath(dir_path)
	done_path := DoneFilePath(dir_path)
	active_projects := []*Project{}
	done_projects := []*Project{}

	for _, p := range projects {
		if p.Done {
			done_projects = append(done_projects, p)
		} else {
			active_projects = append(active_projects, p)
		}
	}

	// write the active projects
	err := WriteProjectFile(active_path, active_projects)
	if err != nil {
		return fmt.Errorf("Failed to write active projects %s, %w", active_path, err)
	}

	// write the done projects
	err = WriteProjectFile(done_path, done_projects)
	if err != nil {
		return fmt.Errorf("Failed to write done projects %s, %w", done_path, err)
	}

	return nil
}

// write projects to a file
func WriteProjectFile(path string, projects []*Project) error {
	// convert the projects to todos
	todos := []*todo.Todo{}
	for _, p := range projects {
		t := p.ToTodo()
		todos = append(todos, t)
	}
	// write the file
	return todo.WriteTodoFile(path, todos)

}

// sort a list of projects
func SortProjects(projects []*Project) {
	// sort the projects
	slices.SortFunc(projects, func(a *Project, b *Project) int {
		// sort by active, then by name
		if a.Active && !b.Active {
			return -1
		}
		if !a.Active && b.Active {
			return 1
		}

		return strings.Compare(a.Name, b.Name)
	})
}

// get active file path
func ActiveFilePath(dir_path string) string {
	// get the active file
	active_path := fmt.Sprintf("%s/projects.txt", dir_path)

	return active_path
}

// get done file
func DoneFilePath(dir_path string) string {
	// get the done file
	done_path := fmt.Sprintf("%s/done.txt", dir_path)

	return done_path
}
