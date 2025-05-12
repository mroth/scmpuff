function scmpuff_status {
  $scmpuff_env_char = "e"

  # Run scmpuff status command and capture the output
  $cmd_output = & scmpuff status --filelist @args

  # if there was an error, exit prematurely, and pass along the exit code
  $es = $LastExitCode
  if ($es -ne 0) {
    return $es
  }

  # Fetch list of files (from first line of script output)
  $files = ($cmd_output | Select-Object -First 1) -split '\s+'

  # Export numbered env variables for each file
  scmpuff_clear_vars
  $e = 1
  foreach ($file in $files) {
    Set-Item "env:$($scmpuff_env_char + $e)" $file
    $e++
  }

  # Print status (from line two onward)
  $cmd_output | Select-Object -Skip 1
}

# Clear numbered env variables
function scmpuff_clear_vars {
  # Define the environment variable character
  $scmpuff_env_char = "e"

  # Iterate through environment variables and unset those starting with 'e'
  for ($i = 1; $i -le 999; $i++) {
    $env_var_i = $scmpuff_env_char + $i
    if (Get-Item "env:$($env_var_i)" -ErrorAction Ignore) {
      Remove-Item "env:$env_var_i"
    } else {
      break
    }
  }
}
