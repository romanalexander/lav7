`Copyright (C) 2016 cr0sh(Nam JH.)`

`This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under certain conditions; read LICENSE.txt in root of the repository.`

# lav7
[![GoDoc](https://godoc.org/github.com/L7-MCPE/lav7?status.svg)](https://godoc.org/github.com/L7-MCPE/lav7)
[![Go Report Card](https://goreportcard.com/badge/github.com/L7-MCPE/lav7)](https://goreportcard.com/report/github.com/L7-MCPE/lav7)

lav7 aims to be Lightweight, rapid, concurrent Minecraft:PE server.
The project's main goal is to provide gameplay features close to vanilla Minecraft:PE server, and handle 15~20 players on ARM11 CPU server, such as Raspberry Pi B+ Model.

## Current project status
lav7 needs much more feature implementations, like level generator, or player movements, etc. A short-term goal for this project is to implement functions on the same level as PocketMine-MP 1.3.1, until March.

## Contributions
Pull requests are always welcome, but please check these before writing pull request:
 - **Format your codes.** Unifying coding styles are important to collaborate. Please follow suggestions from `gofmt`, `golint`, `go vet` if you can.
  - Exception: You can omit documentation comment of exported items if it could be useless.

## Installation
 - Requirements: Latest Go installation
  - Windows: [Download installation package from site](https://golang.org/dl/)
  - Linux: Use `apt-get` or `yum`
 - Add GOPATH and set PATH to GOPATH/bin
 - To install or update lav7, run `go get -u github.com/L7-MCPE/lav7/l7start && go install github.com/L7-MCPE/lav7/l7start`
 - To run lav7, run `l7start`
