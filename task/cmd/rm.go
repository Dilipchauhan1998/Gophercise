package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"task/db"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove a task from todo list",
	Run: func(cmd *cobra.Command, args []string) {
		ids := []int{}

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("\"%s\" is not a valid id.\n", arg)
			} else {
				ids = append(ids, id)
			}
		}

		if len(ids) == 0 {
			fmt.Println("No task to remove")
			return
		}

		tasks, err := db.FetchAllTask()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, id := range ids {

			if id <= 0 || id > len(tasks) {
				fmt.Printf("Invalid ids: %d\n", id)
				continue
			}

			task := tasks[id-1]
			err := db.RemoveTask(task.ID)
			if err != nil {
				fmt.Printf("failed to delete \"%s\" , %v\n", task.Name, err)
				os.Exit(1)
			} else {
				fmt.Printf("You have deleted the \"%s\" task \n", task.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
