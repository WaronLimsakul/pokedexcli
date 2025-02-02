package main

import (
	"bufio"
	"fmt"
	pokeApi "github.com/WaronLimsakul/pokedexcli/internal"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeApi.Config) error
}

var commandsMap map[string]cliCommand

func cleanInput(text string) []string {
	cleaned := strings.Fields(strings.ToLower(text))
	return cleaned
}

func commandExit(configp *pokeApi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configp *pokeApi.Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commandsMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(configp *pokeApi.Config) error {
	locationAreasResponse, err := pokeApi.FetchLocationAreas(configp.Next)
	if err != nil {
		return fmt.Errorf("Error fetching location areas: %v", err)
	}
	configp.Next = locationAreasResponse.Next
	configp.Previous = locationAreasResponse.Previous

	for _, area := range locationAreasResponse.Results {
		fmt.Printf("%s\n", area.Name)
	}
	return nil
}

func commandMapBack(configp *pokeApi.Config) error {
	if configp.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locationAreasResponse, err := pokeApi.FetchLocationAreas(configp.Previous)

	if err != nil {
		return fmt.Errorf("Error fetching location areas: %v", err)
	}
	configp.Next = locationAreasResponse.Next
	configp.Previous = locationAreasResponse.Previous

	for _, area := range locationAreasResponse.Results {
		fmt.Printf("%s\n", area.Name)
	}
	return nil
}

func main() {

	commandsMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokenmon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapBack,
		},
	}

	configp := &pokeApi.Config{Next: "https://pokeapi.co/api/v2/location-area/", Previous: ""}

	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Printf("Pokedex > ")
		scanSuccess := scanner.Scan()
		if !scanSuccess {
			return
		}
		text := scanner.Text()
		inputCommand := cleanInput(text)[0]
		command, ok := commandsMap[inputCommand]
		if ok {
			command.callback(configp)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
