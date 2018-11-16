#compdef kit

_kit() {
  eval $(${words[@]:0:1} --autocomplete zsh ${words[@]:1:#words[@]-1})
}

_kit
