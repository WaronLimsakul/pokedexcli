package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	pokeApi "github.com/WaronLimsakul/pokedexcli/internal/pokeapi"
	"github.com/WaronLimsakul/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	// in Go, we don't have to use all parameters.
	callback func(*pokeApi.Config, []string) error
}

var commandsMap map[string]cliCommand
var mainCache *pokecache.Cache
var pokedex map[string]*pokeApi.PokemonDetail

func cleanInput(text string) []string {
	cleaned := strings.Fields(strings.ToLower(text))
	return cleaned
}

func commandExit(configp *pokeApi.Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configp *pokeApi.Config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commandsMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(configp *pokeApi.Config, args []string) error {
	locationAreasResponse, err := pokeApi.FetchLocationAreas(configp.Next, mainCache)
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

func commandMapBack(configp *pokeApi.Config, args []string) error {
	if configp.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locationAreasResponse, err := pokeApi.FetchLocationAreas(configp.Previous, mainCache)

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

func commandExplore(configp *pokeApi.Config, args []string) error {
	if len(args) == 0 || args[0] == "" {
		fmt.Println("location name required")
		return fmt.Errorf("location name required")
	}
	locationInput := args[0]

	fmt.Printf("Exploring %s...\n", locationInput)

	exploreAreaResponse, err := pokeApi.ExploreAreaPokemons(locationInput, mainCache)
	if err != nil {
		fmt.Printf("%v\n", err)
		return fmt.Errorf("error explore area's pokemon: %v", err)
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range exploreAreaResponse.PokemonsEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(configp *pokeApi.Config, args []string) error {
	if len(args) == 0 || args[0] == "" {
		fmt.Printf("pokemon name required\n")
		return fmt.Errorf("pokemon name required")
	}

	nameInput := args[0]

	_, ok := pokedex[nameInput]

	if ok {
		fmt.Printf("we already have %s!\n", nameInput)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", nameInput)

	// get pokemon data and random the possibility
	pokemonInfo, err := pokeApi.GetPokemonDetail(nameInput, mainCache)
	if err != nil {
		fmt.Printf("error cathcing pokemon: %v\n", err)
		return fmt.Errorf("error cathcing pokemon")
	}

	baseExp := pokemonInfo.BaseExperience
	pokemonName := pokemonInfo.Name
	randNum := rand.Intn(370)

	// this means catch success
	if baseExp < randNum {
		pokedex[pokemonName] = pokemonInfo
		fmt.Printf("%s was caught!\n", pokemonName)
		return nil
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}
}

func commandInspect(configp *pokeApi.Config, args []string) error {
	if len(args) == 0 || args[0] == "" {
		fmt.Printf("pokemon name required\n")
		return fmt.Errorf("pokemon name required")
	}

	nameInput := args[0]
	pokemonInfo, ok := pokedex[nameInput]
	if !ok {
		fmt.Printf("you have not caught %s yet...\n", nameInput)
		return nil
	}

	fmt.Printf("Name: %s\n", pokemonInfo.Name)
	fmt.Printf("Height: %v\n", pokemonInfo.Height)
	fmt.Printf("Weight: %v\n", pokemonInfo.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonInfo.Stats {
		fmt.Printf("    -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemonInfo.Types {
		fmt.Printf("    -%s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(configp *pokeApi.Config, args []string) error {
	if len(pokedex) == 0 {
		fmt.Println("You have not caught a Pokemon yet.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemonDetail := range pokedex {
		fmt.Printf(" - %s\n", pokemonDetail.Name)
	}
	return nil
}

func main() {
	pokedex = map[string]*pokeApi.PokemonDetail{}
	mainCache = pokecache.NewCache(10 * time.Second)
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
		"explore": {
			name:        "explore",
			description: "Displays a list of all Pokemon located in the area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays your caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of all the names of the Pokemon the user has caught.",
			callback:    commandPokedex,
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
		fullCommand := cleanInput(text)
		inputCommand := fullCommand[0]
		command, ok := commandsMap[inputCommand]
		if ok {
			command.callback(configp, fullCommand[1:])
		} else {
			fmt.Println("Unknown command")
		}
	}
}
