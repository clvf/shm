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
Open Files by User
=====================

lsof -u $username
lsof -i -u $username

Exclude Specific User (with '^')
--------------------------------
lsof -i -u^root

Process On a Specific Port
============================

lsof -i TCP:9999

Open Files For Port Ranges
----------------------------

lsof -i TCP:9000-10000

Kill All Processes of User
===========================

kill $(lsof -t -u $username)

$ shm add web/curl

$ shm
home/${USER}/.shm
├── network
│   └── lsof
└── web
    └── curl
```
