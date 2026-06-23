package main

import (
	"fmt"
	"bufio"
	"os"
	"github.com/ajaxx86/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func(*cmdConfig) error
	config *cmdConfig
}
type cmdConfig struct {
	nextAreaURL string
	prevAreaURL string
}

var commands map[string]cliCommand
var client pokeapi.Client


func main() {
	scanner := bufio.NewScanner(os.Stdin)
	client = pokeapi.NewClient()
	cfg := &cmdConfig{
		nextAreaURL: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		prevAreaURL: "",
	}
	commands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
			config: cfg,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
			config: cfg,
		},
		"map": {
			name: "map",
			description: "Gets the next 20 areas from the Pokemon world",
			callback: commandMap,
			config: cfg,
		},
		"mapb": {
			name: "mapb",
			description: "Gets the last 20 areas from the Pokemon world",
			callback: commandMapBack,
			config: cfg,
		},
	}

	for {
		fmt.Print("\nPokedex > ")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println("err:", scanner.Err())
			continue
		}
		rawInput := scanner.Text()
		safeInput := cleanInput(rawInput)
		if len(safeInput) == 0 {
			fmt.Println("err: empty string entered")
			continue
		}

		if command, ok := commands[safeInput[0]]; ok {
			if err := command.callback(cfg); err != nil {
				fmt.Println("err:", err)
			}
			continue
		}
		fmt.Println("Unknown command:", safeInput[0])
	}
}


func commandExit(cfg *cmdConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp(cfg *cmdConfig) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range commands {
		fmt.Println(command.name + ":", command.description)
	}
	return nil
}


func commandMap(cfg *cmdConfig) error {
	areas, nextURL, prevURL, err := client.GetAreas(cfg.nextAreaURL)
	if err != nil {
		return err
	}

	cfg.nextAreaURL = nextURL
	cfg.prevAreaURL = prevURL
	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}


func commandMapBack(cfg *cmdConfig) error {
	areas, nextURL, prevURL, err := client.GetAreas(cfg.prevAreaURL)
	if err != nil {
		return err
	}

	cfg.nextAreaURL = nextURL
	cfg.prevAreaURL = prevURL
	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}
