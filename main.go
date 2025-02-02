package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Welcome to the Quiz Game!")
	fmt.Println("What is your name?")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')

	name = strings.TrimSuffix(name, "\r\n")

	if err != nil {
		panic("Error reading name")
	}

	g.Name = name

	fmt.Printf("Hello %s, let's start the game!\n", g.Name)
}

func (g *GameState) ProcessCSV() {
	file, err := os.Open("quiz-go.csv")

	if err != nil {
		panic("Error reading file")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading csv")
	}

	for index, record := range records {
		// fmt.Println(index, record)
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}

}

func (g *GameState) Run() {
	// Display question to user
	for index, question := range g.Questions {
		fmt.Printf("\033[33m %d. %s \033[0m\n", index+1, question.Text)

		// Display options for the user
		for opt, option := range question.Options {
			fmt.Printf("[%d] %s\n", opt+1, option)
		}

		fmt.Println("Enter your answer:")
		// Get user input
		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			// Validate user input
			answer, err = toInt(strings.TrimSuffix(read, "\r\n"))
			// answer, err = toInt(read[:len(read)-2])

			// if input is wrong type again
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			break
		}

		// Check if answer is correct
		// say to user if answer is correct or not
		// if correct increment points
		if answer == question.Answer {
			fmt.Printf("\033[32m %s \033[0m\n", "Congratulations, the answer is correct!")
			g.Points += 10
		} else {
			fmt.Printf("\033[31m %s \033[0m\n", "Sorry, the answer is wrong!")
			fmt.Println("---------------------------")
		}
	}

}

func (g *GameState) End() {
	if g.Points < 20 {
		fmt.Printf("%s, you have scored less than 20 points, you are reproved!", g.Name)
	} else {
		fmt.Printf("%s, you have scored more than 20 points, you are approved!", g.Name)
	}
}

func main() {
	game := &GameState{}
	game.ProcessCSV()
	game.Init()
	game.Run()

	fmt.Printf("Congratulations %s, you have scored %d points\n", game.Name, game.Points)
	game.End()
}

func toInt(s string) (int, error) {
	int, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("character is not a number, please enter a number")
	}
	return int, err
}
