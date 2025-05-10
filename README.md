# Gozelle

**Gozelle** is a lightning-fast, minimal directory-jumping tool written in Go — inspired by [`zoxide`](https://github.com/ajeetdsouza/zoxide). Jump to frequently used directories with just a keyword, powered by **frecency scoring**, **fuzzy matching**, and **shell integration**.

---

## Table of Contents

- [Features](#Features)
- [Requirements](#Requirements)
- [Installation](#Installation)
- [Usage](#Usage)
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

- [Bash](https://www.gnu.org/software/bash/) or [Zsh](https://www.zsh.org/) shell  
- Gozelle binary in your `$PATH`  
- go version 1.24+

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
`/usr/local/bin` can be replace by any direcotry in your `$PATH`

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

[↑ Back to top](#Gozelle)

---

## Usage

### Jump to a Directory

```bash
gz proj       # jumps to the best match (e.g., ~/projects)
```

### Show Matching Directories (without jumping)

```bash
gozelle query proj
```

### Add a Directory Manually

```bash
gozelle add /some/path/to/add
```

[↑ Back to top](#Gozelle)

---

## How It Works

- Enables shell hooks for automatic logging via `init bash` or `init zsh`  
- Tracks every visited directory using the shell hook  
- Stores them in a gob-encoded file under your user data directory (default is `~/.local/share/Gozelle` for Linux users)  
- Finds all matches for keywords entered, e.g., `gz keywords`  
- Ranks them using a **frecency** score (frequency + recency)

[↑ Back to top](#Gozelle)

---

## Learnings & Concepts

This project is a hands-on learning opportunity for:

- **Concurrency** — handling simultaneous updates and queries efficiently  
- **Worker Pools** — to process background updates to scoring  
- **Mutexes** — for safe access to shared resources (like the gob database)  
- **Command Line Hooks** — shell integration and behavior injection  
- **Gob Encoding** — simple and efficient binary data serialization in Go  

[↑ Back to top](#Gozelle)

---

## Roadmap

- [x] Zsh support  
- [ ] Fish shell support  
- [ ] Interactive `fzf`-style selector  
- [ ] Configurable data file location  
- [ ] Directory expiration / pruning logic  
- [ ] Man Pages  
- [ ] Completion support

[↑ Back to top](#Gozelle)

---

## 📄 License


GNU General Public License 3.0 — see [LICENSE](LICENSE) for details.


[↑ Back to top](#Gozelle)
