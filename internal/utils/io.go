package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func GetConditionalPlayerInput(expectedAnswers []string, msg string) string {
	var playerInput string

	for !contains(expectedAnswers, playerInput) {
		playerInput, _ = GetPlayerInput(msg)
	}
	return playerInput
}

func GetPlayerInput(msg string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("something went wrong, try again. cause: %w", err)
	}
	input = strings.TrimSpace(input)
	return input, nil
}
