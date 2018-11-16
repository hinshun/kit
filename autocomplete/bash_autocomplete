#! /bin/bash

_kit() {
    local cur opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( ${COMP_WORDS[0]} --autocomplete bash ${COMP_WORDS[@]:1:$COMP_CWORD} )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}

complete -F _kit kit
