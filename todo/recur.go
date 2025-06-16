package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type RecurringTask struct {
	Schedule    cron.Schedule
	ScheduleStr string // Keep original string for display/serialization
	Todo        *Todo
}

func NewRecurringTask(scheduleStr string, todo *Todo) (*RecurringTask, error) {
	schedule, err := parseSchedule(scheduleStr)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule '%s': %w", scheduleStr, err)
	}
	
	return &RecurringTask{
		Schedule:    schedule,
		ScheduleStr: scheduleStr,
		Todo:        todo,
	}, nil
}

// parseSchedule converts schedule string to cron.Schedule
func parseSchedule(scheduleStr string) (cron.Schedule, error) {
	// Handle simple formats first - only allow specific ones we support
	switch scheduleStr {
	case "@daily":
		return cron.ParseStandard("0 0 * * *")
	case "@weekly":
		return cron.ParseStandard("0 0 * * 1") // Monday
	case "@monthly":
		return cron.ParseStandard("0 0 1 * *") // 1st of month
	default:
		// Check if it's an unsupported @ format
		if strings.HasPrefix(scheduleStr, "@") {
			return nil, fmt.Errorf("unsupported schedule format '%s' (supported: @daily, @weekly, @monthly)", scheduleStr)
		}
		// Try parsing as standard cron expression
		return cron.ParseStandard(scheduleStr)
	}
}

// Parse a recurring task line from recur.txt
// Format: @daily Task description +project @context key:value
//         0 0 * * * Task description +project @context key:value
func RecurringTaskFromString(line string) (*RecurringTask, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, nil
	}

	var scheduleStr string
	var todoText string

	// Check if it's a simple format (@daily, @weekly, @monthly)
	if strings.HasPrefix(line, "@") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid simple schedule format: '%s' (expected format: @daily Task description)", line)
		}
		scheduleStr = parts[0]
		todoText = parts[1]
	} else {
		// Handle cron format (5 fields: minute hour day month weekday)
		parts := strings.Fields(line)
		if len(parts) < 6 { // At least 5 cron fields + task description
			return nil, fmt.Errorf("invalid cron format: '%s' (expected format: minute hour day month weekday Task description)", line)
		}
		
		// First 5 parts are the cron schedule
		scheduleStr = strings.Join(parts[0:5], " ")
		// Rest is the todo text
		todoText = strings.Join(parts[5:], " ")
	}

	// Parse the todo part using existing todo parser
	todo := FromString(todoText)

	// Create RecurringTask using NewRecurringTask which validates the schedule
	return NewRecurringTask(scheduleStr, todo)
}

// Convert RecurringTask to string format for recur.txt
func (rt *RecurringTask) String() string {
	return rt.ScheduleStr + " " + rt.Todo.String()
}

// Check if a recurring task should generate a todo for the given date
func (rt *RecurringTask) ShouldGenerateForDate(date time.Time) bool {
	// Get the previous scheduled time before the given date
	// and the next scheduled time after that
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	previous := rt.Schedule.Next(startOfDay.Add(-24 * time.Hour))
	
	// Check if the previous scheduled time falls on the given date
	return previous.Year() == date.Year() && 
		   previous.Month() == date.Month() && 
		   previous.Day() == date.Day()
}




// Generate a Todo from this recurring task for a specific date
func (rt *RecurringTask) GenerateTodo(date time.Time) *Todo {
	// Create a copy of the template todo
	newTodo := &Todo{
		Done:           false,
		CreationDate:   time.Time{}, // Don't set creation date - use recur: label instead
		CompletionDate: time.Time{},
		Priority:       rt.Todo.Priority,
		Description:    rt.Todo.Description,
		Projects:       make([]string, len(rt.Todo.Projects)),
		Contexts:       make([]string, len(rt.Todo.Contexts)),
		Labels:         make(map[string]string),
	}

	// Copy projects and contexts
	copy(newTodo.Projects, rt.Todo.Projects)
	copy(newTodo.Contexts, rt.Todo.Contexts)

	// Copy labels
	for k, v := range rt.Todo.Labels {
		newTodo.Labels[k] = v
	}

	// Add recur metadata to track this was generated
	newTodo.Labels["recur"] = date.Format("2006-01-02")

	return newTodo
}

// Load recurring tasks from recur.txt file
func LoadRecurringTasks(path string) ([]*RecurringTask, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*RecurringTask{}, nil // Return empty slice if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	var tasks []*RecurringTask
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		task, err := RecurringTaskFromString(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line %d: %w", lineNum, err)
		}
		if task != nil {
			tasks = append(tasks, task)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Write recurring tasks to recur.txt file
func WriteRecurringTasks(path string, tasks []*RecurringTask) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Backup existing file
	if _, err := os.Stat(path); err == nil {
		if err := os.Rename(path, path+".bak"); err != nil {
			return fmt.Errorf("failed to backup existing file: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		// Restore backup if creation failed
		if _, err := os.Stat(path + ".bak"); err == nil {
			os.Rename(path+".bak", path)
		}
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, task := range tasks {
		if _, err := writer.WriteString(task.String() + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Path to recur.txt file in a given directory
func RecurringTasksPath(dir string) string {
	return filepath.Join(dir, "recur.txt")
}

// Generate todos from recurring tasks for a specific date, avoiding duplicates
func GenerateTodosFromRecurring(todoDir string, date time.Time) ([]*Todo, error) {
	recurPath := RecurringTasksPath(todoDir)
	tasks, err := LoadRecurringTasks(recurPath)
	if err != nil {
		return nil, err
	}

	// Load existing todos to check for duplicates
	existingTodos, err := LoadTodoDir(todoDir)
	if err != nil {
		// If we can't load existing todos, continue anyway
		existingTodos = []*Todo{}
	}

	var newTodos []*Todo
	dateStr := date.Format("2006-01-02")

	for _, task := range tasks {
		if task.ShouldGenerateForDate(date) {
			// Check if we already generated this task for this date
			if !todoExistsForRecurringTask(existingTodos, task.Todo.Description, dateStr) {
				todo := task.GenerateTodo(date)
				newTodos = append(newTodos, todo)
			}
		}
	}

	return newTodos, nil
}

// Check if a todo was already generated for a specific recurring task and date
func todoExistsForRecurringTask(todos []*Todo, description string, dateStr string) bool {
	for _, todo := range todos {
		if todo.Description == description {
			if recur, exists := todo.Labels["recur"]; exists && recur == dateStr {
				return true
			}
		}
	}
	return false
}

// Add generated todos to the todo directory
func AddRecurringTodosToDir(todoDir string, date time.Time) error {
	newTodos, err := GenerateTodosFromRecurring(todoDir, date)
	if err != nil {
		return err
	}

	if len(newTodos) == 0 {
		return nil // No new todos to add
	}

	// Load existing todos
	existingTodos, err := LoadTodoDir(todoDir)
	if err != nil {
		// If we can't load, start with empty slice
		existingTodos = []*Todo{}
	}

	// Append new todos
	allTodos := append(existingTodos, newTodos...)

	// Write back to directory
	return WriteTodoDir(todoDir, allTodos)
}