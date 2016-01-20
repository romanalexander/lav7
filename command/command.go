package command

import (
	"bufio"
	"os"
	"strings"

	"github.com/L7-MCPE/lav7"
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
		}
	}
}
