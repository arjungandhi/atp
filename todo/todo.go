package todo

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Todo struct {
	Done           bool
	CreationDate   time.Time
	CompletionDate time.Time
	Priority       string
	Description    string
	Projects       []string
	Contexts       []string
	Labels         map[string]string
}

// Creates a new todo with defaults
func NewTodo() *Todo {
	return &Todo{
		Done:           false,
		CreationDate:   time.Time{},
		CompletionDate: time.Time{},
		Priority:       "",
		Description:    "",
		Projects:       []string{},
		Contexts:       []string{},
		Labels:         map[string]string{},
	}
}

// Parse a todo.txt todo line
func FromString(line string) *Todo {
	todo := NewTodo()

	// Step 1: Parse if the todo is done
	if strings.HasPrefix(line, "x ") {
		todo.Done = true
		line = strings.TrimPrefix(line, "x ")
	}

	// Match completion date after the x (e.g., '2025-02-15') - only for completed todos
	if todo.Done {
		reCompletionDate := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
		completionDateMatch := reCompletionDate.FindStringSubmatch(getXChar(line, 10))
		if len(completionDateMatch) > 0 {
			completionDate, err := time.Parse("2006-01-02", completionDateMatch[1])
			if err == nil {
				todo.CompletionDate = completionDate
			}
			line = strings.Replace(line, completionDateMatch[0], "", 1)
		}
	}

	// Match priority, e.g., (A), (B), etc.
	rePriority := regexp.MustCompile(`\((\w)\)`)
	priorityMatch := rePriority.FindStringSubmatch(getXChar(line, 3))
	if len(priorityMatch) > 0 {
		todo.Priority = priorityMatch[1]
		line = strings.Replace(line, priorityMatch[0], "", 1)
	}

	// Match creation date (e.g., 2025-02-15)
	reCreationDate := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	creationDateMatch := reCreationDate.FindStringSubmatch(getXChar(line, 11))
	if len(creationDateMatch) > 0 {
		creationDate, err := time.Parse("2006-01-02", creationDateMatch[1])
		if err == nil {
			todo.CreationDate = creationDate
		}
		line = strings.Replace(line, creationDateMatch[0], "", 1)
	}

	// Match projects (e.g., +phone)
	reProjects := regexp.MustCompile(`\+(\w+)`)
	projectMatches := reProjects.FindAllStringSubmatch(line, -1)
	for _, match := range projectMatches {
		todo.Projects = append(todo.Projects, match[1])
		line = strings.Replace(line, match[0], "", -1)
	}

	// Match contexts (e.g., @home)
	reContexts := regexp.MustCompile(`@(\w+)`)
	contextMatches := reContexts.FindAllStringSubmatch(line, -1)
	for _, match := range contextMatches {
		todo.Contexts = append(todo.Contexts, match[1])
		line = strings.Replace(line, match[0], "", -1)
	}

	// Match key-value pairs (e.g., due:2025-03-01, url:https://...)
	// Use \S+ (non-whitespace) to capture the full value including URLs
	reKeyValue := regexp.MustCompile(`(\w+):(\S+)`)
	keyValueMatches := reKeyValue.FindAllStringSubmatch(line, -1)
	todo.Labels = make(map[string]string)
	for _, match := range keyValueMatches {
		todo.Labels[match[1]] = match[2]
		line = strings.Replace(line, match[0], "", -1)
	}

	// The remaining part of the line should be the todo description
	todo.Description = strings.TrimSpace(line)

	return todo
}

// Todo object to a todo.txt
func (todo *Todo) String() string {
	var sb strings.Builder

	// Add 'x' if the todo is completed
	if todo.Done {
		sb.WriteString("x ")
	}

	// Add completion date if available
	if !todo.CompletionDate.IsZero() {
		sb.WriteString(todo.CompletionDate.Format("2006-01-02") + " ")
	}

	// Add priority (e.g., (A), (B), etc.)
	if todo.Priority != "" {
		sb.WriteString("(" + todo.Priority + ") ")
	}

	// Add creation date if available
	if !todo.CreationDate.IsZero() {
		sb.WriteString(todo.CreationDate.Format("2006-01-02") + " ")
	}

	// Add description
	sb.WriteString(todo.Description)

	// Add projects (e.g., +phone, +work)
	for _, project := range todo.Projects {
		sb.WriteString(" +" + project)
	}

	// Add contexts (e.g., @home, @office)
	for _, context := range todo.Contexts {
		sb.WriteString(" @" + context)
	}

	// Add extra key-value pairs (e.g., due:2025-03-01)
	// Sort keys for consistent ordering
	keys := make([]string, 0, len(todo.Labels))
	for key := range todo.Labels {
		keys = append(keys, key)
	}
	// Sort with priority order: repo, issue, pr, url, then alphabetical for the rest
	sortLabels(keys)
	for _, key := range keys {
		sb.WriteString(" " + key + ":" + todo.Labels[key])
	}

	return sb.String()
}

// Load a todo.txt file into todos
func LoadTodoFile(path string) ([]*Todo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	todos := []*Todo{}
	for scanner.Scan() {
		line := scanner.Text()
		todo := FromString(line)
		todos = append(todos, todo)
	}

	return todos, nil
}

// Write a todo.txt file from a list of todos
func WriteTodoFile(path string, todos []*Todo) error {
	// if the file exists, back it up to a .bak file
	if _, err := os.Stat(path); err == nil {
		err := os.Rename(path, path+".bak")
		if err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		// if there was an error, restore the backup
		if _, err := os.Stat(path + ".bak"); err == nil {
			os.Rename(path+".bak", path)
		}
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, todo := range todos {
		writer.WriteString(todo.String())
		writer.WriteString("\n")
	}
	writer.Flush()

	return nil
}

// Load todo dir
func LoadTodoDir(path string) ([]*Todo, error) {
	todo_path := ActiveTodoPath(path)
	done_path := DoneTodoPath(path)

	todos, err := LoadTodoFile(todo_path)
	if err != nil {
		return nil, err
	}

	done_todos, err := LoadTodoFile(done_path)
	if err != nil {
		return nil, err
	}

	// append done todos to the list
	todos = append(todos, done_todos...)

	return todos, nil
}

// Write todo dir
func WriteTodoDir(path string, todos []*Todo) error {
	todo_path := ActiveTodoPath(path)
	done_path := DoneTodoPath(path)

	// split the todos into active and done
	active_todos := []*Todo{}
	done_todos := []*Todo{}

	for _, todo := range todos {
		if todo.Done {
			done_todos = append(done_todos, todo)
		} else {
			active_todos = append(active_todos, todo)
		}
	}

	err := WriteTodoFile(todo_path, active_todos)
	if err != nil {
		return err
	}

	err = WriteTodoFile(done_path, done_todos)
	if err != nil {
		return err
	}

	return nil
}

// todo file path
func ActiveTodoPath(dir string) string {
	return filepath.Join(dir, "todo.txt")
}

// done file path
func DoneTodoPath(dir string) string {
	return filepath.Join(dir, "done.txt")
}

// gets Next X characters of a string up with a max of the full string
func getXChar(s string, chars int) string {
	min_chars := min(len(s), chars)
	return s[:min_chars]
}

// sortLabels sorts label keys with a priority order for GitHub-related labels
func sortLabels(keys []string) {
	sort.Slice(keys, func(i, j int) bool {
		// Priority order: repo, issue, pr, url
		priority := map[string]int{
			"repo":  1,
			"issue": 2,
			"pr":    2, // issue and pr have same priority (mutually exclusive)
			"url":   3,
		}

		pi, oki := priority[keys[i]]
		pj, okj := priority[keys[j]]

		// If both have priority, sort by priority
		if oki && okj {
			if pi != pj {
				return pi < pj
			}
			// If same priority, sort alphabetically
			return keys[i] < keys[j]
		}

		// If only one has priority, it comes first
		if oki {
			return true
		}
		if okj {
			return false
		}

		// Neither has priority, sort alphabetically
		return keys[i] < keys[j]
	})
}
