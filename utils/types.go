package utils

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

var urgentThreshold time.Duration = 24 * time.Hour

type Todo struct {
	Id          string
	Title       string
	Description string
	Deadline    time.Time
	IsDone      bool
	Todos       *TodoList
}

/*
Given subtodos, Ensures parent deadline has the latest deadline of its sub todos
*/
func (todo *Todo) AdaptDeadlines() {
	if todo.Todos == nil || len(*todo.Todos) == 0 {
		return
	}
	latestDeadline := todo.Deadline
	for _, el := range *todo.Todos {
		if latestDeadline.After(el.Deadline) {
			latestDeadline = el.Deadline
		}
	}
	fmt.Println(latestDeadline)
	todo.Deadline = latestDeadline
}

func (todo *Todo) AddTodo(subtodo *Todo) {
	subTodoSlice := append(*(todo.Todos), subtodo)
	todo.Todos = &subTodoSlice

	// Update the deadline
	if (*subtodo).Deadline.Before(todo.Deadline) {
		todo.Deadline = (*subtodo).Deadline
	}
}

func (todo *Todo) ChangeStatus(status bool) {
	todo.IsDone = status
	if status {
		for i := 0; i < len(*todo.Todos); i++ {
			t := (*todo.Todos)[i]
			t.ChangeStatus(status)
		}
	}
}

func (todo *Todo) PrintTodoDetail() {
	// Prefix
	spaceCount := strings.Count(todo.Id, "-")
	indentSpaceCount := 4
	prefix := strings.Repeat(" ", spaceCount*indentSpaceCount)

	// Format status
	status := "[ ]"
	if !todo.IsDone {
		status = "[X]"
	}

	// Format deadline
	deadlineStr := "No deadline"
	if !todo.Deadline.IsZero() {
		deadlineStr = todo.Deadline.Local().String()
	}

	// Build the output with consistent indentation
	output := fmt.Sprintf(`
%s- %s Todo: %s
%s  ├── ID:          %s
%s  ├── Description: %s
%s  └── Deadline:    %s`,
		prefix, status, todo.Title,
		prefix, todo.Id,
		prefix, todo.Description,
		prefix, deadlineStr,
	)

	fmt.Println(output)
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Abs().Minutes())
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dd", hours/24)
}

func (todo *Todo) PrintTodoOneLine() {
	isDone := "- [ ]" // Red for incomplete
	deadline := ""
	spacing := strings.Repeat("\t", strings.Count(todo.Id, "-"))

	id := todo.Id
	title := todo.Title
	colorPrefix := ""
	if todo.IsDone {
		isDone = "- [✓]"
		colorPrefix = "\033[32m" // Green for complete
	} else {
		if !todo.Deadline.IsZero() {
			timeLeft := time.Until(todo.Deadline)
			if timeLeft < 0 {
				deadline = "[OVERDUE " + formatDuration(timeLeft) + "]\033[0m"
				colorPrefix = "\033[31m" // Red for overdue
			} else if timeLeft < 24*time.Hour {
				deadline = "[" + formatDuration(timeLeft) + " left]\033[0m"
				colorPrefix = "\033[33m" // yellow for close deadline
			}
		}
	}
	formatString := colorPrefix + "%s %s %3s %-32s %s \033[0m\n"
	fmt.Printf(
		formatString,
		spacing, isDone, id, title, deadline,
	)
}

type TodoList []*Todo

func (todos *TodoList) SortByDeadline() *TodoList {
	var sortedTodos = make(TodoList, len(*todos))
	copy(sortedTodos, *todos)
	sort.Slice(sortedTodos, func(i, j int) bool {
		t1 := sortedTodos[i]
		t2 := sortedTodos[j]

		// Todos without deadlines
		if t1.Deadline.IsZero() {
			return false
		}
		if t2.Deadline.IsZero() {
			return true
		}

		// Compare deadlines directly (oldest first)
		return t1.Deadline.Before(t2.Deadline)
	})
	return &sortedTodos
}

func (todos TodoList) FilterTodosDone() *TodoList {
	filteredTodos := TodoList{}
	for i := 0; i < len(todos); i++ {
		todo := (todos)[i]
		if todo.IsDone {
			filteredTodos = append(filteredTodos, todo)
		}
	}
	return &filteredTodos
}

/*
Filter urgent todos, depending on the time left to deadline. Currently an urgent todo has a threshold of 1 day
*/
func (todos TodoList) FilterTodosUrgent() *TodoList {
	filteredTodos := TodoList{}
	for i := 0; i < len(todos); i++ {
		todo := (todos)[i]
		if todo.Deadline.IsZero() {
			continue
		}
		if time.Until(todo.Deadline) < urgentThreshold {
			filteredTodos = append(filteredTodos, todo)
		}
	}
	return &filteredTodos
}
