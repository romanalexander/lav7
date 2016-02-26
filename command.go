package lav7

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"

	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/vector"
)

const skin = "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"

// HandleCommand handles command input from stdin.
func HandleCommand() {
	for {
		r := bufio.NewReader(os.Stdin)
		text, _ := r.ReadString('\n')
		texts := strings.Split(strings.Replace(text[:len(text)-1], "\r", "", -1), " ")
		switch texts[0] {
		case "stop", "exit":
			go Stop(strings.Join(texts[1:], " "))
		case "spawn":
			b, _ := base64.StdEncoding.DecodeString(skin)
			b, _ = util.DecodeDeflate(b)
			BroadcastPacket(&proto.PlayerList{
				Type: proto.PlayerListAdd,
				PlayerEntries: []proto.PlayerListEntry{
					{
						RawUUID:  [16]byte{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 5, 6},
						EntityID: 99,
						Username: "Test",
						Skinname: "Festive_FestiveSweaterSteve",
						Skin:     b,
					},
				},
			})
			SpawnPlayer(&Player{
				UUID:     [16]byte{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 5, 6},
				Username: "Test",
				EntityID: 99,
				Position: vector.Vector3{X: 0, Y: 70, Z: 0},
			})

		case "move":
			BroadcastPacket(&proto.MovePlayer{
				EntityID: 199,
				X:        1,
				Y:        64,
				Z:        3,
			})
		case "block":
			br := make([]proto.BlockRecord, 20)
			for i := 0; i < 20; i++ {
				br[i] = proto.BlockRecord{
					X:     0,
					Y:     byte(i) + 55,
					Z:     0,
					Block: types.Block{ID: 4},
				}
			}
			BroadcastPacket(&proto.UpdateBlock{
				BlockRecords: br,
			})
		case "trace":
			fmt.Print(util.GetTrace())
		case "gc":
			Message("[system] Cleaning server memory...")
			c := GetDefaultLevel().Clean()
			runtime.GC()
			debug.FreeOSMemory()
			Message(fmt.Sprintf("[system] Done. %d chunks saved/unloaded.", c))
		case "netbytes":
			bs := atomic.LoadUint64(&raknet.GotBytes)
			log.Printf("%dKBs", bs>>10)
		case "dump":
			f, _ := os.Create("heapdump")
			debug.WriteHeapDump(f.Fd())
			log.Printf("Done")
		case "sendchunk":
			BroadcastCallback(PlayerCallback{
				Call: func(p *Player, args interface{}) {
					p.SendNearChunk(nil)
				},
			})
		default:
			log.Println("?")
		}
	}
}
