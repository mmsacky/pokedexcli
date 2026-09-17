package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		userInput := scanner.Text()

		words := cleanInput(userInput)

		supportedCommands := getCommands()

		switch words[0] {
		case supportedCommands["exit"].name:
			supportedCommands["exit"].callback()
		case supportedCommands["help"].name:
			supportedCommands["help"].callback()
		default:
			fmt.Println("Unknown command.")
		}

	}
}
