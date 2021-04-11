#compdef cvpn

# Zsh completion for cvpn
#
# INSTALLATION
#
# Make sure autocompletion is enabled in your shell, typically
# by adding this to your .zshrc:
#
#     autoload -U compinit && compinit
#
# Then copy this file somewhere (e.g. ~/.zsh/completions/_cvpn) and put the
# following in your .zshrc:
#
#     fpath=(${HOME}/.zsh/completions $fpath)
#

compdef _cvpn cvpn

function _cvpn() {
  _arguments -C \
    '(- *)'{-h,--help}'[Show help]' \
    '1: :__cvpn_subcmds' \
    '*::subcmd_arg:->subcmd_args'

  case $state in
    subcmd_args)
      case $line[1] in
        completion)
          if [ ${#line[@]} -eq 2 ]; then  # if (num of elements of line) == 2
            _values 'shells' 'bash' 'zsh'
          fi
          ;;
        download)
          _arguments \
            '(-o --output)'{-o,--output}'[Destination path of downloaded file]:file:_files' \
            '(-v --volume)'{-v,--volume}'[Volume ID on FS (Default: fsshare)]:volumeID:__volume_id' \
            '(- *)'{-h,--help}'[Show help]' \
            ':download target path on FS:__cached_fs_pathes'
          ;;
        find)
          _arguments \
            '(-r --recursive)'{-r,--recursive}'[Search recursively]' \
            '(-v --volume)'{-v,--volume}'[Volume ID on FS (Default: fsshare)]:volumeID:__volume_id' \
            '(--name)'--name'[Regex pattern of path''s name (last path component)]:regex' \
            '(--path)'--path'[Regex pattern of path]:regex' \
            '(- *)'{-h,--help}'[Show help]' \
            ':starting directory on FS:__cached_fs_pathes'
          ;;
        list|ls)
          _arguments \
            '(-v --volume)'{-v,--volume}'[Volume ID on FS (Default: fsshare)]:volumeID:__volume_id' \
            '(--json)'--json'[Print with JSON format]' \
            '(--path)'--path'[Also shows path]' \
            '(- *)'{-h,--help}'[Show help]' \
            ':directory path on FS:__cached_fs_pathes'
          ;;
        upload)
          _arguments \
            '(-v --volume)'{-v,--volume}'[Volume ID on FS (Default: fsshare)]:volumeID:__volume_id' \
            '(- *)'{-h,--help}'[Show help]' \
            ':local file:_files' \
            ':destination path on FS:__cached_fs_pathes'
          ;;
      esac
      ;;
  esac
}

function __cvpn_subcmds() {
  local cvpn_subcmds=(
    completion':Print the competion script'
    download':Download a file on the FS'
    find':Find files and dirs with using regex'
    list':List up directory contents on the FS'
    login':Login to FS (Required only for first time)'
    upload':Upload a local file into FS:file:_files'
    help':Show help'
  )
  _describe 'command' cvpn_subcmds
}

function __volume_id() {
  _values 'volumeID' 'fsshare' 'fs'
}

function __cached_fs_pathes() {
  readonly cachefile="%s"  # <- replaced by Sprintf() in Go program
  local cache_entries=()

  # Return if cache file is not readable
  [[ -r "$cachefile" ]] || return

  # Collect each line into array 
  for l in $(cat "$cachefile"); do
    cache_entries=($cache_entries "$l")
  done

  compadd $cache_entries
}
