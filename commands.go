package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jcoughlin11/pokedexcli/internal/pokeapi"
)

type config struct {
	client  pokeapi.Client
	nextUrl *string
	prevUrl *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var KnownCommands = make(map[string]cliCommand)

func registerCommands() {
	KnownCommands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	KnownCommands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	KnownCommands["map"] = cliCommand{
		name:        "map",
		description: "Displays the next page of location areas.",
		callback:    commandMap,
	}

	KnownCommands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous page of location areas.",
		callback:    commandMapB,
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for name, cmd := range KnownCommands {
		fmt.Printf("%v: %v\n", name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	response, err := cfg.client.ListLocations(cfg.nextUrl)
	if err != nil {
		return err
	}

	cfg.nextUrl = response.Next
	cfg.prevUrl = response.Prev

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(cfg *config) error {
	if cfg.prevUrl == nil {
		return errors.New("you're on the first page")
	}

	response, err := cfg.client.ListLocations(cfg.prevUrl)
	if err != nil {
		return err
	}

	cfg.nextUrl = response.Next
	cfg.prevUrl = response.Prev

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	return nil
}
