package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jcoughlin11/pokedexcli/internal/pokeapi"
)

func main() {
	registerCommands()
	scanner := bufio.NewScanner(os.Stdin)

	client := pokeapi.NewClient(5 * time.Second)

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	pokedex := make(map[string]pokeapi.Pokemon, 0)

	cfg := config{client: client, rng: rng, pokedex: pokedex}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		cleaned := cleanInput(line)
		cmd, ok := KnownCommands[cleaned[0]]

		if ok {
			if len(cleaned) == 2 {
				err := cmd.callback(&cfg, cleaned[1])
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err := cmd.callback(&cfg, "")
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
