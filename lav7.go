// Package lav7 is not only a lightweight Minecraft:PE server, but provides Minecraft:PE protocol/gameplay mechanics.
package lav7

import "sync"

const (
	// Version is a version of this server.
	Version = "0.1.0 in-dev"
	// ServerName contains human readable server name
	ServerName = "Lav7 - lightweight MCPE server"
	// MaxPlayers is count of maximum available players
	MaxPlayers = 20
)

// GitVersion is a git commit hash for this project.
// You should set this with -ldflags "-X github.com/L7-MCPE/lav7.GitVersion="
var GitCommit = "unknown"

// BuildTime is a timestamp when the program is built.
var BuildTime = "unknown"

// Players is a map containing Player structs.
var Players = make(map[string]*Player)

var iteratorLock = new(sync.Mutex)

var LastEntityID uint64

var levels = map[string]*Level{
	defaultLvl: {Name: "dummy"},
}
var defaultLvl = "default"
