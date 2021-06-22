package utils

import (
	"bytes"
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

func EncodeVarint(i int64, buf *[]byte) {
	var bignum, ok = new(big.Int).SetString("0x10000000000000000", 0)
	ibytes := ToHex(i)
	lebytes := ToLittleEndian(ibytes, len(ibytes))
	if !ok {
		log.Panic("fails to create the big number")
	}
	if  i < 0xfd {
		*buf = ibytes
	} else if i < 0x10000 {
		*buf = lebytes
		*buf = append([]byte{0xfd}, *buf...)
	} else if  i < 0x100000000 {
		*buf = lebytes
		*buf = append([]byte{0xfe}, *buf...)
	} else if  big.NewInt(i).Cmp(bignum) < 0 {
		*buf = lebytes
		*buf = append([]byte{0xff}, *buf...)
	} else {
		log.Panic("integer too large")
	}
}

func ToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}	

	return buff.Bytes()
}