package main

import (
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {

	supportedCommands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	return supportedCommands
}

func cleanInput(text string) []string {
	var words []string

	loweredText := strings.ToLower(text)
	cleanedText := strings.TrimSpace(loweredText)
	words = strings.Split(cleanedText, " ")

	return words
}

func commandExit() error {
	fmt.Println()
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, value := range getCommands() {

		fmt.Printf("%s: %s\n", value.name, value.description)

	}
	fmt.Println()
	return nil
}
