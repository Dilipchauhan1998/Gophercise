package main

import (
	"fmt"
	"os"
	"task/cmd"
	"task/db"
)

func main() {
	home, err := os.UserHomeDir()
	if err !=nil{
		handleError(fmt.Errorf("could not found home dir, %v",err))
	}

	err = db.SetupDB(home+"/my.db")
	if err != nil {
		handleError(err)
	}

	cmd.Execute()
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
