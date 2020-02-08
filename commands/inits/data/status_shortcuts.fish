# Based on https://github.com/arbelt/fish-plugin-scmpuff,
# with fish3 fix https://github.com/arbelt/fish-plugin-scmpuff/pull/3
function scmpuff_status
    scmpuff_clear_vars
    set -lx scmpuff_env_char "e"
    set -l cmd_output (/usr/bin/env scmpuff status --filelist $argv ^/dev/null)
    set -l es "$status"

    if test $es -ne 0
        git status
        return $status
    end

    set -l files (string split \t $cmd_output[1])
    if test (count $files) -gt 0
        for e in (seq (count $files))
            set -gx "$scmpuff_env_char""$e" "$files[$e]"
        end
    end

    for line in $cmd_output[2..-1]
        echo $line
    end
end

function scmpuff_clear_vars
    set -l scmpuff_env_char "e"
    set -l scmpuff_env_vars (set -x | awk '{print $1}' | grep -E '^'$scmpuff_env_char'\d+')

    for v in $scmpuff_env_vars
        set -e $v
    end
end
