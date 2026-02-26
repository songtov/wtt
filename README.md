# wtt — Git Worktree Manager

A fast, opinionated Git worktree manager.

**wtt** = **w**ork**t**ree-**tov**

---

## Installation

```sh
brew tap songtov/tap
brew install wtt
```

To update:

```sh
brew update && brew upgrade wtt
```

[fzf](https://github.com/junegunn/fzf) is optional but recommended — used for interactive pickers.

### Shell integration

The shell wrapper is what lets `wtt` `cd` into worktrees. Add one of these to your shell config and reload.

**zsh** — `~/.zshrc`:
```sh
eval "$(wtt-bin --init zsh)"
```

**bash** — `~/.bashrc` or `~/.bash_profile`:
```sh
eval "$(wtt-bin --init bash)"
```

**fish** — `~/.config/fish/config.fish`:
```sh
wtt-bin --init fish | source
```

`--init` also wires up [shell prompt integration](#shell-prompt) automatically — no extra config needed.

---

## Quick Start

```sh
# Create a worktree with a random branch name, then navigate into it
wtt create

# Create a worktree for a specific branch
wtt create feature/login

# Create a worktree branching from a specific ref
wtt create feature/login -b main

# Pick a worktree interactively and navigate to it
wtt list

# Remove a worktree (interactive picker if no branch given)
wtt remove feature/login
```

---

## Commands

| Command | Description |
|---|---|
| `wtt create [branch]` | Create a new worktree |
| `wtt list` | List worktrees and navigate interactively |
| `wtt remove [branch]` | Remove a worktree |
| `wtt <branch>` | Navigate directly to a worktree |
| `wtt init` | Scaffold a `.wtt.toml` config file |
| `wtt repo list` | Switch the active repository context |
| `wtt repo remove` | Remove a repo from the known repos list |
| `wtt version` | Print version |

### `wtt create [branch]`

Creates a new worktree and navigates into it.

```sh
wtt create                      # random name (e.g. myrepo-crisp-summit)
wtt create feature/login        # specific branch
wtt create feature/login -b main  # branch from main instead of HEAD
```

| Flag | Description |
|---|---|
| `-b, --base <ref>` | Base commit/branch/ref (default: `HEAD`) |

### `wtt list`

Opens an fzf picker to select and navigate to a worktree. Falls back to a numbered list if fzf is not installed.

```sh
wtt list
```

### `wtt remove [branch]`

Removes a worktree. Opens an interactive picker if no branch is given.

```sh
wtt remove                      # pick interactively
wtt remove feature/login
wtt remove -f feature/login     # skip confirmation + force-remove from git
```

| Flag | Description |
|---|---|
| `-f, --force` | Skip confirmation prompt and pass `--force` to `git worktree remove` |

### `wtt <branch>`

Navigate directly to a worktree by branch name.

```sh
wtt feature/login
```

### `wtt init`

Scaffolds a `.wtt.toml` in the repo root with commented-out defaults.

```sh
wtt init          # fails if .wtt.toml already exists
wtt init -f       # overwrite existing file
```

| Flag | Description |
|---|---|
| `-f, --force` | Overwrite an existing `.wtt.toml` |

### `wtt repo list`

Switches the active repository context — kubens-style. After switching, all `wtt` commands (`create`, `list`, `remove`) operate on the selected repo, even when run from outside it.

```sh
wtt repo list
```

### `wtt repo remove`

Removes a repository from the known repos list.

```sh
wtt repo remove
wtt repo remove -f    # skip confirmation
```

| Flag | Description |
|---|---|
| `-f, --force` | Skip confirmation prompt |

---

## Configuration

Run `wtt init` to create a `.wtt.toml` in your repo root, then edit it:

```sh
wtt init
```

### Config keys

| Key | Type | Default | Description |
|---|---|---|---|
| `worktree_dir` | string | `../<repo>-worktrees` | Directory where worktrees are created |
| `copy_files` | list | `[".gitignore"]` | Files copied from the main worktree into each new worktree |
| `copy_dirs` | list | `[]` | Directories copied recursively into each new worktree |
| `symlink_files` | list | `[]` | Files symlinked (not copied) — changes in one worktree are shared across all |
| `post_create` | list | `[]` | Shell commands run inside the new worktree after creation |

### Example `.wtt.toml`

```toml
# Where to create worktrees (default: ../<repo>-worktrees)
worktree_dir = "../myproject-worktrees"

# Files copied into each new worktree
copy_files = [".gitignore", ".env.local"]

# Directories copied recursively
copy_dirs = ["scripts"]

# Files symlinked — edits in any worktree are reflected everywhere
symlink_files = [".env.secrets"]

# Commands run inside the new worktree after creation
post_create = ["npm install"]
```

---

## Shell Prompt

`--init` automatically injects the active repository name into your shell prompt — no manual configuration needed.

| Shell | Where it appears | Framework notes |
|---|---|---|
| zsh | Right prompt (`RPROMPT`) | Auto-injects into Powerlevel10k as a custom segment; falls back to generic `RPROMPT` for plain zsh and oh-my-zsh |
| bash | Left prompt (`PS1`) | Prepended via `PROMPT_COMMAND` |
| fish | Right prompt | Defines `fish_right_prompt` if not already set; for Tide/Starship call `wtt_segment` from your theme hook |

The prompt segment is powered by `wtt-bin context` internally. There is no separate `wtt context` command intended for direct use.

---

## How It Works

A child process can't change the parent shell's directory, so `wtt` ships as two parts:

1. **`wtt-bin`** — the Go binary that does the work
2. **`wtt`** — a shell function (installed by `--init`) that runs `wtt-bin` and `cd`s into the output if it's a directory

This is the same pattern used by tools like `nvm` and `direnv`.

Worktrees land at `../<repo>-worktrees/<branch>/` by default. Slashes in branch names become dashes (`feature/login` → `feature-login`).

When fzf is not installed, interactive pickers fall back to a numbered list — `wtt list` and `wtt remove` always work regardless.

---

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
