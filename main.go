package main

import (
	"fmt"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commandMap = map[string]cliCommand{}


func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commandMap = map[string]cliCommand {
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
	}
	
	for {
		fmt.Print("Pokedex > ")
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
		
		if command, ok := commandMap[safeInput[0]]; ok {
			if err := command.callback(); err != nil {
				fmt.Println("err:", err)
			}
			continue
		}
		fmt.Println("Unknown command:", safeInput[0])
	}
}


func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range commandMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}