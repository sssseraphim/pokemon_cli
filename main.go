package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sssseraphim/pokedex/internal/pokeapi"
	"math/rand"
	"os"
	"strings"
	"time"
)

type config struct {
	clientPoke   pokeapi.Client
	pokedex      map[string]pokeapi.Pokemon
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
		pokedex:    map[string]pokeapi.Pokemon{},
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect stats of a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all caught pokemon",
			callback:    commandPokedex,
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("catch who?")
	}
	name := args[0]
	pokemon, err := cfg.clientPoke.GetPokemon(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	attempt := rand.Intn(pokemon.BaseExperience)
	if attempt > 40 {
		fmt.Printf("%s escaped!\n", name)
	} else {
		fmt.Printf("%s is caught!\n", name)
		cfg.pokedex[name] = pokemon
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("inspect who?")
	}
	name := args[0]
	pokemon, ok := cfg.pokedex[name]
	if !ok {
		return errors.New("no such pokemon caught")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, poketype := range pokemon.Types {
		fmt.Printf("  -%s\n", poketype.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, _ ...string) error {
	fmt.Println("Your Pokedex:")
	for _, val := range cfg.pokedex {
		fmt.Printf(" - %s\n", val.Name)
	}
	return nil
}
