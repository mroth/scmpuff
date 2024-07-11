Remove-Item alias:\git
Remove-Item function:\git

$SCMPUFF_GIT_CMD = Get-Command git | Select-Object -ExpandProperty Definition

function git {
  switch -regex -casesensitive($args[0]) {
    "^(commit|blame|log|rebase|merge)$" {
      & scmpuff expand -- $SCMPUFF_GIT_CMD $args
    }
    "^(checkout|diff|rm|reset)$" {
      & scmpuff expand --relative -- $SCMPUFF_GIT_CMD $args
    }
    "^add$" {
      & scmpuff expand -- $SCMPUFF_GIT_CMD $args
      scmpuff status
    }
    default {
      & $SCMPUFF_GIT_CMD $args
    }
  }
}
