package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"task/db"
	"time"
)

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your completed tasks",

	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.FetchAllCompletedTask()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("No completed task found!")
			return
		}

        fmt.Println("you completed following tasks today:")
		for _, task := range tasks{
			if time.Now().Second() - task.CompletionTime.Second() <= 86400 {

				fmt.Printf("-%s\n",task.Name)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
