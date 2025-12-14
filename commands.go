package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/jcoughlin11/pokedexcli/internal/pokeapi"
)

type config struct {
	client  pokeapi.Client
	pokedex map[string]pokeapi.Pokemon
	nextUrl *string
	prevUrl *string
	rng     *rand.Rand
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

	KnownCommands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempts to catch the given pokemon.",
		callback:    commandCatch,
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

func commandCatch(cfg *config, pokemonName string) error {
	response, err := cfg.client.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Generate a random number in [0, 1). If that number is > baseExp
	// / 100, we catch it. This makes it so the higher the baseExp,
	// the lower the chance of catching it
	randN := cfg.rng.Float64()
	fmt.Printf("rand: %v\n", randN)
	fmt.Printf("baseexp: %v\n", response.BaseExperience)

	if randN > (float64(response.BaseExperience) / 1000.0) {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.pokedex[pokemonName] = response
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}
