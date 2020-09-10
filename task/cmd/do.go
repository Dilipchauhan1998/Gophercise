package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"task/db"
	"time"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
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

        if len(ids) == 0{
        	fmt.Println("No task to mark complete")
        	return
		}

		tasks , err := db.FetchAllTask()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _ , id := range ids {

			if  id<=0 || id>len(tasks) {
				fmt.Printf("Invalid ids: %d\n",id)
				continue
			}

			task := tasks[id-1]
			err := db.RemoveTask(task.ID)
			if err != nil{
				fmt.Printf("failed to mark \"%d\" as complted, %v\n", id, err)
				os.Exit(1)
			}else {
				_, err = db.AddCompletedTask(db.CompletedTask{ID: task.ID,Name: task.Name,CompletionTime: time.Now() })
				if err != nil{
					_, err = db.AddTask(task)
				}else{
					fmt.Printf("marked \"%d\" as completed\n", id)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
