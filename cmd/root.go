package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Common flags
var isDone bool

var rootCmd = &cobra.Command{
	Use:   "mango",
	Short: "Todo manager for busys terminal users",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
Mango - The todo manager for busy terminal users
usage: mango create "My todo"
		`)

		fmt.Println("To list the todos: mango list")
		fmt.Println("To Create a todo: mango create")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
