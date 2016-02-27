// Package config stores global server config for lov7.
// Each properties file line should be key=value(value could be empty).
// You must not put space or other characters between keys/values and =.
// Comment lines sould start with #.
//
// Changing config values at runtime is not recommended because of data race.
package config

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

// DefaultConfig is a default config string.
var DefaultConfig = `# Default lav7 properties
server-port=19132
server-name=lav7 - lightweight MCPE server
max-players=20
generator-name=flat
generator-args=
level-format=vilan
chunk-radius=6
`

// Port is a port number of the server.
var Port uint16

// ServerName is a server name displayed on lists.
var ServerName string

// MaxPlayers is a maximum player available on server.
var MaxPlayers int32

// Generator is a name of world generator.
var Generator string

// GeneratorArgs is additional arguments handled by (Generator).Init()
var GeneratorArgs string

// Format is a name of level format provider.
var Format string

// ChunkRadius is a default chunk send radius for client.
var ChunkRadius int32

// Parse parses the config with given reader interface.
func Parse(rd io.Reader) {
	scanner := bufio.NewScanner(rd)
	cfg := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0:1] == "#" {
			continue
		}
		split := strings.SplitN(line, "=", 2)
		if len(split) < 2 {
			split = append(split, "")
		}
		cfg[split[0]] = split[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Error while reading config:", err)
	}

	port, err := strconv.Atoi(getString(cfg, "server-port", "19132"))
	if err != nil {
		log.Fatalln("Invalid server port")
	}
	Port = uint16(port)

	ServerName = getString(cfg, "server-name", "lav7 - lightweight MCPE server")
	m, err := strconv.Atoi(getString(cfg, "max-players", "20"))
	if err != nil {
		log.Fatalln("Invalid max players")
	}
	MaxPlayers = int32(m)

	Generator = getString(cfg, "generator-name", "flat")
	GeneratorArgs = getString(cfg, "generator-args", "")
	Format = getString(cfg, "level-format", "vilan")

	chunkRadius, err := strconv.Atoi(getString(cfg, "chunk-radius", "6"))
	if err != nil {
		log.Fatalln("Invalid chunk radius")
	}
	ChunkRadius = int32(chunkRadius)
}

func getString(m map[string]string, key string, def string) string {
	val, ok := m[key]

	if !ok {
		m[key] = def
		return def
	}
	return val
}
