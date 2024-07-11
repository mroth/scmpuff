$SCMPUFF_GIT_CMD = Get-Command git | Select-Object -ExpandProperty Definition

function git {
  switch -regex -casesensitive($args[0]) {
    "^(commit|blame|log|rebase|merge)$" {
      & scmpuff exec -- $SCMPUFF_GIT_CMD $args
    }
    "^(checkout|diff|rm|reset)$" {
      & scmpuff exec --relative -- $SCMPUFF_GIT_CMD $args
    }
    "^add$" {
      & scmpuff exec -- $SCMPUFF_GIT_CMD $args
      scmpuff_status
    }
    default {
      & $SCMPUFF_GIT_CMD $args
    }
  }
}
