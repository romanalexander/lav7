// Package lav7 is not only a lightweight Minecraft:PE server, but provides Minecraft:PE protocol/gameplay mechanics.
package lav7

import (
	"reflect"
	"strings"
	"sync"

	"github.com/L7-MCPE/lav7/format"
)

const (
	// Version is a version of this server.
	Version = "1.1.0 alpha-dev"
	// ServerName contains human readable server name
	ServerName = "lav7 - lightweight MCPE server"
	// MaxPlayers is count of maximum available players
	MaxPlayers = 20
)

// GitCommit is a git commit hash for this project.
// You should set this with -ldflags "-X github.com/L7-MCPE/lav7.GitVersion="
var GitCommit = "unknown"

// BuildTime is a timestamp when the program is built.
var BuildTime = "unknown"

// Players is a map containing Player structs.
var Players = make(map[string]*Player)

var iteratorLock = new(sync.Mutex)

var lastEntityID uint64

var levels = map[string]*Level{
	defaultLvl: {Name: "dummy"},
}

var levelProviders = map[string]format.Provider{}

var defaultLvl = "default"

// RegisterProvider adds level format provider for server.
func RegisterProvider(provider format.Provider) {
	typname := reflect.TypeOf(provider)
	typsl := strings.Split(typname.String(), ".")
	name := strings.ToLower(typsl[len(typsl)-1])
	if _, ok := levelProviders[name]; !ok {
		levelProviders[name] = provider
	}
}

// GetProvider finds the provider with given name.
// If it doesn't present, returns nil.
func GetProvider(name string) format.Provider {
	if pv, ok := levelProviders[name]; ok {
		return pv
	}
	return nil
}
