package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {

	csvFilePath, timeLimit, shuffle := setFlag()
	quizProblems, err := readProblemsFromCSVFile(csvFilePath)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if len(quizProblems) == 0{
		fmt.Println("No problems found in the file")
		os.Exit(1)
	}

	if shuffle == true {
		shuffleProblems(quizProblems)
	}

	gameOver := make(chan bool)
	correctAnswerCount := 0

	fmt.Print("Start Game !")
	fmt.Scanf("%s")

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	go playGame(quizProblems, &correctAnswerCount, gameOver)

outer:
	for {
		select {
		case <-timer.C:
			fmt.Printf("\nYour Scored %d out of  %d\n", correctAnswerCount, len(quizProblems))
			break outer
		case <-gameOver:
			fmt.Printf("Your Scored %d out of  %d\n", correctAnswerCount, len(quizProblems))
			break outer

		}
	}
}

func setFlag() (string, int, bool) {
	fPath := flag.String("fpath", "problems.csv", "path to the quiz probelms csv file in format 'question, ans'")
	limit := flag.Int("limit", 30, "game time limit in seconds")
	shuffle := flag.Bool("shuffle", false, "True if want to shuffle the questions order false otherwise")
	flag.Parse()

	return *fPath, *limit, *shuffle
}

func readProblemsFromCSVFile(csvFilePath string) ([]problem, error) {

	if filepath.Ext(csvFilePath) !=".csv" {
		return []problem{}, errors.New("Not a CSV file")
	}

	quizProblems := []problem{}
	//Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		return []problem{}, err
	}

	// Create a new reader.
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		if err != nil {
			return []problem{}, err
		}

		if len(record) == 2 {
			quizProblems = append(quizProblems, problem{question: record[0], answer: strings.TrimSpace(record[1])})
		}
	}
	return quizProblems, nil
}

func shuffleProblems(quizProblems []problem) {
	rand.Seed(time.Now().UnixNano())

	for i := range quizProblems {
		randIndex := rand.Intn(len(quizProblems))
		quizProblems[i], quizProblems[randIndex] = quizProblems[randIndex], quizProblems[i]
	}
}

func playGame(quizProblems []problem, correctAnswerCount *int, gameOver chan bool) {

	for i, value := range quizProblems {
		question := value.question
		correctAnswer := value.answer

		var userAnswer string
		fmt.Printf("Problem #%d: %s ", i+1, question)
		fmt.Scanf("%s", &userAnswer)

		if strings.TrimSpace(userAnswer) == correctAnswer {
			*correctAnswerCount++
		}
	}

	gameOver <- true
}
