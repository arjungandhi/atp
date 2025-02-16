package todo

import (
	"bufio"
	"os"
	"regexp"
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

	// Match completion date after the x (e.g., '2025-02-15')
	reCompletionDate := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	completionDateMatch := reCompletionDate.FindStringSubmatch(getXChar(line, 10))
	if len(completionDateMatch) > 0 {
		completionDate, err := time.Parse("2006-01-02", completionDateMatch[1])
		if err == nil {
			todo.CompletionDate = completionDate
		}
		line = strings.Replace(line, completionDateMatch[0], "", 1)
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

	// Match key-value pairs (e.g., due:2025-03-01)
	reKeyValue := regexp.MustCompile(`(\w+):([\w\-\/]+)`)
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
	for key, value := range todo.Labels {
		sb.WriteString(" " + key + ":" + value)
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
	file, err := os.Create(path)
	if err != nil {
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

// gets Next X characters of a string up with a max of the full string
func getXChar(s string, chars int) string {
	min_chars := min(len(s), chars)
	return s[:min_chars]
}
