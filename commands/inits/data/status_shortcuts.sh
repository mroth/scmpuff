scmpuff_status() {
  local scmpuff_env_char="e"

  # Ensure shwordsplit is on for zsh
  if [ -n "$ZSH_VERSION" ]; then setopt shwordsplit; fi;

  # Run scmpuff status, store output
  # (`local` needs to be on its own line otherwise exit code is swallowed!)
  local cmd_output
  cmd_output="$(/usr/bin/env scmpuff status --filelist $@)"

  # if there was an error, exit prematurely, and pass along the exit code
  # (STDOUT was swallowed but not STDERR, so user should still see error msg)
  local es=$?
  if [ $es -ne 0 ]; then
    return $es
  fi

  # Fetch list of files (from first line of script output)
  files="$(echo "$cmd_output" | head -n 1)"

  # Export numbered env variables for each file
  scmpuff_clear_vars
  IFS_CUR=$IFS
  # Some shells (pdksh) do not support ANSI C Quoting, i.e. $'\t'
  # This temporary IFS is a literal tab character, not a space. 
  IFS='	' 
  local e=1
  for file in $files; do
    export $scmpuff_env_char$e="$file"
    let e++
  done
  IFS=$IFS_CUR

  # Print status (from line two onward)
  echo "$cmd_output" | tail -n +2

  # Reset zsh environment to default
  if [ -n "$ZSH_VERSION" ]; then unsetopt shwordsplit; fi;
}


# Clear numbered env variables
scmpuff_clear_vars() {
  local scmpuff_env_char="e"
  local i=0
  local env_var_i=${scmpuff_env_char}${i}

  until [ "${env_var_i+is_null}" ]; do
    unset ${env_var_i}
    env_var_i=${scmpuff_env_char}$(( ++i ))
  done
}
