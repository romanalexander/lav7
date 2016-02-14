package command

import (
	"bufio"
	"os"
	"strings"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/util"
)

// HandleCommand handles command input from stdin.
func HandleCommand() {
	for {
		r := bufio.NewReader(os.Stdin)
		text, _ := r.ReadString('\n')
		texts := strings.Split(strings.Replace(text[:len(text)-1], "\r", "", -1), " ")
		switch texts[0] {
		case "stop", "exit":
			lav7.Stop(strings.Join(texts[1:], " "))
		case "spawn":
			lav7.SpawnPlayer(&lav7.Player{
				Username: "Test",
				EntityID: 99,
				Position: util.Vector3{0, 65, 0},
			})
		case "move":
			lav7.BroadcastPacket(&proto.MovePlayer{
				EntityID: 199,
				X:        1,
				Y:        64,
				Z:        3,
			})
		}
	}
}
