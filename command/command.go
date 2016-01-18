package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HandleCommand handles command input from stdin.
func HandleCommand() {
	for {
		r := bufio.NewReader(os.Stdin)
		text, _ := r.ReadString('\n')
		texts := strings.Split(strings.Replace(text[:len(text)-1], "\r", "", -1), " ")
		fmt.Println(texts)
	}
}
