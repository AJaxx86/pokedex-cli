package main

import (
	"fmt"
	"bufio"
	"os"
	"github.com/ajaxx86/pokedex-cli/internal/pokeapi"
	"math/rand"
)

type cliCommand struct {
	name string
	description string
	callback func(*cmdConfig, []string) error
	config *cmdConfig
}
type cmdConfig struct {
	nextAreaURL string
	prevAreaURL string
	catchPercent int
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
		"explore": {
			name: "explore",
			description: "Lists Pokemon encounters in the area (i.e. explore pastoria-city-area)",
			callback: commandExplore,
			config: cfg,
		},
		"catch": {
			name: "catch",
			description: "Tries to catch the Pokemon and add it to your Pokedex (i.e. catch pikachu)",
			callback: commandCatch,
			config: cfg,
		}
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
			if err := command.callback(cfg, safeInput[1:]); err != nil {
				fmt.Println("err:", err)
			}
			continue
		}
		fmt.Println("Unknown command:", safeInput[0])
	}
}


func commandExit(cfg *cmdConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp(cfg *cmdConfig, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range commands {
		fmt.Println(command.name + ":", command.description)
	}
	return nil
}


func commandMap(cfg *cmdConfig, args []string) error {
	areas, nextURL, prevURL, err := client.GetAreas(cfg.nextAreaURL)
	if err != nil {
		return err
	}

	cfg.nextAreaURL = nextURL
	cfg.prevAreaURL = prevURL
	for _, area := range areas {
		fmt.Println("-", area)
	}

	return nil
}


func commandMapBack(cfg *cmdConfig, args []string) error {
	areas, nextURL, prevURL, err := client.GetAreas(cfg.prevAreaURL)
	if err != nil {
		return err
	}

	cfg.nextAreaURL = nextURL
	cfg.prevAreaURL = prevURL
	for _, area := range areas {
		fmt.Println("-", area)
	}

	return nil
}


func commandExplore(cfg *cmdConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No area name entered")
	}

	fmt.Println("Exploring", args[0] + "...")
	pokemon, err := client.GetEncounters(args[0])
	if err != nil {
		return err
	}
	if len(pokemon) == 0 {
		fmt.Println("No encounters in this area.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, p := range pokemon {
		fmt.Println("-", p)
	}
	return nil
}


func commandCatch(cfg *cmdConfig, args []string) error {
	
	return nil
}
