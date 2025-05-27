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
	s := fmt.Sprintf(`
Todo '%s':
    id: %s
    title: %s
    description: %s
    isDone: %t
    deadline: %v
	`, todo.Title, todo.Id, todo.Title, todo.Description, todo.IsDone, todo.Deadline)

	fmt.Println(s)
}

func formatTime(d time.Duration) string {
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

func (todo *Todo) PrintTodo() {
	spacing := strings.Repeat("\t", strings.Count(todo.Id, "-"))
	isDone := "- [ ]"
	if todo.IsDone {
		isDone = "- [x]"
	}
	id := todo.Id
	title := todo.Title
	deadline := ""
	if !todo.Deadline.IsZero() {
		timeLeft := time.Until(todo.Deadline)
		if timeLeft < 0 {
			deadline = "[OVERDUE " + formatTime(timeLeft) + "]"
		} else if timeLeft < 1*time.Hour {
			deadline = "[" + formatTime(timeLeft) + " left]"
		} else if timeLeft < 24*time.Hour {
			deadline = "[" + formatTime(timeLeft) + " left]"
		}
	}
	fmt.Println(spacing, isDone, id, title, deadline)
}

type TodoList []*Todo

func (todos TodoList) SortByDeadline() *TodoList {
	var sortedTodos = make(TodoList, len(todos))
	copy(sortedTodos, todos)
	sort.Slice(sortedTodos, func(i, j int) bool {
		t1 := sortedTodos[i]
		t2 := sortedTodos[j]
		timeUntilT1 := time.Until(t1.Deadline).Abs()
		timeUntilT2 := time.Until(t2.Deadline).Abs()
		return timeUntilT1 < timeUntilT2
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
