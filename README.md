# Gozelle

[![Go Report Card](https://goreportcard.com/badge/github.com/ATLIOD/Gozelle?t=1)](https://goreportcard.com/report/github.com/ATLIOD/Gozelle)

**Gozelle** is a lightning-fast and minimal smart cd command written in Go — inspired by [`zoxide`](https://github.com/ajeetdsouza/zoxide). Jump to frequently used directories with just a keyword, powered by **frecency scoring**, **fuzzy matching**, and **shell integration**.

---

## Table of Contents

- [Features](#Features)
- [Requirements](#Requirements)
- [Installation](#Installation)
- [Usage](#Usage)
- [Environment Variables](#Environment-Variables)
- [How It Works](#How-it-works)
- [Learnings & Concepts](#Learnings--concepts)
- [Roadmap](#Roadmap)
- [License](#License)

---

## Features

- **Frecency Scoring** — jump history is ranked by frequency and recency  
- **Fuzzy Matching** — jump with just a keyword or part of a directory name  
- **Smart Ranking** — most relevant paths surface first  
- **Manual Add** — add directories to the index yourself  
- **Query Mode** — list matching directories without jumping  
- **Compact Storage** — gob-encoded data stored locally  
- **Shell Integration** — Bash and Zsh command-line hooks for seamless tracking

[↑ Back to top](#Gozelle)

---

## Requirements

- [Bash](https://www.gnu.org/software/bash/), [Zsh](https://www.zsh.org/), or [Fish](https://fishshell.com/) shell  
- Gozelle binary in your `$PATH`  
- go version 1.24+
- [`fzf`](https://github.com/junegunn/fzf) installed for interactive mode

**Platform Support Notice:**  
> Gozelle has currently only been tested and verified on Linux systems. While it may work on other Unix-like OSes or Windows, no official support or testing has been done outside of Linux.


[↑ Back to top](#Gozelle)

---

## Installation

### Option 1: Build from Source

```bash
git clone https://github.com/yourusername/gozelle
cd gozelle
go build -o gozelle .
sudo mv gozelle /usr/local/bin/
```
`/usr/local/bin` can be replace by any directory in your `$PATH`

### Option 2: Use `go install`

```bash
go install github.com/ATLIOD/Gozelle@latest
```

Make sure your `$GOPATH/bin` is in your `$PATH` (commonly `~/go/bin`):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Add to Shell Startup File

For **Bash**:

```bash
echo 'eval "$(gozelle init bash)"' >> ~/.bashrc
source ~/.bashrc
```

For **Zsh**:

```bash
echo 'eval "$(gozelle init zsh)"' >> ~/.zshrc
source ~/.zshrc
```

For **Fish**:

```bash
echo 'gozelle init fish | source' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```

[↑ Back to top](#Gozelle)

---

## Usage

### Jump to a Directory

```bash
gz projects       # jumps to the best match (e.g., ~/Documents/School/Programming/projects)
```

### Show Matching Directories (without jumping)

```bash
gozelle query projects
```

### Add a Directory Manually

```bash
gozelle add /some/path/to/add
```
### Interactive Mode
```bash
gi
```

[↑ Back to top](#Gozelle)

---
## Environment Variables

| Variable          | Description                                                                                 | Default                                                                                           |
|-------------------|---------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| `GOZELLE_ECHO`    | Whether to print the target directory path to stdout after jumping. Must be `"true"` or `"false"`. | `"false"` (default)                                                                             |
| `GOZELLE_DATA_DIR`| Path to the directory where Gozelle stores its data file (`db.gob`). If not set, defaults to: <br> `$XDG_DATA_HOME/gozelle/db.gob` <br> or `<home>/.local/share/gozelle/db.gob` if `$XDG_DATA_HOME` is unset. | `~/.local/share/gozelle/db.gob` (default)                                                       |

### Notes

- `GOZELLE_ECHO` must be set to exactly `"true"` or `"false"`. Any other value will reset it to `"false"` and print a warning.
- When `GOZELLE_DATA_DIR` is set to a non-existent directory, Gozelle will attempt to create the directory automatically in the uses default data directory. ~/.local/share in linux)
- The environment variable `GOZELLE_DATA_DIR` refers to a directory; the actual data file is stored inside it as `db.gob`.

### Example usage

```bash
export GOZELLE_ECHO=true
export GOZELLE_DATA_DIR="$HOME/.config/gozelle"
```

---

## How It Works

- Enables shell hooks for automatic logging via `init bash` or `init zsh`  
- Tracks every visited directory using the shell hook  
- Stores them in a gob-encoded file under your user data directory (default is `~/.local/share/Gozelle` for Linux users)  
- Finds all matches for keywords entered, e.g., `gz keywords`  
- Ranks them using a **frecency** score (frequency + recency)
- Uses fzf to provide an interactive selection UI when requested

[↑ Back to top](#Gozelle)

---

## Learnings & Concepts

This project is a hands-on learning opportunity for:

- **Concurrency** — handling simultaneous updates and queries efficiently  
- **Worker Pools** — to process background updates to scoring  
- **Mutexes** — for safe access to shared resources (like the gob database)  
- **Command Line Hooks** — shell integration and behavior injection  
- **Gob Encoding** — simple and efficient binary data serialization in Go
- **Integration with External Tools** — incorporating fzf for interactive mode

[↑ Back to top](#Gozelle)

---

## Roadmap

- [X] Zsh support  
- [X] Fish shell support  
- [X] Interactive `fzf`-style selector  
- [X] Configurable data file location  
- [X] Directory expiration / pruning logic  
- [X] Man Page  
- [X] Completion support
- [ ] Better pruning logic
- [X] Higher weight to paths where the keyword is closer to the end

[↑ Back to top](#Gozelle)

---

## License


GNU General Public License 3.0 — see [LICENSE](LICENSE) for details.


[↑ Back to top](#Gozelle)
