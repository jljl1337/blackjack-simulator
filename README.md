# blackjack-simulator

[![Release](https://github.com/jljl1337/blackjack-simulator/actions/workflows/release.yml/badge.svg)](https://github.com/jljl1337/blackjack-simulator/actions/workflows/release.yml)
[![Source](https://img.shields.io/badge/Source-GitHub-blue?logo=github)](https://github.com/jljl1337/blackjack-simulator)
[![Docker](https://img.shields.io/badge/Docker-jljl1337%2Fblackjack--simulator-blue?logo=docker)](https://hub.docker.com/r/jljl1337/blackjack-simulator)
[![License](https://img.shields.io/github/license/jljl1337/blackjack-simulator?label=License
)](https://github.com/jljl1337/blackjack-simulator/blob/main/LICENSE)

## Features

- **Efficient**: Concurrency enabled by default, able to simulate millions of games in seconds.
- **Customizable**: Easily change the number of decks, penetration rate, and more.
- **Reproducible**: Exact same result can be achieved by using the same seed and version.
- **Simple**: No dependencies, released as just a single binary to run.

## Usage

### Flags

| Flag | Description |
| ---- | ----------- |
| `-config` | Path to the configuration file (default: `config.json`), see [Configuration](#configuration) for details. |
| `-csv` | Path to the CSV file to write results to if specified. |
| `-num-workers` | Number of concurrent shuffles to run (default: number of CPU cores). |
| `-verbose` | Enable verbose output. |

### Configuration

Sample configuration file can be found in the repository root as `config.json`.

| Field | Type | Description |
| ----- | ---- | ----------- |
| `seed` | `int64` | Random seed for reproducible results. If not provided or set to 0, uses current timestamp. |
| `numShuffles` | `uint` | Number of shuffles to simulate. |
| `numRounds` | `uint` | Number of rounds to simulate. |
| `numHands` | `uint` | Number of hands to simulate. |
| `numDecks` | `uint` | Number of decks in the shoe. Must be greater than 0. |
| `penetration` | `float64` | Shoe penetration percentage with a range of (0, 1]. Determines portion of the shoe that is dealt before reshuffling. |

> [!IMPORTANT]  
> The `numShuffles`, `numRounds`, and `numHands` fields are mutually exclusive, exactly one must be specified with a value greater than 0.