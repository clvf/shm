# SHow Me (shm)

SHow Me (`shm`) is a small unix command line utility to manage recipes and
knowledge base articles inspired by Pass, "the standard unix password
manager".

```sh
$ shm add network/lsof

$ shm
    home/${USER}/.shm
    └── network
        └── lsof

$ shm network/lsof

    Open Network Files by User
    ===========================

    lsof -i -u $username

    Exclude Specific User (with '^')
    --------------------------------
    lsof -i -u^root

$ shm add web/curl

$ shm
    home/${USER}/.shm
    ├── network
    │   └── lsof
    └── web
        └── curl
```

#### Shell Completion

If you wanted to have shell completion then download the appropriate completion
file from the [urfave/cli](https://github.com/urfave/cli/tree/main/autocomplete "urfave/cli v2")
repo and copy it to `~/local/share/bash-completion/completions/shm`.

Eg.:

```sh
# alternatively you can clone the repo https://github.com/urfave/cli.git
$ wget https://raw.githubusercontent.com/urfave/cli/refs/heads/main/autocomplete/bash_autocomplete

$ mkdir -p ~/local/share/bash-completion/completions
$ mv bash_autocomplete ~/local/share/bash-completion/completions/shm
```

***

Mind you the whole Go code is equivalent with this simple and short shell
script. In fact this is the code tagged by `LAST_SHELL_VERSION` in the
git history.

```sh
#!/bin/env bash

set -e

: ${SHM_STORE:=~/.shm}
: ${EDITOR:=vi}

if which tree >/dev/null 2>&1;
then
    LIST="tree --noreport --prune --"
else
    LIST=find
fi

if which pygmentize >/dev/null 2>&1 \
   && pygmentize -l md </dev/null >/dev/null 2>&1;
then
    CAT="pygmentize -l md"
else
    CAT=cat
fi

if [[ -z "${1}" ]];
then
    mkdir -p ${SHM_STORE}
    ${LIST} ${SHM_STORE}
    echo
    exit
fi

case ${1} in
    "add" | "edit" ) 
        mkdir -p $(dirname ${SHM_STORE}/${2})
        ${EDITOR} ${SHM_STORE}/${2}
        ;;
    "search" | "find" )
        find ${SHM_STORE} -name "*${2}*"
        ;;
    "rm" ) 
        rm ${SHM_STORE}/${2}
        ;;
    "mv" )
        mv ${SHM_STORE}/${2} ${SHM_STORE}/${3}
        ;;
    "cp" )
        cp ${SHM_STORE}/${2} ${SHM_STORE}/${3}
        ;;
    *) ${CAT} ${SHM_STORE}/${1};;
esac
```
