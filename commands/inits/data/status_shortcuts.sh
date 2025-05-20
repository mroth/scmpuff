# source file: status_shortcuts.sh
# shellcheck shell=sh

scmpuff_status() {
  __scmpuff_env_char="e"

  # Ensure shwordsplit is on for zsh
  if [ -n "$ZSH_VERSION" ]; then setopt shwordsplit; fi;

  # Run scmpuff status, store output
  __scmpuff_status_cmd_output=$(/usr/bin/env scmpuff status --filelist "$@")
  __scmpuff_status_exit=$?

  # if there was an error, exit prematurely, and pass along the exit code
  # (STDOUT was swallowed but not STDERR, so user should still see error msg)
  if [ $__scmpuff_status_exit -ne 0 ]; then
    return $__scmpuff_status_exit
  fi

  # Fetch list of files (from first line of script output)
  __scmpuff_files="$(echo "$__scmpuff_status_cmd_output" | head -n 1)"

  # Export numbered env variables for each file
  scmpuff_clear_vars
  IFS=$(printf '\t')
  __scmpuff_loop_e=1
  for __scmpuff_file in $__scmpuff_files; do
    export $__scmpuff_env_char$__scmpuff_loop_e="$__scmpuff_file"
    __scmpuff_loop_e=$((__scmpuff_loop_e+1)) #e++
  done
  IFS=$(printf ' \t\n')

  # Print status (from line two onward)
  echo "$__scmpuff_status_cmd_output" | tail -n +2

  # Reset zsh environment to default
  if [ -n "$ZSH_VERSION" ]; then unsetopt shwordsplit; fi;
}

# Clear numbered env variables
scmpuff_clear_vars() {
  __scmpuff_env_char="e"
  __scmpuff_loop_i=0
  while [ $__scmpuff_loop_i -le 999 ]; do
    __scmpuff_env_var_i=${__scmpuff_env_char}${__scmpuff_loop_i}
    if [ -n "$__scmpuff_env_var_i" ]; then
      unset "$__scmpuff_env_var_i"
    else
      break
    fi
    __scmpuff_loop_i=$((__scmpuff_loop_i+1)) #i++
  done
}
