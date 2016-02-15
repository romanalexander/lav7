package util

import (
	"log"
	"runtime/debug"
	"sync"
)

// MutexTracer is a simple wrapper for sync.Mutex to debug lock issues
type MutexTracer struct {
	lock *sync.Mutex
}

// NewMutexTracer returns new MutexTracer with valid initial value.
func NewMutexTracer() *MutexTracer {
	m := new(MutexTracer)
	m.Init()
	return m
}

func (m *MutexTracer) Init() {
	m.lock = new(sync.Mutex)
}

func (m MutexTracer) Lock() {
	log.Print("Aquiring lock")
	m.lock.Lock()
	log.Print("Lock grabbed")
	debug.PrintStack()
}

func (m MutexTracer) Unlock() {
	log.Print("Releasing lock")
	debug.PrintStack()
	m.lock.Unlock()
}
