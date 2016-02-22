package lav7

import "github.com/L7-MCPE/lav7/types"

// SimpleInventory is just a set of items, for containers or inventory holder entities.
type SimpleInventory []types.Item

// PlayerInventory is a inventory holder for players.
type PlayerInventory struct {
	*SimpleInventory
	Hotbars []types.Item
	Hand    types.Item
	Holder  *Player
}

// Init initializes the inventory.
func (pi PlayerInventory) Init() {

}
