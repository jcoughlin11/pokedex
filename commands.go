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
	callback    func(*config, string) error
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

	KnownCommands["explore"] = cliCommand{
		name:        "explore",
		description: "Lists pokemon available in the given area.",
		callback:    commandExplore,
	}
}

func commandExit(cfg *config, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, _ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for name, cmd := range KnownCommands {
		fmt.Printf("%v: %v\n", name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, _ string) error {
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

func commandMapB(cfg *config, _ string) error {
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

func commandExplore(cfg *config, areaName string) error {
	response, err := cfg.client.ListPokemon(areaName)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range response.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
