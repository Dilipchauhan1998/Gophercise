package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"task/db"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		newTask := strings.Join(args, " ")
		_, err := db.AddTask(db.Task{Name: newTask})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Added \"%s\" to your task list\n", newTask)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
