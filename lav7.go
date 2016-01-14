package lav7

// ServerName contains human readable server name
var ServerName = "Lav7 - lightweight MCPE server"

// Players is a map containing Player structs.
var Players = make(map[string]interface{})

// MaxPlayers is count of maximum available players
var MaxPlayers = 20

const (
	// MinecraftProtocol is a mojang network protocol version.
	MinecraftProtocol = 38
	// MinecraftVersion is a human readable minecraft version.
	MinecraftVersion = "0.13.1"
)
