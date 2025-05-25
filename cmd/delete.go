package cmd

import (
	"log"

	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the todo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Command require a positional argument <todo_id>")
		}
		todoId := args[0]
		todos := utils.GetTodos()
		todos = utils.DeleteTodo(todoId, todos)
		utils.SaveTodos(todos)
	},
}
