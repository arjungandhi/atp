package atp

import (
	"fmt"
	"slices"
	"strings"

	"github.com/arjungandhi/atp/pkg/repo"
	"github.com/arjungandhi/atp/pkg/todo"
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

// Load a file contianing projects
func LoadProjectFile(path string, repos []*repo.Repo) ([]*Project, error) {
	todos, err := todo.LoadTodoFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load file %w", err)
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

// write projects to a file
func WriteProjectFile(path string, projects []*Project) error {
	// convert all projects todos
	todos := []*todo.Todo{}
	for _, p := range projects {
		todos = append(todos, p.ToTodo())
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
