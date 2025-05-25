package cmd

import (
	"log"

	"github.com/schrodingdong/mango/utils"
	"github.com/spf13/cobra"
)

var sike bool

func init() {
	doneCmd.PersistentFlags().BoolVar(&sike, "sike", false, "Un-done-ify the todo")
	rootCmd.AddCommand(doneCmd)
}

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark a todo as 'done'",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Command require a positional argument <todo_id>")
		}
		todoId := args[0]
		utils.ChangeTodoState(!sike, todoId)
	},
}
