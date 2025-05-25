package cmd

import (
	"log"
	"time"

	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

var title string
var description string
var deadlineString string
var isDone bool
var parentId string

func init() {
	// Flags
	createCmd.PersistentFlags().StringVar(&title, "title", "", "Todo title")
	createCmd.PersistentFlags().StringVar(&description, "description", "", "Todo description")
	createCmd.PersistentFlags().StringVar(&deadlineString, "deadline", "", "Absolute deadline (DateTime format, e.g., '2025-05-24 21:41:23')")
	createCmd.PersistentFlags().StringVar(&parentId, "parent", "", "Parent todo id")
	createCmd.PersistentFlags().BoolVar(&isDone, "done", false, "Todo status")
	rootCmd.AddCommand(createCmd)
}

var zero time.Time

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new todo",
	Run: func(cmd *cobra.Command, args []string) {
		todoDeadline := parseDeadline(deadlineString)
		if time.Now().After(todoDeadline) {
			log.Fatal("Can't have a deadline before now")
		}
		todo := utils.Todo{
			Title:       title,
			Description: description,
			Deadline:    todoDeadline,
			IsDone:      isDone,
			Todos:       &utils.TodoList{},
		}
		todo.Id = utils.AssignId(parentId)

		todos := utils.GetTodos()
		if len(parentId) == 0 {
			todos = utils.AddTodo(&todo, todos)
		} else {
			parentTodo, err := utils.GetTodo(parentId, todos)
			if err != nil {
				log.Fatal(err)
			}
			parentTodo.Todos = utils.AddTodo(&todo, parentTodo.Todos)
		}
		utils.SaveTodos(todos)
		todo.PrintTodoDetail()
	},
}

func parseDeadline(deadlineString string) time.Time {
	deadline, err := time.ParseInLocation(time.DateTime, deadlineString, time.Local)
	if err != nil {
		return zero
	}
	return deadline
}
