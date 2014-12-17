git () {
case $1 in
  (commit|blame|add|log|rebase|merge) scmpuff expand "$_git_cmd" "$@" ;;
  (checkout|diff|rm|reset) scmpuff expand --relative "$_git_cmd" "$@" ;;
  (branch) _scmb_git_branch_shortcuts "${@:2}" ;;
  (*) "$_git_cmd" "$@" ;;
esac
}
