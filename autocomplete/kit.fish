function __kit
  set -l cmd (commandline -poc)
  echo (eval $cmd[1] --autocomplete fish $cmd[2..-1]) | xargs printf '%s\n'
end

complete -f -c kit -a '(__kit)'
