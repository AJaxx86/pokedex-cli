package main

import (
	"fmt"
	"bufio"
	"os"
)


func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		rawInput := scanner.Text()
		safeInput := cleanInput(rawInput)
		if len(safeInput) == 0 {
			fmt.Println("err: empty string entered")
			continue
		}
		
		fmt.Println("Your command was:", safeInput[0])
	}
}
