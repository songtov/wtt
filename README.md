# wtt — Git Worktree Manager

`wtt` wraps `git worktree` to create, navigate, and remove worktrees with minimal friction.

## Installation

```sh
brew tap songtov/tap
brew install wtt
```

To update to the latest version:

```sh
brew update && brew upgrade wtt
```

Then add the shell integration so `wtt` can `cd` for you.

**zsh** — add to `~/.zshrc`:
```sh
eval "$(wtt-bin --init zsh)"
```

**bash** — add to `~/.bashrc` or `~/.bash_profile`:
```sh
eval "$(wtt-bin --init bash)"
```

**fish** — add to `~/.config/fish/config.fish`:
```sh
wtt-bin --init fish | source
```

Reload your shell after (`source ~/.zshrc` or equivalent).

## Usage

```sh
# Create a worktree — auto-generates a branch name if omitted
wtt create
wtt create feature/login

# Navigate to a worktree (cds into the directory)
wtt feature/login

# Pick a worktree interactively with fzf, then navigate
wtt list

# Remove a worktree — pick interactively if no branch given
wtt remove
wtt remove feature/login
wtt remove --force feature/login   # skip confirmation
```

## Configuration

Create `.wtt.toml` in your repo root to customize behavior:

```toml
# Where to put worktrees (default: ../<repo>-worktrees)
worktree_dir = "../myproject-worktrees"

# Files to copy from main worktree into each new worktree (default: [".gitignore"])
copy_files = [".gitignore", ".env.local"]

# Directories to copy recursively
copy_dirs = []

# Commands to run inside the new worktree after creation
post_create = ["npm install"]
```

## How It Works

A child process can't change the parent shell's directory, so `wtt` ships as `wtt-bin` plus a shell wrapper function. When a command prints a directory path, the wrapper `cd`s into it — the same pattern used by `nvm` and `direnv`.

Worktrees land at `../<repo>-worktrees/<branch>/` by default. Slashes in branch names become dashes (`feature/login` → `feature-login`).

## Contributing

Contributions are welcome! Here's how to get started:

```sh
git clone https://github.com/songtov/wtt.git
cd wtt
go build -o wtt-bin .
```

- Bug reports and feature requests → [open an issue](https://github.com/songtov/wtt/issues)
- Pull requests → fork the repo, make your changes on a branch, and open a PR against `main`
- Please keep changes focused — one feature or fix per PR

## License

MIT — see [LICENSE](LICENSE) for details.
