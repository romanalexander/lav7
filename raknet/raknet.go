// Package raknet implements MCPE network protocol, both internal Raknet and Mojang-implemented protocols.
package raknet

import "strconv"

const (
	// RaknetMagic is a magic bytes for internal Raknet protocol.
	RaknetMagic = "\x00\xff\xff\x00\xfe\xfe\xfe\xfe\xfd\xfd\xfd\xfd\x12\x34\x56\x78"
	// RaknetProtocol is a internal Raknet protocol version.
	RaknetProtocol = 6
	// RaknetVersion is a version of this library.
	RaknetVersion = "1.0.0"
	// MinecraftProtocol is a mojang network protocol version.
	MinecraftProtocol = 43
	// MinecraftVersion is a human readable minecraft version.
	MinecraftVersion = "0.14.0.4"
)

// Players is a reference of lav7.players; it is a dirty trick to avoid import cycle.
var Players map[string]interface{}

// ServerName contains human readable server name
var ServerName = "Lav7 - lightweight MCPE server"

// OnlinePlayers is count of online players
var OnlinePlayers = 0

// MaxPlayers is count of maximum available players
var MaxPlayers = 20

// GetServerString returns server status message for unconnected pong
func GetServerString() string {
	return "MCPE;" + ServerName + ";" +
		strconv.Itoa(MinecraftProtocol) + ";" +
		MinecraftVersion + ";" +
		strconv.Itoa(OnlinePlayers) + ";" +
		strconv.Itoa(MaxPlayers)
}
