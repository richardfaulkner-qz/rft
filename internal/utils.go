package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetUserChatInput() (string, error) {
	s := "\n\nUser: "
	return GetUserInput(s)
}

func GetUserInput(s string) (string, error) {
	// Print the passed in string
	fmt.Println(s)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
		return "", scanner.Err()
	}
	return scanner.Text(), nil
}

func GetYNUserInput() (bool, error) {
	s := "continue? (Y/n)"

	response, err := GetUserInput(s)
	if err != nil {
		return false, err
	}

	if response == "" || strings.ToUpper(response) == "Y" || strings.ToUpper(response) == "N" {
		return strings.ToUpper(response) != "N", nil
	}

	return false, nil
}
