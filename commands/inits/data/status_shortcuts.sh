scmpuff_status_shortcuts() {
  local scmpuff_env_char="e"

  # fail_if_not_git_repo || return 1

  # Ensure shwordsplit is on for zsh
  if [ -n "$ZSH_VERSION" ]; then setopt shwordsplit; fi;

  scmpuff_clear_vars

  # Run scmpuff status, store output
  local cmd_output="$(/usr/bin/env scmpuff status --filelist $@)"

  # Fetch list of files from first line of script output
  files="$(echo "$cmd_output" | head -n 1)"

  # Export numbered env variables for each file
  IFS="|"
  local e=1
  for file in $files; do
    export $scmpuff_env_char$e="$file"
    if [ "${scmpuffDebug:-}" = "true" ]; then echo "Set \$$scmpuff_env_char$e  => $file"; fi
    let e++
  done
  IFS=$' \t\n'

  if [ "${scmpuffDebug:-}" = "true" ]; then echo "------------------------"; fi

  # Print status (from line two onward)
  echo "$cmd_output" | tail -n +2

  # Reset zsh environment to default
  if [ -n "$ZSH_VERSION" ]; then unsetopt shwordsplit; fi;
}


# Clear numbered env variables
scmpuff_clear_vars() {
  local scmpuff_env_char="e"
  local i

  for (( i=1; i<=999; i++ )); do
    local env_var_i=${scmpuff_env_char}${i}
    if [[ -n ${env_var_i} ]]; then
      unset ${env_var_i}
    else
      break
    fi
  done
}
