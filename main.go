package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jcoughlin11/pokedexcli/internal/pokeapi"
)

func main() {
	registerCommands()
	scanner := bufio.NewScanner(os.Stdin)

	client := pokeapi.NewClient(5 * time.Second)
	cfg := config{client: client}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		cleaned := cleanInput(line)
		cmd, ok := KnownCommands[cleaned[0]]

		if ok {
			err := cmd.callback(&cfg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
