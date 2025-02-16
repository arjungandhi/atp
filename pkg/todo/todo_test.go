package todo

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		line         string
		expectedTodo Todo
	}{
		{
			// Test with no priority or dates
			line: "Call Mom @phone +Family",
			expectedTodo: Todo{
				Done:        false,
				Priority:    "",
				Description: "Call Mom",
				Projects:    []string{"Family"},
				Contexts:    []string{"phone"},
				Labels:      map[string]string{},
			},
		},
		{
			// Test with priority and creation date
			line: "(A) 2025-02-15 Call Mom +Family @phone",
			expectedTodo: Todo{
				Done:         false,
				Priority:     "A",
				CreationDate: parseDate("2025-02-15"),
				Description:  "Call Mom",
				Projects:     []string{"Family"},
				Contexts:     []string{"phone"},
				Labels:       map[string]string{},
			},
		},
		{
			// Test with completion date
			line: "x 2025-02-16 2025-02-15 Call Mom +Family @phone",
			expectedTodo: Todo{
				Done:           true,
				CreationDate:   parseDate("2025-02-15"),
				CompletionDate: parseDate("2025-02-16"),
				Description:    "Call Mom",
				Projects:       []string{"Family"},
				Contexts:       []string{"phone"},
				Labels:         map[string]string{},
			},
		},
		{
			// Test with metadata
			line: "x 2025-02-15 2025-02-14 Call Dad +Home @phone due:2025-02-20",
			expectedTodo: Todo{
				Done:           true,
				CreationDate:   parseDate("2025-02-14"),
				CompletionDate: parseDate("2025-02-15"),
				Priority:       "",
				Description:    "Call Dad",
				Projects:       []string{"Home"},
				Contexts:       []string{"phone"},
				Labels:         map[string]string{"due": "2025-02-20"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			todo := FromString(test.line)
			if !reflect.DeepEqual(*todo, test.expectedTodo) {
				t.Errorf("FromString() failed for %s: got %v, want %v", test.line, *todo, test.expectedTodo)
			}
		})
	}
}

func TestToString(t *testing.T) {
	todo := &Todo{
		Done:           true,
		CreationDate:   parseDate("2025-02-15"),
		CompletionDate: parseDate("2025-02-16"),
		Priority:       "A",
		Description:    "Call Mom",
		Projects:       []string{"Family"},
		Contexts:       []string{"phone"},
		Labels:         map[string]string{"due": "2025-02-20"},
	}

	expectedStr := "x 2025-02-16 (A) 2025-02-15 Call Mom +Family @phone due:2025-02-20"
	actualStr := todo.String()

	if expectedStr != actualStr {
		t.Errorf("String() failed: got %s, want %s", actualStr, expectedStr)
	}
}

func TestLoadAndSave(t *testing.T) {
	// Prepare test data
	todos := []*Todo{
		{
			Done:        false,
			Priority:    "A",
			Description: "Call Mom",
			Projects:    []string{"Family"},
			Contexts:    []string{"phone"},
			Labels:      map[string]string{},
		},
		{
			Done:           true,
			CreationDate:   parseDate("2025-02-15"),
			CompletionDate: parseDate("2025-02-16"),
			Priority:       "B",
			Description:    "Call Dad",
			Projects:       []string{"Home"},
			Contexts:       []string{"phone"},
			Labels:         map[string]string{"due": "2025-02-20"},
		},
	}
	// create a temp file
	file, err := os.CreateTemp("", "txt")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	// Save to file
	err = WriteTodoFile(file.Name(), todos)
	if err != nil {
		t.Fatalf("Error writing file: %v", err)
	}

	// Load from file
	_, err = LoadTodoFile(file.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

}

// Helper function to parse dates
func parseDate(dateStr string) time.Time {
	date, _ := time.Parse("2006-01-02", dateStr)
	return date
}
