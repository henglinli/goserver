package server

import (
	"encoding/binary"
)

// message header
type DynamicHeader struct {
	magiclen int    // magic lenght
	size     uint32 // message length
	buffer   []byte
}

// get
func (this *DynamicHeader) Get() []byte {
	return this.buffer
}

//
func NewDynamicHeader(m string, s uint32) *DynamicHeader {
	header := &DynamicHeader{
		magiclen: len(m),
		size:     s,
		buffer:   make([]byte, len(m)+4),
	}
	// set magic
	copy(header.buffer[0:header.magiclen], m)
	// set size
	binary.BigEndian.PutUint32(
		header.buffer[header.magiclen:header.magiclen+4],
		s)
	//
	return header
}

type DefaultHeader [len(kDefaultMagic) + 4]byte

// check magic
func (this *DefaultHeader) CheckMagic() bool {
	magic := string(this[0:kDefaultMagicLen])
	return magic == kDefaultMagic
}

// get magic
func (this *DefaultHeader) GetMagic() string {
	return string(this[0:kDefaultMagicLen])
}

// get size
func (this *DefaultHeader) GetSize() uint32 {
	return binary.BigEndian.Uint32(this[kDefaultMagicLen : kDefaultMagicLen+4])
}

// set size
func (this *DefaultHeader) SetSize(s uint32) {
	binary.BigEndian.PutUint32(this[kDefaultMagicLen:kDefaultMagicLen+4], s)
}

// Len
func (this *DefaultHeader) Len() int {
	return len(this)
}

// set magic
func NewDefaultHeader() *DefaultHeader {
	var header DefaultHeader
	copy(header[0:kDefaultMagicLen], kDefaultMagic)
	return &header
}

// tiny header

type TinyHeader [4]byte

// get size
func (this *TinyHeader) GetSize() uint32 {
	return binary.BigEndian.Uint32(this[0:4])
}

// set size
func (this *TinyHeader) SetSize(s uint32) {
	binary.BigEndian.PutUint32(this[0:4], s)
}

//
func (this *TinyHeader) CheckMagic() bool {
	return true
}

// get magic
func (this *TinyHeader) GetMagic() string {
	return ""
}

// len
func (this *TinyHeader) Len() int {
	return 4
}

// new
func NewTinyHeader() *TinyHeader {
	header := &TinyHeader{}
	return header
}
