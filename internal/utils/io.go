package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func GetPlayerInput(mess string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(mess)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("something went wrong, try again. cause: %w", err)
	}
	input = strings.TrimSpace(input)
	return input, nil
}
