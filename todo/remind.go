package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// LoadReminderTasks loads reminder tasks from reminders.txt file
func LoadReminderTasks(path string) ([]*Todo, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Todo{}, nil // Return empty slice if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	var reminders []*Todo
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		todo := FromString(line)
		reminders = append(reminders, todo)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

// WriteReminderTasks writes reminder tasks to reminders.txt file
func WriteReminderTasks(path string, reminders []*Todo) error {
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

	for _, reminder := range reminders {
		if _, err := writer.WriteString(reminder.String() + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// ReminderTasksPath returns the path to reminders.txt file in a given directory
func ReminderTasksPath(dir string) string {
	return filepath.Join(dir, "reminders.txt")
}

// GetDueReminders returns reminders that are due on or before the given date
func GetDueReminders(reminders []*Todo, date time.Time) []*Todo {
	var dueReminders []*Todo
	dateStr := date.Format("2006-01-02")
	
	for _, reminder := range reminders {
		if remindDate, exists := reminder.Labels["remind"]; exists {
			if remindDate <= dateStr {
				dueReminders = append(dueReminders, reminder)
			}
		}
	}
	
	return dueReminders
}

// RemoveDueReminders removes the due reminders from the reminder list
func RemoveDueReminders(reminders []*Todo, dueReminders []*Todo) []*Todo {
	// Create a map for quick lookup of due reminders
	dueMap := make(map[string]bool)
	for _, due := range dueReminders {
		dueMap[due.String()] = true
	}
	
	var remaining []*Todo
	for _, reminder := range reminders {
		if !dueMap[reminder.String()] {
			remaining = append(remaining, reminder)
		}
	}
	
	return remaining
}

// ProcessReminders moves due reminders from reminders.txt to active todos
func ProcessReminders(todoDir string, date time.Time) error {
	reminderPath := ReminderTasksPath(todoDir)
	
	// Load existing reminders
	reminders, err := LoadReminderTasks(reminderPath)
	if err != nil {
		return err
	}
	
	if len(reminders) == 0 {
		return nil // No reminders to process
	}
	
	// Get due reminders
	dueReminders := GetDueReminders(reminders, date)
	if len(dueReminders) == 0 {
		return nil // No due reminders
	}
	
	// Create copies of due reminders with remind labels removed and creation date set
	processedTodos := make([]*Todo, len(dueReminders))
	for i, reminder := range dueReminders {
		// Create a copy of the reminder
		processedTodo := &Todo{
			Done:           reminder.Done,
			CreationDate:   date,
			CompletionDate: reminder.CompletionDate,
			Priority:       reminder.Priority,
			Description:    reminder.Description,
			Projects:       make([]string, len(reminder.Projects)),
			Contexts:       make([]string, len(reminder.Contexts)),
			Labels:         make(map[string]string),
		}
		
		// Copy projects and contexts
		copy(processedTodo.Projects, reminder.Projects)
		copy(processedTodo.Contexts, reminder.Contexts)
		
		// Copy labels except remind
		for k, v := range reminder.Labels {
			if k != "remind" {
				processedTodo.Labels[k] = v
			}
		}
		
		processedTodos[i] = processedTodo
	}
	
	// Load existing todos
	existingTodos, err := LoadTodoDir(todoDir)
	if err != nil {
		// If we can't load, start with empty slice
		existingTodos = []*Todo{}
	}
	
	// Add processed todos to active todos
	allTodos := append(existingTodos, processedTodos...)
	
	// Remove due reminders from reminder list
	remainingReminders := RemoveDueReminders(reminders, dueReminders)
	
	// Write updated files
	if err := WriteTodoDir(todoDir, allTodos); err != nil {
		return fmt.Errorf("failed to write todos: %w", err)
	}
	
	if err := WriteReminderTasks(reminderPath, remainingReminders); err != nil {
		return fmt.Errorf("failed to write reminders: %w", err)
	}
	
	return nil
}

// SortRemindersByDate sorts reminders by their remind date
func SortRemindersByDate(reminders []*Todo) {
	sort.Slice(reminders, func(i, j int) bool {
		dateI := reminders[i].Labels["remind"]
		dateJ := reminders[j].Labels["remind"]
		return dateI < dateJ
	})
}

// AddReminderTask adds a new reminder task to reminders.txt
func AddReminderTask(todoDir string, reminder *Todo) error {
	reminderPath := ReminderTasksPath(todoDir)
	
	// Load existing reminders
	reminders, err := LoadReminderTasks(reminderPath)
	if err != nil {
		return err
	}
	
	// Add new reminder
	reminders = append(reminders, reminder)
	
	// Write updated reminders
	return WriteReminderTasks(reminderPath, reminders)
}