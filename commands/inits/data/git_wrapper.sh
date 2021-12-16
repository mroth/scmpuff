# Remove any existing git alias or function
unalias git > /dev/null 2>&1
unset -f git > /dev/null 2>&1

# Use the full path to git to avoid infinite loop with git function
SCMPUFF_GIT_CMD=${SCMPUFF_GIT_CMD:-"$(\which git)"}
export SCMPUFF_GIT_CMD

function git() {
  case $1 in
    commit|blame|log|rebase|merge)
      scmpuff exec -- "$SCMPUFF_GIT_CMD" "$@";;
    checkout|diff|rm|reset|restore)
      scmpuff exec --relative -- "$SCMPUFF_GIT_CMD" "$@";;
    add)
      scmpuff exec -- "$SCMPUFF_GIT_CMD" "$@"
      scmpuff_status;;
    *)
      "$SCMPUFF_GIT_CMD" "$@";;
  esac
}
