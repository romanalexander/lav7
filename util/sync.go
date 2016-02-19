package util

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var MutexDebug bool

// NewMutex returns new Locker; sync.Mutex or MutexTracer.
// new(sync.Mutex) should be changed to this for supporting mutex debugging.
func NewMutex() Locker {
	if MutexDebug {
		return NewMutexTracer()
	}
	return new(sync.Mutex)
}

// NewRWMutex returns new RWLocker.
// See NewMutex comment for more informations.
func NewRWMutex() RWLocker {
	if MutexDebug {
		return NewRWMutexTracer()
	}
	return new(sync.RWMutex)
}

// Locker is a alias for sync.Locker; this is just for a consistency with RWLocker interface.
type Locker sync.Locker

// MutexTracer is a simple wrapper for sync.Mutex to debug lock issues
type MutexTracer struct {
	*sync.Mutex
	lastLock  time.Time
	lockTrace string
}

// NewMutexTracer returns new MutexTracer with valid initial value.
func NewMutexTracer() *MutexTracer {
	m := new(MutexTracer)
	m.Init()
	return m
}

// Init initializes the MutexTracer struct.
func (m *MutexTracer) Init() {
	m.Mutex = new(sync.Mutex)
}

// Lock locks the hidden mutex for tracer, and records when the lock grabbed.
func (m *MutexTracer) Lock() {
	// log.Print("Aquiring lock")
	m.Mutex.Lock()
	m.lastLock = time.Now()
	m.lockTrace = GetTrace()
}

// Unlock unlocks the hidden mutex for tracer, and calculates how long the lock grabbed.
func (m *MutexTracer) Unlock() {
	grabtime := time.Since(m.lastLock).Seconds()
	if grabtime > 1 {
		log.Println("!!! Lock grabbed for > 1 seconds?!?! !!!")
		fmt.Print(m.lockTrace)
		fmt.Print(GetTrace())
	}
	m.Mutex.Unlock()
}

// RWLocker is a interface for holding RWMutex or RWMutexTracer.
type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

// RWMutexTracer is a simple wrapper for sync.RWMutex to debug lock issues
type RWMutexTracer struct {
	*sync.RWMutex
	lastLock  time.Time
	lockTrace string
}

// NewRWMutexTracer returns new RWMutexTracer with valid initial value.
func NewRWMutexTracer() *RWMutexTracer {
	m := new(RWMutexTracer)
	m.Init()
	return m
}

// Init initializes the RWMutexTracer struct.
func (m *RWMutexTracer) Init() {
	m.RWMutex = new(sync.RWMutex)
}

// Lock locks the hidden mutex for tracer, and records when the lock grabbed.
func (m *RWMutexTracer) Lock() {
	// log.Print("Aquiring lock")
	m.RWMutex.Lock()
	m.lastLock = time.Now()
	m.lockTrace = GetTrace()
}

// Unlock unlocks the hidden mutex for tracer, and calculates how long the lock grabbed.
func (m *RWMutexTracer) Unlock() {
	grabtime := time.Since(m.lastLock).Seconds()
	if grabtime > 1 {
		log.Println("!!! Lock grabbed for > 1 seconds?!?! !!!")
		fmt.Print(m.lockTrace)
		fmt.Print(GetTrace())
	}
	m.RWMutex.Unlock()
}
