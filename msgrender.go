package g53

import (
	"log"
	"os"

	"github.com/mistletoeChao/g53/util"
)

type offsetItem struct {
	hash uint32
	pos  uint16
	l    uint16
}

type nameComparator struct {
	buffer        *util.OutputBuffer
	nameBuf       *util.InputBuffer
	hash          uint32
	caseSensitive bool
}

var logger *log.Logger = log.New(os.Stdout, "[debug]", 0)

func (c *nameComparator) compare(item *offsetItem) bool {
	if item.hash != c.hash || item.l != uint16(c.nameBuf.Len()) {
		return false
	}

	itemPos := item.pos

	var itemLabelLen uint16 = 0
	for i := uint16(0); i < item.l; i++ {
		itemLabelLen, itemPos = nextPos(c.buffer, itemPos, itemLabelLen)
		ch1, _ := c.buffer.At(uint(itemPos))
		ch2, _ := c.nameBuf.ReadUint8()
		if c.caseSensitive {
			if ch1 != ch2 {
				return false
			}
		} else {
			if maptolower[int(ch1)] != maptolower[int(ch2)] {
				return false
			}
		}
		itemPos++
	}

	return true
}

func nextPos(buffer *util.OutputBuffer, pos uint16, llen uint16) (uint16, uint16) {
	if llen == 0 {
		i := 0
		b, _ := buffer.At(uint(pos))
		for ; b&COMPRESS_POINTER_MARK8 == COMPRESS_POINTER_MARK8; i += 2 {
			nb, _ := buffer.At(uint(pos + 1))
			pos = uint16((b & ^uint8(COMPRESS_POINTER_MARK8)))*256 + uint16(nb)
			b, _ = buffer.At(uint(pos))
		}
		return uint16(b), pos
	} else {
		return llen - 1, pos
	}
}

const (
	BUCKETS        uint   = 64
	RESERVED_ITEMS uint   = 16
	NO_OFFSET      uint16 = 65535
)

type MsgRender struct {
	buffer        *util.OutputBuffer
	truncated     bool
	LenLimit      uint32
	caseSensitive bool
	table         [BUCKETS][]offsetItem
	seqHashs      [MAX_LABELS]uint32
}

func NewMsgRender() *MsgRender {
	render := MsgRender{
		buffer:        util.NewOutputBuffer(512),
		truncated:     false,
		LenLimit:      512,
		caseSensitive: false,
	}
	for i := uint(0); i < BUCKETS; i++ {
		render.table[i] = make([]offsetItem, 0, RESERVED_ITEMS)
	}
	return &render
}

func (r *MsgRender) IsTrancated() bool {
	return r.truncated
}

func (r *MsgRender) SetTrancated() {
	r.truncated = true
}

func (r *MsgRender) findOffset(buffer *util.OutputBuffer, nameBuf *util.InputBuffer, hash uint32, caseSensitive bool) uint16 {
	bucketId := hash % uint32(BUCKETS)
	comparator := nameComparator{buffer, nameBuf, hash, caseSensitive}
	found := false

	items := r.table[bucketId]

	i := int(len(items))
	for i -= 1; i >= 0; i-- {
		found = comparator.compare(&items[i])
		if found {
			break
		}
	}

	if found {
		return uint16(items[i].pos)
	} else {
		return NO_OFFSET
	}
}

func (r *MsgRender) addOffset(hash, offset, length uint32) {
	index := hash % uint32(BUCKETS)
	r.table[index] = append(r.table[index], offsetItem{hash, uint16(offset), uint16(length)})
}

func (r *MsgRender) Clear() {
	r.buffer.Clear()
	r.LenLimit = 512
	r.truncated = false
	r.caseSensitive = false
	for i := uint(0); i < BUCKETS; i++ {
		r.table[i] = r.table[i][0:RESERVED_ITEMS]
	}
}

func (r *MsgRender) WriteName(name *Name, compress bool) {
	nlables := name.LabelCount()
	var nlabelsUncomp uint
	ptrOffset := NO_OFFSET

	parent := name
	for nlabelsUncomp = 0; nlabelsUncomp < nlables; nlabelsUncomp++ {
		if nlabelsUncomp > 0 {
			parent, _ = parent.StripLeft(1)
		}

		if parent.Length() == 1 {
			nlabelsUncomp += 1
			break
		}

		r.seqHashs[nlabelsUncomp] = parent.Hash(r.caseSensitive)
		if compress {
			ptrOffset = r.findOffset(r.buffer, util.NewInputBuffer(parent.raw), r.seqHashs[nlabelsUncomp], r.caseSensitive)
			if ptrOffset != NO_OFFSET {
				break
			}
		}
	}

	offset := r.buffer.Len()
	if compress == false || nlabelsUncomp == nlables {
		r.buffer.WriteData(name.raw)
	} else if nlabelsUncomp > 0 {
		compLabelOffset := name.offsets[nlabelsUncomp]
		r.buffer.WriteData(name.raw[0:compLabelOffset])
	}

	if compress && (ptrOffset != NO_OFFSET) {
		ptrOffset |= COMPRESS_POINTER_MARK16
		r.buffer.WriteUint16(ptrOffset)
	}

	nameLen := name.length
	for i := uint(0); i < nlabelsUncomp; i++ {
		labelLen, _ := r.buffer.At(offset)
		if labelLen == 0 {
			break
		}

		if offset > MAX_COMPRESS_POINTER {
			break
		}

		r.addOffset(r.seqHashs[i], uint32(offset), uint32(nameLen))
		offset += uint(labelLen + 1)
		nameLen -= uint(labelLen + 1)
	}
}

func (r *MsgRender) Data() []uint8 {
	return r.buffer.Data()
}

func (r *MsgRender) Len() uint {
	return r.buffer.Len()
}

func (r *MsgRender) Skip(length uint) {
	r.buffer.Skip(length)
}

func (r *MsgRender) Trim(length uint) error {
	return r.buffer.Trim(length)
}

func (r *MsgRender) WriteUint8(data uint8) {
	r.buffer.WriteUint8(data)
}

func (r *MsgRender) WriteUint16(data uint16) {
	r.buffer.WriteUint16(data)
}

func (r *MsgRender) WriteUint16At(data uint16, pos uint) error {
	return r.buffer.WriteUint16At(data, pos)
}

func (r *MsgRender) WriteUint32(data uint32) {
	r.buffer.WriteUint32(data)
}

func (r *MsgRender) WriteData(data []uint8) {
	r.buffer.WriteData(data)
}
