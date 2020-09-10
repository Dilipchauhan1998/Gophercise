package main

import (
	"testing"
)

type TestFile struct{
	filePath string
	quizProblems []problem
	error bool
}

func TestReadProblemsFromCSVFile(t *testing.T){

	testFileData := []TestFile{
		TestFile{
			filePath: "",
			quizProblems: []problem{},
			error: true,
		},
		TestFile{
			filePath: "testdata/problems.txt",
			quizProblems: []problem{},
			error: true,
		},
		TestFile{
			filePath: "testdata/problems_1.txt",
			quizProblems: []problem{},
			error: true,
		},
		TestFile{
			filePath: "testdata/problems.csv",
			quizProblems: []problem{},
			error: false,
		},
		TestFile{
			filePath: "testdata/problems_1.csv",
			quizProblems: []problem{
				problem{
					question: "5+5",
					answer: "10",
				},
				problem{
					question: "7+3",
					answer: "10",
				},
				problem{
					question: "8+3",
					answer: "11",
				},
			},
			error: false,
		},
	}

   for _, value := range testFileData{

   		quizProblems , err := readProblemsFromCSVFile(value.filePath)

   		if value.error == true {
   			if err == nil {
   				t.Error("Expected error")
			}
		} else{

			if len(quizProblems) == len(value.quizProblems) {

				for i, _ := range quizProblems{
					if quizProblems[i] != value.quizProblems[i]{
						t.Error("Question doesn't match ")
					}
				}

			}else{
				t.Error("Content don't match")
			}

		}
   }

}




