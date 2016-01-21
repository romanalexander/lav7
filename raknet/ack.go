package raknet

import (
	"sort"

	"github.com/L7-MCPE/lav7/util/buffer"
)

type ackTable []uint32

func (t ackTable) Len() int           { return len(t) }
func (t ackTable) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ackTable) Less(i, j int) bool { return t[i] < t[j] }

// EncodeAck packs packet sequence numbers into Raknet acknowledgement format.
func EncodeAck(t ackTable) (b *buffer.Buffer) {
	b = buffer.FromBytes(make([]byte, 0))
	sort.Sort(t)
	packets := t.Len()
	records := uint16(0)
	if packets > 0 {
		pointer := 1
		start, last := t[0], t[0]
		for pointer < packets {
			current := t[pointer]
			pointer++
			diff := current - last
			if diff == 1 {
				last = current
			} else if diff > 1 {
				if start == last {
					b.WriteByte(1)
					b.WriteLTriad(start)
					last = current
					start = last
				} else {
					b.WriteByte(0)
					b.WriteLTriad(start)
					b.WriteLTriad(last)
					last = current
					start = last
				}
				records++
			}
		}
		if start == last {
			b.WriteByte(1)
			b.WriteLTriad(start)
		} else {
			b.WriteByte(0)
			b.WriteLTriad(start)
			b.WriteLTriad(last)
		}
		records++
	}
	tmp := buffer.FromBytes(make([]byte, 0))
	tmp.WriteShort(records)
	tmp.Append(b)
	b = tmp
	return
}

// DecodeAck unpacks packet sequence numbers from given buffer.
func DecodeAck(b *buffer.Buffer) (t []uint32) {
	var records uint16
	records = b.ReadShort()
	count := 0
	for i := 0; uint16(i) < records && b.Require(1) && count < 4096; i++ {
		if f := b.ReadByte(); f == 0 {
			start := b.ReadLTriad()
			last := b.ReadLTriad()
			if (last - start) > 512 {
				last = start + 512
			}
			for c := start; c <= last; c++ {
				t = append(t, c)
				count++
			}
		} else {
			p := b.ReadLTriad()
			t = append(t, p)
			count++
		}
	}
	return
}
