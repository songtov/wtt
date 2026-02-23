# Claude Instructions for wtt

## Always QA After Making Changes
After implementing anything, **run the actual commands** to verify it works before calling it done.

Examples:
```sh
make build && ./wtt-bin version   # verify build + version
make test-local                   # verify test-local target output
make clean && ls wtt-bin          # verify clean removes binary
./wtt-bin --help                  # verify help output
```

Never just create a file and report success — run it and confirm the output is correct.

## Local Testing Workflow
```sh
make build                        # build local binary (v0.3.0-dev)
export PATH=$(pwd):$PATH          # put local binary first
eval "$(wtt-bin --init zsh)"      # redefine wtt shell function
wtt version                       # should show v0.3.0-dev
```

## Build
```sh
make build             # default: VERSION=v0.3.0-dev
make build VERSION=v0.3.0  # release build
make clean             # remove wtt-bin
```

## Project Structure
- Module: `github.com/songtov/wtt`
- Binary: `wtt-bin` (shell wrapper `wtt` calls this)
- Shell wrapper pattern: binary prints path → shell function `cd`s to it
