# Based on https://github.com/arbelt/fish-plugin-scmpuff,
# with scmpuff-exec support (https://github.com/mroth/scmpuff/pull/49)
functions -e git

set -q SCMPUFF_GIT_CMD; or set -x SCMPUFF_GIT_CMD (which git)

function git
    if test (count $argv) -eq 0
        eval $SCMPUFF_GIT_CMD
        set -l s $status
        return $s
    end

    switch $argv[1]
    case commit blame log rebase merge
        scmpuff exec -- "$SCMPUFF_GIT_CMD" $argv
    case checkout diff rm reset restore
        scmpuff exec --relative -- "$SCMPUFF_GIT_CMD" $argv
    case add
        scmpuff exec -- "$SCMPUFF_GIT_CMD" $argv
        scmpuff_status
    case '*'
        eval command "$SCMPUFF_GIT_CMD" (string escape -- $argv)
    end
end
