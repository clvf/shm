# SHow Me (shm)

SHow Me (`shm`) is a small unix command line utility to manage recipes and
knowledge base articles inspired by Pass, "the standard unix password
manager".

```
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
