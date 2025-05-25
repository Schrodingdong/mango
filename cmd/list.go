package cmd

import (
	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

var isUrgent bool

// var isDone bool

func init() {
	listCmd.PersistentFlags().BoolVar(&isUrgent, "urgent", false, "List urgent todos")
	listCmd.PersistentFlags().BoolVar(&isDone, "done", false, "List done todos")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all the todos",
	Long:  "Print all the todos",
	Run: func(cmd *cobra.Command, args []string) {
		todos := utils.GetTodos()
		todos = todos.SortByDeadline()
		printTodos(todos)
	},
}

func printTodos(todos *utils.TodoList) {
	if isUrgent {
		todos = todos.FilterTodosUrgent()
	}
	if isDone {
		todos = todos.FilterTodosDone()
	}

	for i := 0; i < len(*todos); i++ {
		todo := (*todos)[i]
		todo.PrintTodo()
		if len(*todo.Todos) != 0 {
			printTodos(todo.Todos)
		}
	}
}
