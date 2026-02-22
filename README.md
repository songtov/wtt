# wtt — Git Worktree Manager

`wtt` wraps `git worktree` to create, navigate, and remove worktrees with minimal friction.

## Features

- **`wtt create [branch]`** — Create a worktree. Auto-generates a branch name if omitted.
- **`wtt <branch>`** — Navigate (cd) to a worktree by branch name.
- **`wtt list`** — Interactive picker (fzf) to select and navigate to a worktree.
- **`wtt remove <branch>`** — Remove a worktree with confirmation.
- **Config file** (`.wtt.toml`) — Customize worktree directory, files to copy, and post-create hooks.

## Installation

### Homebrew (macOS / Linux)

```sh
brew install songtov/wtt/wtt
```

### Shell integration (required for navigation)

Add to your shell rc file:

**zsh** (`~/.zshrc`):
```sh
eval "$(wtt-bin --init zsh)"
```

**bash** (`~/.bashrc` or `~/.bash_profile`):
```sh
eval "$(wtt-bin --init bash)"
```

**fish** (`~/.config/fish/config.fish`):
```sh
wtt-bin --init fish | source
```

Then reload your shell: `source ~/.zshrc` (or equivalent).

## Usage

```sh
# Create a worktree with a given branch name
wtt create feature/login

# Create a worktree with an auto-generated branch name (e.g. myrepo-bright-falcon)
wtt create

# Navigate to a worktree (shell function cds into the directory)
wtt feature/login

# Interactively pick a worktree (requires fzf)
wtt list

# Remove a worktree (with confirmation prompt)
wtt remove feature/login

# Remove without confirmation
wtt remove --force feature/login

# Print version
wtt version
```

## Configuration

Create `.wtt.toml` in your repository root to customize behavior:

```toml
# Default: ../<repo-name>-worktrees
worktree_dir = "../myproject-worktrees"

# Files to copy from the main worktree to each new worktree
# Default: [".gitignore"]
copy_files = [".gitignore", ".env.local"]

# Directories to copy (recursive)
copy_dirs = []

# Commands to run in the new worktree after creation
post_create = ["npm install"]
```

## How It Works

A child process cannot change the parent shell's directory. `wtt` ships a binary named `wtt-bin` and a shell wrapper function called `wtt`. When a command outputs a directory path to stdout, the shell wrapper `cd`s into it. This is the same pattern used by `nvm` and `direnv`.

Worktrees are created at `../<repo-name>-worktrees/<branch>/` by default. Branch slashes are replaced with dashes (e.g. `feature/login` → `feature-login`).

## Building from Source

```sh
go build -o wtt-bin .
```

To inject the version at build time:

```sh
go build -ldflags "-X github.com/songtov/wtt/cmd.Version=v1.0.0" -o wtt-bin .
```

To build release artifacts with goreleaser:

```sh
goreleaser release --snapshot --clean
```

## License

MIT
