package network

import (
	"bytes"
	"encoding/gob"
	
)


type NetworkEnvelope struct {
	Peerid		 []byte
	Command      []byte // 12 bytes
	Payload      []byte
}

func (ne *NetworkEnvelope) Parse(s []byte) *NetworkEnvelope {
	var m NetworkEnvelope
    b := bytes.Buffer{}
    b.Write(s)
    d := gob.NewDecoder(&b)
    err := d.Decode(&m)
	Handle(err)
    return &m
}

func (ne *NetworkEnvelope) Serialize() []byte {
	b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(ne)
	Handle(err)
    return b.Bytes()
}

func command(com []byte) []byte{
	cl := len(com)
	for i := 0;i <= cl-12;i++{
		com = append(com, 0x00)
	}
	return com
}