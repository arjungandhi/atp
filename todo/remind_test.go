package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadReminderTasks(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	reminderPath := ReminderTasksPath(tempDir)

	// Test loading from non-existent file
	reminders, err := LoadReminderTasks(reminderPath)
	if err != nil {
		t.Fatalf("Expected no error loading non-existent file, got: %v", err)
	}
	if len(reminders) != 0 {
		t.Fatalf("Expected empty slice, got %d reminders", len(reminders))
	}

	// Create test reminder content
	content := `evaluate if system is still working remind:2025-07-15
cancel renters insurance remind:2025-11-01 @home +personal
(A) important reminder remind:2025-06-01 +work`

	// Write test file
	err = os.WriteFile(reminderPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Load reminders
	reminders, err = LoadReminderTasks(reminderPath)
	if err != nil {
		t.Fatalf("Failed to load reminders: %v", err)
	}

	if len(reminders) != 3 {
		t.Fatalf("Expected 3 reminders, got %d", len(reminders))
	}

	// Check first reminder
	if reminders[0].Description != "evaluate if system is still working" {
		t.Errorf("Expected description 'evaluate if system is still working', got '%s'", reminders[0].Description)
	}
	if reminders[0].Labels["remind"] != "2025-07-15" {
		t.Errorf("Expected remind date '2025-07-15', got '%s'", reminders[0].Labels["remind"])
	}

	// Check second reminder
	if len(reminders[1].Contexts) != 1 || reminders[1].Contexts[0] != "home" {
		t.Errorf("Expected context @home, got %v", reminders[1].Contexts)
	}
	if len(reminders[1].Projects) != 1 || reminders[1].Projects[0] != "personal" {
		t.Errorf("Expected project +personal, got %v", reminders[1].Projects)
	}

	// Check third reminder
	if reminders[2].Priority != "A" {
		t.Errorf("Expected priority A, got '%s'", reminders[2].Priority)
	}
}

func TestWriteReminderTasks(t *testing.T) {
	tempDir := t.TempDir()
	reminderPath := ReminderTasksPath(tempDir)

	// Create test reminders
	reminder1 := FromString("test task remind:2025-01-01")
	reminder2 := FromString("another task remind:2025-02-01 @home")
	reminders := []*Todo{reminder1, reminder2}

	// Write reminders
	err := WriteReminderTasks(reminderPath, reminders)
	if err != nil {
		t.Fatalf("Failed to write reminders: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(reminderPath); os.IsNotExist(err) {
		t.Fatalf("Reminder file was not created")
	}

	// Load and verify content
	loadedReminders, err := LoadReminderTasks(reminderPath)
	if err != nil {
		t.Fatalf("Failed to load written reminders: %v", err)
	}

	if len(loadedReminders) != 2 {
		t.Fatalf("Expected 2 reminders, got %d", len(loadedReminders))
	}

	if loadedReminders[0].Description != "test task" {
		t.Errorf("Expected 'test task', got '%s'", loadedReminders[0].Description)
	}
	if loadedReminders[1].Contexts[0] != "home" {
		t.Errorf("Expected context 'home', got '%s'", loadedReminders[1].Contexts[0])
	}
}

func TestGetDueReminders(t *testing.T) {
	// Create test reminders with different dates
	reminder1 := FromString("past task remind:2025-01-01")
	reminder2 := FromString("today task remind:2025-06-15")
	reminder3 := FromString("future task remind:2025-12-31")
	reminders := []*Todo{reminder1, reminder2, reminder3}

	// Test date: 2025-06-15
	testDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)

	dueReminders := GetDueReminders(reminders, testDate)

	// Should get past and today tasks (2 total)
	if len(dueReminders) != 2 {
		t.Fatalf("Expected 2 due reminders, got %d", len(dueReminders))
	}

	// Check that we got the right ones
	descriptions := make(map[string]bool)
	for _, reminder := range dueReminders {
		descriptions[reminder.Description] = true
	}

	if !descriptions["past task"] {
		t.Error("Expected 'past task' to be due")
	}
	if !descriptions["today task"] {
		t.Error("Expected 'today task' to be due")
	}
	if descriptions["future task"] {
		t.Error("Expected 'future task' to NOT be due")
	}
}

func TestRemoveDueReminders(t *testing.T) {
	reminder1 := FromString("keep me remind:2025-12-31")
	reminder2 := FromString("remove me remind:2025-01-01")
	reminder3 := FromString("also keep me remind:2025-11-30")
	allReminders := []*Todo{reminder1, reminder2, reminder3}

	dueReminders := []*Todo{reminder2}

	remaining := RemoveDueReminders(allReminders, dueReminders)

	if len(remaining) != 2 {
		t.Fatalf("Expected 2 remaining reminders, got %d", len(remaining))
	}

	// Check that the right ones remain
	descriptions := make(map[string]bool)
	for _, reminder := range remaining {
		descriptions[reminder.Description] = true
	}

	if descriptions["remove me"] {
		t.Error("Expected 'remove me' to be removed")
	}
	if !descriptions["keep me"] {
		t.Error("Expected 'keep me' to remain")
	}
	if !descriptions["also keep me"] {
		t.Error("Expected 'also keep me' to remain")
	}
}

func TestSortRemindersByDate(t *testing.T) {
	reminder1 := FromString("third remind:2025-03-01")
	reminder2 := FromString("first remind:2025-01-01")
	reminder3 := FromString("second remind:2025-02-01")
	reminders := []*Todo{reminder1, reminder2, reminder3}

	SortRemindersByDate(reminders)

	expectedOrder := []string{"first", "second", "third"}
	for i, expected := range expectedOrder {
		if reminders[i].Description != expected {
			t.Errorf("Position %d: expected '%s', got '%s'", i, expected, reminders[i].Description)
		}
	}
}

func TestProcessReminders(t *testing.T) {
	tempDir := t.TempDir()

	// Create test reminder content
	reminderContent := `past task remind:2025-01-01
today task remind:2025-06-15
future task remind:2025-12-31`

	reminderPath := ReminderTasksPath(tempDir)
	err := os.WriteFile(reminderPath, []byte(reminderContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test reminders file: %v", err)
	}

	// Create empty todo files
	todoPath := ActiveTodoPath(tempDir)
	donePath := DoneTodoPath(tempDir)
	err = os.WriteFile(todoPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create todo file: %v", err)
	}
	err = os.WriteFile(donePath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create done file: %v", err)
	}

	// Process reminders for 2025-06-15
	testDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	err = ProcessReminders(tempDir, testDate)
	if err != nil {
		t.Fatalf("Failed to process reminders: %v", err)
	}

	// Check that due reminders were moved to todos
	todos, err := LoadTodoFile(todoPath)
	if err != nil {
		t.Fatalf("Failed to load todos: %v", err)
	}

	if len(todos) != 2 {
		t.Fatalf("Expected 2 todos after processing, got %d", len(todos))
	}

	// Check that remind labels were removed and creation dates set
	for _, todo := range todos {
		if _, hasRemind := todo.Labels["remind"]; hasRemind {
			t.Error("Expected remind label to be removed from processed todo")
		}
		if todo.CreationDate.IsZero() {
			t.Error("Expected creation date to be set on processed todo")
		}
		if todo.CreationDate != testDate {
			t.Errorf("Expected creation date %v, got %v", testDate, todo.CreationDate)
		}
	}

	// Check that only future reminder remains
	remainingReminders, err := LoadReminderTasks(reminderPath)
	if err != nil {
		t.Fatalf("Failed to load remaining reminders: %v", err)
	}

	if len(remainingReminders) != 1 {
		t.Fatalf("Expected 1 remaining reminder, got %d", len(remainingReminders))
	}

	if remainingReminders[0].Description != "future task" {
		t.Errorf("Expected 'future task' to remain, got '%s'", remainingReminders[0].Description)
	}
}

func TestAddReminderTask(t *testing.T) {
	tempDir := t.TempDir()

	// Add first reminder
	reminder1 := FromString("first task remind:2025-01-01")
	err := AddReminderTask(tempDir, reminder1)
	if err != nil {
		t.Fatalf("Failed to add first reminder: %v", err)
	}

	// Add second reminder
	reminder2 := FromString("second task remind:2025-02-01")
	err = AddReminderTask(tempDir, reminder2)
	if err != nil {
		t.Fatalf("Failed to add second reminder: %v", err)
	}

	// Load and verify
	reminderPath := ReminderTasksPath(tempDir)
	reminders, err := LoadReminderTasks(reminderPath)
	if err != nil {
		t.Fatalf("Failed to load reminders: %v", err)
	}

	if len(reminders) != 2 {
		t.Fatalf("Expected 2 reminders, got %d", len(reminders))
	}

	descriptions := make(map[string]bool)
	for _, reminder := range reminders {
		descriptions[reminder.Description] = true
	}

	if !descriptions["first task"] {
		t.Error("Expected 'first task' to be present")
	}
	if !descriptions["second task"] {
		t.Error("Expected 'second task' to be present")
	}
}

func TestReminderTasksPath(t *testing.T) {
	dir := "/test/dir"
	expected := filepath.Join(dir, "reminders.txt")
	actual := ReminderTasksPath(dir)

	if actual != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, actual)
	}
}

