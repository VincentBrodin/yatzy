package ws

import (
	"bytes"
	"encoding/binary"
	"log"
)

type Packet struct {
	CallId  uint32
	Client  *Client
	Message []byte
}

func (p *Packet) Build() []byte {
	callBuf := new(bytes.Buffer)
	if err := binary.Write(callBuf, binary.BigEndian, uint32(p.CallId)); err != nil {
		log.Println(err)
		return make([]byte, 4)
	}
	return append(callBuf.Bytes(), p.Message...)
}
