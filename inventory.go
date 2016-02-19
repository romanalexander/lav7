package lav7

import "github.com/L7-MCPE/lav7/types"

// SimpleInventory is just a set of items, for containers or inventory holder entities.
type SimpleInventory []types.Item

func (si SimpleInventory) Get(pos uint16) types.Item {
	return si[pos]
}

// PlayerInventory is a inventory holder for players.
type PlayerInventory struct {
	*SimpleInventory
	Hotbars  []types.Item
	Hand     types.Item
	Creative bool
}
