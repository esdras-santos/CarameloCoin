package utils

import (
	"encoding/binary"
	"log"
	"math/big"
)


func ToLittleEndian(bytes []byte, length int) []byte {
	le := make([]byte, length)
	for i := len(le) - 1; i >= 0; i-- {
		if bytes[i] != 0x00 {
			le = append(le, bytes[i])
		}
		le = append(le, 0x00)
	}
	return le
}

func ReadVarint(s []byte, buf *uint){
	i := s[0]
	if i == 0xfd{
		a := binary.LittleEndian.Uint16(s[1:3])
		*buf = uint(a)
	}else if i == 0xfe{
		a := binary.LittleEndian.Uint32(s[1:5])
		*buf = uint(a)
	}else if i == 0xff{
		a := binary.LittleEndian.Uint64(s[1:9])
		*buf = uint(a)
	}else{
		*buf = uint(i)
	}
}

func EncodeVarint(i big.Int, buf *[]byte) {
	var bignum, ok = new(big.Int).SetString("0x10000000000000000", 0)
	ibytes := i.Bytes()
	lebytes := ToLittleEndian(ibytes, 4)
	if !ok {
		log.Panic("fails to create the big number")
	}
	if cmp := i.Cmp(big.NewInt(0xfd)); cmp < 0 {
		*buf = ibytes
	} else if cmp := i.Cmp(big.NewInt(0x10000)); cmp < 0 {
		*buf = lebytes
		*buf = append([]byte{0xfd}, *buf...)
	} else if cmp := i.Cmp(big.NewInt(0x100000000)); cmp < 0 {
		*buf = lebytes
		*buf = append([]byte{0xfe}, *buf...)
	} else if cmp := i.Cmp(bignum); cmp < 0 {
		*buf = lebytes
		*buf = append([]byte{0xff}, *buf...)
	} else {
		log.Panic("integer too large")
	}
}
