package raknet

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/L7-MCPE/lav7/util/buffer"
)

var serverID uint64
var blockList = make(map[string]time.Time)
var blockLock = new(sync.Mutex)
var GotBytes uint64

// Router handles packets from network, and manages sessions.
type Router struct {
	sessions      []Session
	conn          *net.UDPConn
	sendChan      chan Packet
	playerAdder   func(*net.UDPAddr) chan<- *buffer.Buffer
	playerRemover func(*net.UDPAddr) error
}

// CreateRouter create/opens new raknet router with given port.
func CreateRouter(playerAdder func(*net.UDPAddr) chan<- *buffer.Buffer,
	playerRemover func(*net.UDPAddr) error, port uint16) (r *Router, err error) {
	InitProtocol()
	Sessions = make(map[string]*Session)
	r = new(Router)
	serverID = uint64(rand.Int63())
	r.sessions = make([]Session, 0)
	r.sendChan = make(chan Packet, chanBufsize)
	r.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: int(port)})
	r.playerAdder = playerAdder
	r.playerRemover = playerRemover
	return
}

// Start makes router process network I/O operations.
func (r *Router) Start() {
	go r.sendAsync()
	go r.receivePacket()
}

func (r *Router) receivePacket() {
	defer r.conn.Close()
	for {
		var n int
		var addr *net.UDPAddr
		var err error
		recvbuf := make([]byte, 1024*1024)
		if n, addr, err = r.conn.ReadFromUDP(recvbuf); err != nil {
			fmt.Println("Error while reading packet:", err)
			continue
		} else if n > 0 {
			atomic.AddUint64(&GotBytes, uint64(n))
			buf := buffer.FromBytes(recvbuf[0:n])
			pk := Packet{
				Buffer:  buf,
				Address: addr,
			}
			if buf.Payload[0] == 0x01 { // Check if the packet is unconnected ping
				pingID := buf.ReadLong()
				buf := new(buffer.Buffer)
				buf.WriteByte(0x1c)
				buf.WriteLong(pingID)
				buf.WriteLong(serverID)
				buf.Write([]byte(RaknetMagic))
				buf.WriteString(GetServerString())
				pk := Packet{
					Buffer:  buf,
					Address: addr,
				}
				r.sendPacket(pk)
				continue
			}
			func() {
				blockLock.Lock()
				defer blockLock.Unlock()
				if blockList[addr.String()].After(time.Now()) {
					r.conn.WriteToUDP([]byte("\x80\x00\x00\x00\x00\x00\x08\x15"), pk.Address)
				} else {
					delete(blockList, addr.String())
					GetSession(addr, r.sendChan, r.playerAdder, r.playerRemover).ReceivedChan <- pk
				}
			}()
		}
	}
}

func (r *Router) sendAsync() {
	for pk := range r.sendChan {
		r.sendPacket(pk)
	}
}

func (r *Router) sendPacket(pk Packet) {
	r.conn.WriteToUDP(pk.Buffer.Payload, pk.Address)
}
