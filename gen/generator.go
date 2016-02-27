// Package gen provides MCPE level generator for creating worlds.
package gen

import (
	"reflect"
	"strings"

	"github.com/L7-MCPE/lav7/types"
)

// Generator is an interface for MCPE map generator.
type Generator interface {
	Init()
	Gen(int32, int32) *types.Chunk
}

// Seedused is an interface for generators which uses seeds for randomness.
type Seedused interface {
	Generator
	Seed() *int64
}

var levelGenerators = map[string]Generator{}

// RegisterGenerator adds level format Generator for server.
// Generator name must end with "Generator".
func RegisterGenerator(g Generator) {
	typname := reflect.TypeOf(g).String()
	if typname[len(typname)-9:] != "Generator" {
		panic("Invalid generator name: " + typname)
	}
	typsl := strings.Split(typname, ".")
	name := strings.ToLower(typsl[len(typsl)-1])
	name = name[:len(name)-9]
	if _, ok := levelGenerators[name]; !ok {
		levelGenerators[name] = g
	}
}

// GetGenerator finds the Generator with given name.
// If it doesn't present, returns nil.
func GetGenerator(name string) Generator {
	if pv, ok := levelGenerators[name]; ok {
		return pv
	}
	return nil
}
