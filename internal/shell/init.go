package shell

import "fmt"

const zshFunc = `
# wtt shell wrapper — lets "wtt cd/create/list" change directory
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

# ── Prompt integration ───────────────────────────────────────────────────────
# Automatically shows the current repo name in the right prompt.
# No manual RPROMPT configuration needed — just source this file.

# p10k custom segment (used when Powerlevel10k is active).
prompt_wtt() {
  local ctx
  ctx=$(wtt-bin context 2>/dev/null) || return
  [[ -z "$ctx" ]] && return
  p10k segment -f cyan -i $'\ue0a0' -t "$ctx"
}
instant_prompt_wtt() { prompt_wtt; }

# _wtt_precmd — updates _WTT_PS1 before every prompt render.
_wtt_precmd() {
  local ctx
  ctx=$(wtt-bin context 2>/dev/null)
  _WTT_PS1="${ctx:+%F{cyan}$'\ue0a0' ${ctx}%f}"
}

# _wtt_setup runs once on the first prompt render, after all plugins are loaded.
# By that point we know which prompt framework is active and can wire up correctly.
_wtt_setup() {
  # Remove ourselves so we only run once.
  precmd_functions=("${(@)precmd_functions:#_wtt_setup}")

  if typeset -f p10k > /dev/null 2>&1; then
    # Powerlevel10k is active.
    # Inject 'wtt' into the right-prompt elements if not already present,
    # then reload p10k so it picks up the new segment.
    if [[ ${POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS[(r)wtt]} != wtt ]]; then
      # Insert before the trailing 'end' element when possible.
      local end_idx=${POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS[(i)end]}
      if (( end_idx <= ${#POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS} )); then
        POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS=(
          "${POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS[@]:0:$((end_idx-1))}"
          wtt
          "${POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS[@]:$((end_idx-1))}"
        )
      else
        POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS+=(wtt)
      fi
      # Persist so p10k reload sees the change, then reload once.
      typeset -g POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS
      p10k reload 2>/dev/null
    fi
    return
  fi

  # Generic zsh (plain, oh-my-zsh without p10k, etc.):
  # Hook _wtt_precmd and inject _WTT_PS1 into RPROMPT automatically.
  precmd_functions+=(_wtt_precmd)
  _wtt_precmd  # populate immediately for the first prompt

  # Prepend our segment to whatever RPROMPT is already set.
  if [[ -n "$RPROMPT" ]]; then
    RPROMPT='${_WTT_PS1}${_WTT_PS1:+ }'"$RPROMPT"
  else
    RPROMPT='${_WTT_PS1}'
  fi
}

precmd_functions+=(_wtt_setup)
`

const bashFunc = `
# wtt shell wrapper — lets "wtt cd/create/list" change directory
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

# ── Prompt integration ───────────────────────────────────────────────────────
# Automatically prepends the current repo name to PS1.
# No manual PS1 configuration needed — just source this file.

_wtt_update_ps1() {
  local ctx
  ctx=$(wtt-bin context 2>/dev/null)
  if [[ -n "$ctx" ]]; then
    _WTT_PS1="\[\e[36m\]\ue0a0 ${ctx}\[\e[0m\] "
  else
    _WTT_PS1=""
  fi
}

# Prepend our segment to PS1 via PROMPT_COMMAND (runs before every prompt).
_WTT_ORIG_PS1="$PS1"
PS1='$(_wtt_update_ps1)${_WTT_PS1}'"$PS1"
`

const fishFunc = `
# wtt shell wrapper — lets "wtt cd/create/list" change directory
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

# ── Prompt integration ───────────────────────────────────────────────────────
# Automatically shows the current repo name in the right prompt.
# No manual configuration needed — just source this file.
#
# We define fish_right_prompt only if it is not already defined by the theme.
# Tide, Starship, etc. define their own fish_right_prompt; in those cases the
# user should call wtt_segment from within their theme's right-prompt hook.

function wtt_segment
  set ctx (wtt-bin context 2>/dev/null)
  if test -n "$ctx"
    set_color cyan
    printf '\ue0a0 %s ' $ctx
    set_color normal
  end
end

if not functions -q fish_right_prompt
  function fish_right_prompt
    wtt_segment
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
