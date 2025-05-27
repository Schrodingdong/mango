package cmd

import (
	"log"
	"time"

	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

var description string
var deadlineString string
var parentId string

func init() {
	// Flags
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
		todos := utils.GetTodos()
		// Get title from args
		if len(args) == 0 {
			log.Fatal("Not enough args")
		}
		title := args[0]

		// Add deadline
		todoDeadline := parseDeadline(deadlineString)
		if deadlineString != "" {
			if time.Now().After(todoDeadline) {
				log.Fatal("Can't have a deadline before now")
			}
		}

		// Init todo data
		todo := utils.Todo{
			Id:          utils.AssignId(parentId, todos),
			Title:       title,
			Description: description,
			Deadline:    todoDeadline,
			IsDone:      isDone,
			Todos:       &utils.TodoList{},
		}

		// Add todo
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

		// Print added todo details
		todo.PrintTodoDetail()
	},
}

/*
Parse a dedadline string into a `time.Time` object. if an empty string is provided, we get t=0ns.
*/
func parseDeadline(deadlineString string) time.Time {
	deadline, err := time.ParseInLocation(time.DateTime, deadlineString, time.Local)
	if err != nil {
		return zero
	}
	return deadline
}
