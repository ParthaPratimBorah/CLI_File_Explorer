package cli

import (
	"fmt"
	"strings"
)

// askForConfirmation asks the user to enter y or n.
func askForConfirmation(message string) bool {
	var answer string

	fmt.Printf("%s [y/n]: ", message)

	_, err := fmt.Scanln(&answer)

	if err != nil {
		return false
	}

	answer = strings.ToLower(answer)

	return answer == "y" || answer == "yes"
}