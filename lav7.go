package lav7

import "github.com/L7-MCPE/lav7/level"

const (
	// MinecraftProtocol is a mojang network protocol version.
	MinecraftProtocol = 45
	// MinecraftVersion is a human readable minecraft version.
	MinecraftVersion = "0.14.0.7"
)

// ServerName contains human readable server name
var ServerName = "Lav7 - lightweight MCPE server"

// Players is a map containing Player structs.
var Players = make(map[string]*Player)

// MaxPlayers is count of maximum available players
var MaxPlayers = 20

var LastEntityID uint64

var levels = map[string]*level.Level{
	defaultLvl: new(level.Level),
}
var defaultLvl = "default"
