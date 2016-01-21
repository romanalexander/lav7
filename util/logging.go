package util

import "fmt"

// PrintDbg will be true if the program is started with -d flag.
var PrintDbg = true

// Debug prints log if PrintDbg is true
func Debug(args ...interface{}) {
	if PrintDbg {
		fmt.Println(args...)
	}
}
