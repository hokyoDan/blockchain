package v1

import (
	"bytes"
	"encoding/binary"
	"log"
)

//uint64>>>>[]byte

func Uint64ToByte(num uint64)[]byte  {

	var buffer bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}
