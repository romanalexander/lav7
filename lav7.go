package lav7

import (
	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/level/format/dummy"
)

const (
	// MinecraftProtocol is a mojang network protocol version.
	MinecraftProtocol = 38
	// MinecraftVersion is a human readable minecraft version.
	MinecraftVersion = "0.13.1"
)

// ServerName contains human readable server name
var ServerName = "Lav7 - lightweight MCPE server"

// Players is a map containing Player structs.
var Players = make(map[string]interface{})

// MaxPlayers is count of maximum available players
var MaxPlayers = 20

var lastEntityID uint64

var levels = map[string]level.Level{
	defaultLvl: new(dummy.Level),
}
var defaultLvl = "default"
