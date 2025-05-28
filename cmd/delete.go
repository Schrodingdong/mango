package cmd

import (
	"fmt"
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
		todos := utils.GetTodos()
		todoIds := args
		todos = utils.DeleteTodos(todoIds, todos)
		utils.SaveTodos(todos)
		fmt.Println("Deleted: ", todoIds)
	},
}
