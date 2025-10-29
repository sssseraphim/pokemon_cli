package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sssseraphim/pokedex/internal/pokeapi"
	"os"
	"strings"
	"time"
)

type config struct {
	clientPoke   pokeapi.Client
	nextLocation *string
	prevLocation *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

var commands map[string]cliCommand

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		clientPoke: pokeClient,
	}

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Shows all pokemons in the area",
			callback:    commandExplore,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := strings.Fields(scanner.Text())
		if len(words) == 0 {
			continue
		}
		command := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		val, ok := commands[command]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			err := val.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for name, _ := range commands {
		fmt.Printf("%s: %s\n", commands[name].name, commands[name].description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	locations, err := cfg.clientPoke.ListLocations(cfg.nextLocation)
	if err != nil {
		return err
	}

	cfg.nextLocation = locations.Next
	cfg.prevLocation = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocation == nil {
		return errors.New("you're on the first page")
	}

	locations, err := cfg.clientPoke.ListLocations(cfg.prevLocation)
	if err != nil {
		return err
	}

	cfg.nextLocation = locations.Next
	cfg.prevLocation = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("send me location")
	}
	name := args[0]
	location, err := cfg.clientPoke.ListPokemons(name)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", name)
	fmt.Println("Found Pokemon: ")
	for _, enc := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}
