package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandMap map[string]cliCommand

func cleanInput(text string) []string {
	cleaned := strings.Fields(strings.ToLower(text))
	return cleaned
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commandMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func main() {

	commandMap = map[string]cliCommand{
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
	}

	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Printf("Pokedex > ")
		scanSuccess := scanner.Scan()
		if !scanSuccess {
			return
		}
		text := scanner.Text()
		inputCommand := cleanInput(text)[0]
		command, ok := commandMap[inputCommand]
		if ok {
			command.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
	cleaned := cleanInput("  Hello World  ")
	fmt.Printf("text: '  Hello World  ', cleaned: %v\n", cleaned)
	fmt.Printf("length: %v", len(cleaned))
}
