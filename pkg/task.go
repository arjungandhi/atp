package atp

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Task is a superset around at todo.txt task, all atp tasks are todo.txt tasks
type Task struct {
	Done        bool
	Priority    string
	Description string
	Projects    []string
	Contexts    []string
	Labels      map[string]string
}

// Creates a new task with defaults
func NewTask() *Task {
	return &Task{
		Done:        false,
		Priority:    "",
		Description: "",
		Projects:    []string{},
		Contexts:    []string{},
		Labels:      map[string]string{},
	}
}

// Parse a todo.txt task line
func FromString(line string) *Task {
	task := NewTask()

	// Step 1: Parse if the task is done
	if strings.HasPrefix(line, "x ") {
		task.Done = true
		line = strings.TrimPrefix(line, "x ")
	}

	// Match priority, e.g., (A), (B), etc.
	rePriority := regexp.MustCompile(`\((\w)\)`)
	priorityMatch := rePriority.FindStringSubmatch(line)
	if len(priorityMatch) > 0 {
		task.Priority = priorityMatch[1]
		line = strings.Replace(line, priorityMatch[0], "", 1)
	}

	// Match projects (e.g., +phone)
	reProjects := regexp.MustCompile(`\+(\w+)`)
	projectMatches := reProjects.FindAllStringSubmatch(line, -1)
	for _, match := range projectMatches {
		task.Projects = append(task.Projects, match[1])
		line = strings.Replace(line, match[0], "", -1)
	}

	// Match contexts (e.g., @home)
	reContexts := regexp.MustCompile(`@(\w+)`)
	contextMatches := reContexts.FindAllStringSubmatch(line, -1)
	for _, match := range contextMatches {
		task.Contexts = append(task.Contexts, match[1])
		line = strings.Replace(line, match[0], "", -1)
	}

	// Match key-value pairs (e.g., priority:high, location:office)
	reKeyValue := regexp.MustCompile(`(\w+):([\w\-\/]+)`)
	keyValueMatches := reKeyValue.FindAllStringSubmatch(line, -1)
	task.Labels = make(map[string]string)
	for _, match := range keyValueMatches {
		task.Labels[match[1]] = match[2]
		line = strings.Replace(line, match[0], "", -1)
	}

	// The remaining part of the line should be the task description
	task.Description = strings.TrimSpace(line)

	return task
}

// Task object to a todo.txt
func (task *Task) String() string {
	var sb strings.Builder

	// Add 'x' if the task is completed
	if task.Done {
		sb.WriteString("x ")
	}

	// Add priority (e.g., (A), (B), etc.)
	if task.Priority != "" {
		sb.WriteString("(" + task.Priority + ") ")
	}

	// Add description
	sb.WriteString(task.Description)

	// Add projects (e.g., +phone, +work)
	for _, project := range task.Projects {
		sb.WriteString(" +" + project)
	}

	// Add contexts (e.g., @home, @office)
	for _, context := range task.Contexts {
		sb.WriteString(" @" + context)
	}

	// Add extra key-value pairs
	for key, value := range task.Labels {
		sb.WriteString(" " + key + ":" + value)
	}

	return sb.String()

}

// Load a todo.txt file into tasks
func LoadTaskFile(path string) ([]*Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tasks := []*Task{}
	for scanner.Scan() {
		line := scanner.Text()
		task := FromString(line)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Write a todo.txt file from a list of tasks
func WriteTaskFile(path string, tasks []*Task) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, task := range tasks {
		writer.WriteString(task.String())
		writer.WriteString("\n")
	}
	writer.Flush()

	return nil
}
