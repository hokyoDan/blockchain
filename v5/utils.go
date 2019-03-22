package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
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

// 判断文件是否存在
func IsFileExist(name string)bool{
	_,err := os.Stat(name)

	if os.IsNotExist(err){
		return false
	}
	return true
}
