package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Quiz struct {
	Question string
	Answer   string
}

func loadCsv(filename string) ([]Quiz, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, errReadLine := reader.ReadAll()
	if errReadLine != nil {
		return nil, errReadLine
	}

	return parseLines(lines), nil
}

func parseLines(lines [][]string) []Quiz {
	quiz := make([]Quiz, len(lines))
	for index, line := range lines {
		if len(line) != 2 {
			continue
		}
		quiz[index] = Quiz{
			Question: line[0],
			Answer:   line[1],
		}
	}

	return quiz
}

func main() {
	fileName := flag.String("csv", "quiz.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	correctAnswers := 0

	quiz, errLoadQuix := loadCsv(*fileName)
	if errLoadQuix != nil {
		exit(fmt.Sprintf("Failed to load the quiz: %s\n", errLoadQuix))
	}

	for index, q := range quiz {
		fmt.Printf("Problem #%d: %s = ", index+1, q.Question)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correctAnswers, len(quiz))
			return
		case answer := <-answerChannel:
			if q.Answer == answer {
				correctAnswers++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correctAnswers, len(quiz))
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}
