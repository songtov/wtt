package shell

import "fmt"

const zshFunc = `
# wtt shell wrapper — lets "wtt create/list/remove" cd into the selected path
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

# wtt_ps1 — like kube_ps1, works with ANY prompt theme.
#
# Usage (pick one):
#   Plain zsh / oh-my-zsh:    RPROMPT='$(wtt_ps1)'
#   Right side of p10k:       add 'wtt' to POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS
#                             (p10k segment defined below)
#   Inline anywhere:          PS1="$PS1\$(wtt_ps1) "
#
# The function outputs nothing when no repo context is active.
wtt_ps1() {
  local ctx
  ctx=$(wtt-bin context 2>/dev/null) || return
  [[ -z "$ctx" ]] && return
  # Nerd Font git-branch icon (\ue0a0). Falls back gracefully in non-NF fonts.
  echo -n "%F{cyan}\ue0a0 ${ctx}%f"
}

# Powerlevel10k custom segment.
# If you use p10k add 'wtt' to POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS in ~/.p10k.zsh
prompt_wtt() {
  local ctx
  ctx=$(wtt-bin context 2>/dev/null) || return
  [[ -z "$ctx" ]] && return
  p10k segment -f cyan -i $'\ue0a0' -t "$ctx"
}
# Enable instant-prompt support (p10k calls this during startup)
instant_prompt_wtt() { prompt_wtt; }
`

const bashFunc = `
# wtt shell wrapper — lets "wtt create/list/remove" cd into the selected path
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

# wtt_ps1 — like kube_ps1, works with ANY prompt theme.
#
# Usage:
#   PS1="${PS1}\$(wtt_ps1) "
#   or add to PROMPT_COMMAND to update before each prompt.
#
# The function outputs nothing when no repo context is active.
wtt_ps1() {
  local ctx
  ctx=$(wtt-bin context 2>/dev/null) || return
  [[ -z "$ctx" ]] && return
  # Nerd Font git-branch icon; \e[36m = cyan, \e[0m = reset
  printf '\e[36m\ue0a0 %s\e[0m' "$ctx"
}
`

const fishFunc = `
# wtt shell wrapper — lets "wtt create/list/remove" cd into the selected path
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

# wtt_ps1 — call inside fish_right_prompt (or fish_prompt) to show repo context.
#
# Usage: add to ~/.config/fish/functions/fish_right_prompt.fish:
#   wtt_ps1
#
function wtt_ps1
  set ctx (wtt-bin context 2>/dev/null)
  if test -n "$ctx"
    set_color cyan
    printf '\ue0a0 %s' $ctx
    set_color normal
  end
end
`

// InitFunction returns the shell function definitions for the given shell.
// Source the output in your rc file:
//
//	eval "$(wtt-bin --init zsh)"
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
