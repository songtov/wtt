package shell

import "fmt"

const zshFunc = `
wtt() {
  local output
  output=$(wtt-bin "$@" 2>/dev/tty)
  local exit_code=$?
  if [ $exit_code -ne 0 ]; then
    return $exit_code
  fi
  if [ -d "$output" ]; then
    cd "$output" || return 1
  else
    echo "$output"
  fi
}
`

const bashFunc = `
wtt() {
  local output
  output=$(wtt-bin "$@" 2>/dev/tty)
  local exit_code=$?
  if [ $exit_code -ne 0 ]; then
    return $exit_code
  fi
  if [ -d "$output" ]; then
    cd "$output" || return 1
  else
    echo "$output"
  fi
}
`

const fishFunc = `
function wtt
  set output (wtt-bin $argv 2>/dev/tty)
  set exit_code $status
  if test $exit_code -ne 0
    return $exit_code
  end
  if test -d "$output"
    cd "$output"
  else
    echo "$output"
  end
end
`

// InitFunction returns the shell function definition for the given shell.
func InitFunction(sh string) (string, error) {
	switch sh {
	case "zsh":
		return zshFunc, nil
	case "bash":
		return bashFunc, nil
	case "fish":
		return fishFunc, nil
	default:
		return "", fmt.Errorf("unsupported shell %q; supported: zsh, bash, fish", sh)
	}
}
