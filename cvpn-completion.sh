#!/usr/bin/bash

_cvpn() {
    local prev
    _get_comp_words_by_ref -n : prev

    if [ ${prev} = "list" ] || [ ${prev} = "ls" ]; then
        __cvpn_list
    fi  
}

__cvpn_list() {    
    local cur prev
    _get_comp_words_by_ref -n : cur prev

    local IFS=$'\n'
    local separator='/'
    local count_of_slash=$(awk -F"${separator}" '{print NF-1}' <<< "${cur}")
 
    local opts=$(
        while read line
        do
            echo "'${line}'"
        done < index.txt
    );

    local new_opts=$(
        while read line
        do
            if [ "$(awk -F"${separator}" '{print NF-1}' <<< "${line}")" = "${count_of_slash}" ]; then
                echo "${line}"
            fi
        done <<< "${opts}"
    );

    if [ "${#new_opts}" != "0" ]; then
        opts="${new_opts}"
    fi


    COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") ); IFS=${defaultIFS}
}

complete -F _cvpn cvpn