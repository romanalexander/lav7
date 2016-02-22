package lav7

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/types"
)

var creativeInvCache *bytes.Buffer

func init() {
	inv := make(Inventory, len(types.CreativeItems))
	copy(inv, types.CreativeItems)
	creativeInvCache = (&proto.ContainerSetContent{
		WindowID: proto.CreativeWindow,
		Slots:    inv,
	}).Write()
	fmt.Print(hex.Dump(creativeInvCache.Bytes()))
}

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
		pi.Holder.Send(bytes.NewBufferString(creativeInvCache.String()))
	}
}
