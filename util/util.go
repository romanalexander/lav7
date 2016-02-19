// Package util provides some functions used widely.
package util

import (
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

// FmtElement contains type of parsed field and numbers
type FmtElement struct {
	T byte
	C int
}

// ParseBatchFmt parses format string for batch read/write operations
func ParseBatchFmt(format string) (pe []FmtElement) {
	var stage byte
	var num string
	for i := 0; i <= len(format); i++ {
		if i == len(format) {
			if stage != 0 {
				if len(num) > 0 {
					n, _ := strconv.Atoi(num)
					pe = append(pe, FmtElement{stage, n})
					continue
				}
				pe = append(pe, FmtElement{stage, 1})
			}
		} else if _, err := strconv.Atoi(string(format[i])); err == nil { //Number
			num += string(format[i])
		} else {
			if len(num) > 0 { //character+numbers
				n, _ := strconv.Atoi(num)
				pe = append(pe, FmtElement{stage, n})
				num = ""
			} else { //Normal characters
				if stage != 0 {
					pe = append(pe, FmtElement{stage, 1})
					stage = 0
				}
			}
			stage = format[i]
		}
	}
	return
}

// GetSortedKeys will return a sorted slice of 'uint' keys from given map.
func GetSortedKeys(m interface{}) []int {
	mm := reflect.ValueOf(m)
	keys := make([]int, len(mm.MapKeys()))
	i := 0
	for _, k := range mm.MapKeys() {
		keys[i] = int(k.Uint())
		i++
	}
	sort.Ints(keys)
	return keys
}

// GetTrace returns stack trace for all goroutines.
func GetTrace() string {
	var b [1024 * 1024 * 16]byte
	n := runtime.Stack(b[:], true)
	return string(b[:n])
}

// Suspend is a dumb funciton which will make program run forever - hangs current goroutine.
func Suspend() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	wg.Wait()
}
