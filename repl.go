package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/mmsacky/pokedexcli/internal/pokeapi"
)

type config struct {
	commandsRegistry map[string]cliCommand
	nextUrl          string
	previousUrl      string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
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
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Displays the names of the previous 20 location areas in the Pokemon world.",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of pokemon found in that location.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Used to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inpect",
			description: "Used to get details of Pokemon you've caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Used to list all the Pokemon you've caught",
			callback:    commandPokedex,
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

func commandExit(config *config, blank ...string) error {
	fmt.Println()
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandHelp(config *config, blank ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range config.commandsRegistry {

		fmt.Printf("%s: %s\n", command.name, command.description)

	}
	fmt.Println()
	return nil
}

func commandMap(config *config, blank ...string) error {

	locationsAreas, _ := pokeapi.GetLocationAreas(config.nextUrl)

	pokeapi.PrintLocations(locationsAreas)

	config.nextUrl = locationsAreas.Next

	if previous, ok := locationsAreas.Previous.(string); ok {

		config.previousUrl = previous
	}

	return nil
}

func commandMapBack(config *config, blank ...string) error {

	if config.previousUrl == "" {
		fmt.Println()
		fmt.Println("you're on the first page")
		fmt.Println()
		return nil
	}

	locationsAreas, _ := pokeapi.GetLocationAreas(config.previousUrl)

	if previous, ok := locationsAreas.Previous.(string); ok {
		config.previousUrl = previous
		config.nextUrl = locationsAreas.Next
	}

	pokeapi.PrintLocations(locationsAreas)

	if locationsAreas.Previous == nil {
		config.previousUrl = ""
		config.nextUrl = locationsAreas.Next
	}

	return nil
}

func commandExplore(config *config, locationName ...string) error {

	if len(locationName) == 0 || len(locationName) > 1 {
		fmt.Println("Please enter a valid location name")
		return nil
	}

	_, location := pokeapi.GetLocationAreas(pokeapi.LocationAreasURL + locationName[0])

	if location.ID != 0 {
		pokeapi.PrintPokemonNames(location)
	} else {
		fmt.Println("That location doesn't exist")
	}

	return nil
}

func commandCatch(config *config, pokemonName ...string) error {
	if len(pokemonName) == 0 || len(pokemonName) > 1 {
		fmt.Println("Please enter a valid pokemon name")
		return nil
	}

	pokemon := pokeapi.GetPokemon(pokeapi.PokemonURL + pokemonName[0])

	if pokemon.ID != 0 {
		difficulty := 0
		chanceToCatch := rand.Intn(pokemon.BaseExperience)

		switch {
		case pokemon.BaseExperience > 100 && pokemon.BaseExperience < 300:
			difficulty = 30
		case pokemon.BaseExperience > 300:
			difficulty = 50
		default:
			difficulty = 20
		}

		fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

		if pokemon.BaseExperience-chanceToCatch <= difficulty {
			fmt.Println()
			fmt.Printf("%s was caught!\n", pokemon.Name)
			fmt.Println("You may now inspect it with the inspect command.")
			fmt.Println()

			pokeapi.Pokedex[pokemon.Name] = pokemon

		} else {
			fmt.Println()
			fmt.Printf("%s escaped!", pokemon.Name)
			fmt.Println()
		}
	} else {
		fmt.Println("The pokemon you are trying to catch doesnt exist")
	}

	return nil
}

func commandInspect(config *config, pokemonName ...string) error {

	if pokemon, ok := pokeapi.Pokedex[pokemonName[0]]; ok {
		fmt.Println("Name: ", pokemon.Name)
		fmt.Println("Height: ", pokemon.Height)
		fmt.Println("Weight: ", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stats := range pokemon.Stats {
			fmt.Printf(" -%s: %d\n", stats.Stat.Name, stats.BaseStat)
		}
		fmt.Println("Types:")
		for _, types := range pokemon.Types {
			fmt.Println(" -", types.Type.Name)
		}

	} else {
		fmt.Println("You have not caught this Pokemon")
	}

	return nil
}

func commandPokedex(config *config, blank ...string) error {
	fmt.Println("Your Pokedex:")
	if len(pokeapi.Pokedex) == 0 {
		fmt.Println(" - Your Pokedex is empty")
	}
	for _, value := range pokeapi.Pokedex {
		fmt.Println(" - ", value.Name)
	}
	return nil
}

func startRepl(config *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		userInput := scanner.Text()

		words := cleanInput(userInput)

		supportedCommands := config.commandsRegistry

		if command, ok := supportedCommands[words[0]]; ok {

			switch len(words) {
			case 2:
				err := command.callback(config, words[1])
				if err != nil {
					fmt.Println(err)
				}
			default:
				err := command.callback(config)
				if err != nil {
					fmt.Println(err)
				}
			}

		} else {
			fmt.Println("Unknown command.")
		}

	}
}
