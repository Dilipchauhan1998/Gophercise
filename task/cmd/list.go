package cmd

import (
	"fmt"
	"os"
	"task/db"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: " List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.FetchAllTask()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("No task to complete!")
		}

		for i, v := range tasks {
			fmt.Printf("%d  %s\n", i+1, v.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
