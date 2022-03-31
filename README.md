# wheel-of-fortune
Project for 1C

Wheel of fortune client-server game

# Install

## Requirements

- Go 1.18

## Build

```shell
$ make client
$ make server
```

# Run

Client and server will be placed in bin directory

## Server

```
Usage of ./bin/server:
  -port string
        Port for game_server listen (default "8888")
  -tries int
        Number of tries for wheel of fortune game (default 30)
  -words string
        Filename with word for wheel of fortune, separated by newline (default "/text.txt")

```

## Client

```
Usage of ./bin/client:
  -help
        Show usage
  -s string
        Address of game server (default "localhost:8888")
```