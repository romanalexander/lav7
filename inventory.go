package lav7

import (
	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/types"
)

// Inventory is just a set of items, for containers or inventory holder entities.
type Inventory []types.Item

// PlayerInventory is a inventory holder for players.
type PlayerInventory struct {
	*Inventory
	Hotbars []types.Item
	Hand    types.Item
	Holder  *Player
}

// Init initializes the inventory.
func (pi *PlayerInventory) Init() {
	pi.Hotbars = make([]types.Item, 8)
	if true { // No survival inventory now
		inv := make(Inventory, len(types.CreativeItems))
		copy(inv, types.CreativeItems)
		pi.Inventory = &inv
		pi.Holder.SendCompressed(&proto.ContainerSetContent{
			WindowID: proto.CreativeWindow,
			Slots:    inv,
		})
	}
}
