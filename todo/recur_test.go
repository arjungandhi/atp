package todo

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestRecurringTaskFromString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantErr      bool
		expectNil    bool
		validateFunc func(*testing.T, *RecurringTask)
	}{
		{
			name:    "daily task with project and context",
			input:   "@daily Check email +work @office",
			wantErr: false,
			validateFunc: func(t *testing.T, rt *RecurringTask) {
				if rt.ScheduleStr != "@daily" {
					t.Errorf("Expected schedule '@daily', got %s", rt.ScheduleStr)
				}
				if rt.Todo.Description != "Check email" {
					t.Errorf("Expected description 'Check email', got %s", rt.Todo.Description)
				}
			},
		},
		{
			name:    "weekly task with priority",
			input:   "@weekly (A) Review weekly goals +personal",
			wantErr: false,
			validateFunc: func(t *testing.T, rt *RecurringTask) {
				if rt.ScheduleStr != "@weekly" {
					t.Errorf("Expected schedule '@weekly', got %s", rt.ScheduleStr)
				}
				if rt.Todo.Priority != "A" {
					t.Errorf("Expected priority 'A', got %s", rt.Todo.Priority)
				}
			},
		},
		{
			name:    "cron task with labels",
			input:   "0 9 * * 1 Team standup @office +work priority:B",
			wantErr: false,
			validateFunc: func(t *testing.T, rt *RecurringTask) {
				if rt.ScheduleStr != "0 9 * * 1" {
					t.Errorf("Expected schedule '0 9 * * 1', got %s", rt.ScheduleStr)
				}
				if rt.Todo.Description != "Team standup" {
					t.Errorf("Expected description 'Team standup', got %s", rt.Todo.Description)
				}
			},
		},
		{
			name:    "simple daily task",
			input:   "@daily Water plants",
			wantErr: false,
			validateFunc: func(t *testing.T, rt *RecurringTask) {
				if rt.ScheduleStr != "@daily" {
					t.Errorf("Expected schedule '@daily', got %s", rt.ScheduleStr)
				}
				if rt.Todo.Description != "Water plants" {
					t.Errorf("Expected description 'Water plants', got %s", rt.Todo.Description)
				}
			},
		},
		{
			name:     "empty line",
			input:    "",
			wantErr:  false,
			expectNil: true,
		},
		{
			name:     "comment line",
			input:    "# This is a comment",
			wantErr:  false,
			expectNil: true,
		},
		{
			name:    "invalid simple format - missing task",
			input:   "@daily",
			wantErr: true,
		},
		{
			name:    "invalid simple format - with colon",
			input:   "@daily: Water plants",
			wantErr: true,
		},
		{
			name:    "invalid simple schedule",
			input:   "@hourly Water plants",
			wantErr: true,
		},
		{
			name:    "invalid cron format - too few fields",
			input:   "0 9 * Team standup",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RecurringTaskFromString(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("Expected result, got nil")
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

func TestRecurringTaskString(t *testing.T) {
	task, err := NewRecurringTask("@daily", &Todo{
		Priority:    "A",
		Description: "Check email",
		Projects:    []string{"work"},
		Contexts:    []string{"office"},
		Labels:      map[string]string{"priority": "B"},
	})
	if err != nil {
		t.Fatalf("Failed to create RecurringTask: %v", err)
	}

	result := task.String()
	expected := "@daily (A) Check email +work @office priority:B"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestShouldGenerateForDate(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		date     time.Time
		expected bool
	}{
		{
			name:     "daily task should always generate",
			schedule: "@daily",
			date:     time.Date(2025, 6, 14, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "weekly task on Monday",
			schedule: "@weekly",
			date:     time.Date(2025, 6, 16, 9, 0, 0, 0, time.UTC), // Monday
			expected: true,
		},
		{
			name:     "weekly task on Tuesday",
			schedule: "@weekly",
			date:     time.Date(2025, 6, 17, 9, 0, 0, 0, time.UTC), // Tuesday
			expected: false,
		},
		{
			name:     "monthly task on 1st",
			schedule: "@monthly",
			date:     time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "monthly task on 15th",
			schedule: "@monthly",
			date:     time.Date(2025, 6, 15, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "cron every day at midnight",
			schedule: "0 0 * * *",
			date:     time.Date(2025, 6, 14, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "cron Monday at midnight",
			schedule: "0 0 * * 1",
			date:     time.Date(2025, 6, 16, 9, 0, 0, 0, time.UTC), // Monday
			expected: true,
		},
		{
			name:     "cron Monday at midnight (wrong day)",
			schedule: "0 0 * * 1",
			date:     time.Date(2025, 6, 17, 9, 0, 0, 0, time.UTC), // Tuesday
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewRecurringTask(tt.schedule, NewTodo())
			if err != nil {
				t.Fatalf("Failed to create RecurringTask: %v", err)
			}
			result := task.ShouldGenerateForDate(tt.date)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

// Remove TestMatchesCronField since we no longer have custom cron field matching

func TestGenerateTodo(t *testing.T) {
	task, err := NewRecurringTask("@daily", &Todo{
		Priority:    "A",
		Description: "Check email",
		Projects:    []string{"work"},
		Contexts:    []string{"office"},
		Labels:      map[string]string{"priority": "B"},
	})
	if err != nil {
		t.Fatalf("Failed to create RecurringTask: %v", err)
	}

	date := time.Date(2025, 6, 14, 9, 0, 0, 0, time.UTC)
	todo := task.GenerateTodo(date)

	if todo.Done {
		t.Error("Generated todo should not be done")
	}

	if !todo.CreationDate.IsZero() {
		t.Errorf("Expected no creation date, got %v", todo.CreationDate)
	}

	// Should have recur label with the date
	if recur, exists := todo.Labels["recur"]; !exists || recur != date.Format("2006-01-02") {
		t.Errorf("Expected recur label with date %s, got %s (exists: %v)", date.Format("2006-01-02"), recur, exists)
	}

	if todo.Priority != "A" {
		t.Errorf("Expected priority A, got %s", todo.Priority)
	}

	if todo.Description != "Check email" {
		t.Errorf("Expected description 'Check email', got %s", todo.Description)
	}

	if !reflect.DeepEqual(todo.Projects, []string{"work"}) {
		t.Errorf("Expected projects [work], got %v", todo.Projects)
	}

	if !reflect.DeepEqual(todo.Contexts, []string{"office"}) {
		t.Errorf("Expected contexts [office], got %v", todo.Contexts)
	}

	expectedLabels := map[string]string{
		"priority": "B",
		"recur":    "2025-06-14",
	}
	if !reflect.DeepEqual(todo.Labels, expectedLabels) {
		t.Errorf("Expected labels %v, got %v", expectedLabels, todo.Labels)
	}
}

func TestLoadWriteRecurringTasks(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "atp_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	recurPath := filepath.Join(tmpDir, "recur.txt")

	// Create test tasks
	task1, err := NewRecurringTask("@daily", &Todo{
		Description: "Check email",
		Projects:    []string{"work"},
		Contexts:    []string{"office"},
		Labels:      map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	task2, err := NewRecurringTask("@weekly", &Todo{
		Priority:    "A",
		Description: "Weekly review",
		Projects:    []string{"personal"},
		Contexts:    []string{},
		Labels:      map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	tasks := []*RecurringTask{task1, task2}

	// Write tasks
	err = WriteRecurringTasks(recurPath, tasks)
	if err != nil {
		t.Fatalf("Failed to write recurring tasks: %v", err)
	}

	// Load tasks
	loadedTasks, err := LoadRecurringTasks(recurPath)
	if err != nil {
		t.Fatalf("Failed to load recurring tasks: %v", err)
	}

	if len(loadedTasks) != len(tasks) {
		t.Errorf("Expected %d tasks, got %d", len(tasks), len(loadedTasks))
	}

	for i, task := range loadedTasks {
		if task.ScheduleStr != tasks[i].ScheduleStr {
			t.Errorf("Task %d: expected schedule %s, got %s", i, tasks[i].ScheduleStr, task.ScheduleStr)
		}
		if task.Todo.Description != tasks[i].Todo.Description {
			t.Errorf("Task %d: expected description %s, got %s", i, tasks[i].Todo.Description, task.Todo.Description)
		}
	}
}

func TestGenerateTodosFromRecurring(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "atp_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	recurPath := filepath.Join(tmpDir, "recur.txt")

	// Create test recurring tasks
	task1, err := NewRecurringTask("@daily", &Todo{
		Description: "Daily task",
		Projects:    []string{"work"},
		Labels:      map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	task2, err := NewRecurringTask("@weekly", &Todo{
		Description: "Weekly task",
		Projects:    []string{"personal"},
		Labels:      map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	tasks := []*RecurringTask{task1, task2}

	err = WriteRecurringTasks(recurPath, tasks)
	if err != nil {
		t.Fatal(err)
	}

	// Test on Monday (should generate both daily and weekly)
	monday := time.Date(2025, 6, 16, 9, 0, 0, 0, time.UTC)
	todos, err := GenerateTodosFromRecurring(tmpDir, monday)
	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos on Monday, got %d", len(todos))
	}

	// Test on Tuesday (should generate only daily)
	tuesday := time.Date(2025, 6, 17, 9, 0, 0, 0, time.UTC)
	todos, err = GenerateTodosFromRecurring(tmpDir, tuesday)
	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo on Tuesday, got %d", len(todos))
	}

	if todos[0].Description != "Daily task" {
		t.Errorf("Expected 'Daily task', got %s", todos[0].Description)
	}
}

func TestTodoExistsForRecurringTask(t *testing.T) {
	todos := []*Todo{
		{
			Description: "Daily task",
			Labels:      map[string]string{"recur": "2025-06-14"},
		},
		{
			Description: "Other task",
			Labels:      map[string]string{},
		},
	}

	// Should find existing todo
	if !todoExistsForRecurringTask(todos, "Daily task", "2025-06-14") {
		t.Error("Should find existing recurring todo")
	}

	// Should not find non-existing todo
	if todoExistsForRecurringTask(todos, "Daily task", "2025-06-15") {
		t.Error("Should not find todo for different date")
	}

	// Should not find different task
	if todoExistsForRecurringTask(todos, "Different task", "2025-06-14") {
		t.Error("Should not find different task")
	}
}

func TestAddRecurringTodosToDir(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "atp_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create recurring tasks
	recurPath := filepath.Join(tmpDir, "recur.txt")
	task, err := NewRecurringTask("@daily", &Todo{
		Description: "Daily task",
		Projects:    []string{"work"},
		Labels:      map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	tasks := []*RecurringTask{task}

	err = WriteRecurringTasks(recurPath, tasks)
	if err != nil {
		t.Fatal(err)
	}

	// Add recurring todos
	date := time.Date(2025, 6, 14, 9, 0, 0, 0, time.UTC)
	err = AddRecurringTodosToDir(tmpDir, date)
	if err != nil {
		t.Fatal(err)
	}

	// Load and verify todos were added
	todos, err := LoadTodoDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}

	if todos[0].Description != "Daily task" {
		t.Errorf("Expected 'Daily task', got %s", todos[0].Description)
	}

	if todos[0].Labels["recur"] != "2025-06-14" {
		t.Errorf("Expected recur label '2025-06-14', got %s", todos[0].Labels["recur"])
	}

	// Try adding again - should not duplicate
	err = AddRecurringTodosToDir(tmpDir, date)
	if err != nil {
		t.Fatal(err)
	}

	todos, err = LoadTodoDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo after duplicate prevention, got %d", len(todos))
	}
}

func TestCronValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{
			name:    "valid cron expression",
			input:   "0 9 * * 1 Weekly meeting",
			wantErr: false,
		},
		{
			name:    "invalid minute range",
			input:   "60 9 * * 1 Invalid minute",
			wantErr: true,
		},
		{
			name:    "invalid hour range", 
			input:   "0 25 * * 1 Invalid hour",
			wantErr: true,
		},
		{
			name:    "invalid cron format",
			input:   "0-25-30 9 * * 1 Invalid range",
			wantErr: true,
		},
		{
			name:    "valid range",
			input:   "0-30 9-17 1-15 * 1-5 Work hours",
			wantErr: false,
		},
		{
			name:    "valid step",
			input:   "*/15 * * * * Every 15 minutes",
			wantErr: false,
		},
		{
			name:    "valid comma separated",
			input:   "0 9,12,18 * * 1,3,5 Multiple times",
			wantErr: false,
		},
		{
			name:    "invalid simple schedule",
			input:   "@hourly Water plants",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RecurringTaskFromString(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}