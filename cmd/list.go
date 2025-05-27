package cmd

import (
	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

var isUrgent bool
var number int

func init() {
	listCmd.PersistentFlags().BoolVar(&isUrgent, "urgent", false, "List urgent todos")
	listCmd.PersistentFlags().BoolVar(&isDone, "done", false, "List done todos")
	listCmd.PersistentFlags().IntVarP(&number, "number", "n", 0, "number of displayed todos")
	listCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all the todos",
	Long:  "Print all the todos",
	Run: func(cmd *cobra.Command, args []string) {
		todos := utils.GetTodos()
		if isUrgent {
			todos = todos.FilterTodosUrgent()
		}
		if isDone {
			todos = todos.FilterTodosDone()
		}
		if number > 0 {
			slicedTodos := (*todos)[:number]
			todos = &slicedTodos
		}
		printTodos(todos)
	},
}

func printTodos(todos *utils.TodoList) {
	for i := 0; i < len(*todos); i++ {
		todo := (*todos)[i]
		if !verbose {
			todo.PrintTodoOneLine()
		} else {
			todo.PrintTodoDetail()
		}
		if len(*todo.Todos) != 0 {
			printTodos(todo.Todos)
		}
	}
}
