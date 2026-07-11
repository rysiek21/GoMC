package nbt

import (
	"bytes"
	"encoding/binary"
)

type Tag interface {
	ID() byte
	WritePayload(buf *bytes.Buffer)
}

type Byte byte

func (b Byte) ID() byte { return 0x01 }
func (b Byte) WritePayload(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, b)
}

type Int int32

func (i Int) ID() byte { return 0x03 }
func (i Int) WritePayload(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, i)
}

type String string

func (s String) ID() byte { return 0x08 }
func (s String) WritePayload(buf *bytes.Buffer) {
	stringLen := uint16(len(s))
	binary.Write(buf, binary.BigEndian, stringLen)
	buf.WriteString(string(s))
}

type Float float32

func (f Float) ID() byte { return 0x05 }
func (f Float) WritePayload(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, f)
}

type Double float64

func (d Double) ID() byte { return 0x06 }
func (d Double) WritePayload(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, d)
}

type Compound map[string]Tag

func (c Compound) ID() byte { return 0x0A }
func (c Compound) WritePayload(buf *bytes.Buffer) {
	for name, tag := range c {
		buf.WriteByte(tag.ID())
		nameLen := uint16(len(name))
		binary.Write(buf, binary.BigEndian, nameLen)
		buf.WriteString(name)
		tag.WritePayload(buf)
	}
	buf.WriteByte(0x00)
}

type List []Tag

func (l List) ID() byte { return 0x09 }
func (l List) WritePayload(buf *bytes.Buffer) {
	if len(l) <= 0 {
		buf.WriteByte(0x00)
		Int(0).WritePayload(buf)
		return
	}
	buf.WriteByte(l[0].ID())
	Int(len(l)).WritePayload(buf)
	for _, tag := range l {
		tag.WritePayload(buf)
	}
}

func WriteNamedTag(buf *bytes.Buffer, name string, tag Tag) {
	buf.WriteByte(tag.ID())
	nameLen := uint16(len(name))
	binary.Write(buf, binary.BigEndian, nameLen)
	buf.WriteString(name)
	tag.WritePayload(buf)
}
