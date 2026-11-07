# Pokedex

An interactive command-line Pokedex REPL written in Go, backed by [PokeAPI](https://pokeapi.co). Browse locations, explore encounters, and catch Pokemon into your own Pokedex.

## Requirements

- Go 1.22+

## Getting started

```sh
go run .
# or build a binary and start it
go build -o pokedex . && ./pokedex
```

## Usage

```
$ ./pokedex
Pokedex > map
canalave-city-area
eterna-city-area
...
Pokedex > explore canalave-city-area
Tentacool
Tentacruel
...
```

All API responses are cached for 5 minutes (TTL cache in `pokecache/`), so repeat requests resolve instantly and don't re-fetch.

### Commands

| Command | Description |
| --- | --- |
| `help` | Show the usage message |
| `exit` | Quit the Pokedex |
| `map` | Next page of 20 location areas |
| `mapb` | Previous page of location areas |
| `explore <location-area>` | List the Pokemon that can be encountered there |
| `catch <pokemon>` | Try to catch a Pokemon and add it to your Pokedex |
| `inspect <pokemon>` | Show details (height, weight, stats, types) of a caught Pokemon |
| `pokedex` | List all the Pokemon you have caught |

### Catching Pokemon

Each catch throws a Pokeball and rolls `math/rand`: the higher the Pokemon's base experience, the harder it is to catch (caught when `rand.Intn(800) >= base experience`). Magikarp (~95%) falls easily, while species like Blissey (~24%) take several throws.

Once caught, details are stored in the Pokedex — `inspect` works offline from the stored data without additional API calls.

```
Pokedex > inspect pidgey
you have not caught that pokemon
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  -hp: 40
  -attack: 45
  -defense: 40
  -special-attack: 35
  -special-defense: 35
  -speed: 56
Types:
  - normal
  - flying
```

## Testing

```sh
go test ./...
```

Tests are offline and deterministic (API-free); run the tool itself for manual checks against the live API.

## Project layout

- Repo root: `main` package — REPL loop, command handlers (`command_*.go` / `command*.go`), and test helpers
- `pokecache/`: the concurrency-safe TTL response cache (only subpackage)
