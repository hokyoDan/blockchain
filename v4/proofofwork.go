package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	Block *Block
	Target *big.Int
}

//定义右移位常量
const Bits  = 16

func GeneratePow(block *Block) *ProofOfWork {

	/*
	//目标值选定为固定值
	var tmp big.Int
	initTarget := "0000100000000000000000000000000000000000000000000000000000000000"
	tmp.SetString(initTarget,16)
	*/

	//设置目标值为可变

	tmp := big.NewInt(1)
	tmp.Lsh(tmp,256-Bits)

	pow := ProofOfWork{
		Block:block,
		Target:tmp,
	}
	return &pow

}

func (pow *ProofOfWork)Run() ([]byte,uint64) {

	var nonce uint64
	var hash [32]byte
	var bigIntHash big.Int
	for ; ;  {
		hash =  sha256.Sum256(pow.PrepareData(nonce))
		fmt.Printf("当前计算哈希为：%x\r",hash)
		bigIntHash.SetBytes(hash[:])

		if bigIntHash.Cmp(pow.Target) ==-1 {
			fmt.Printf("挖矿成功！当前区块哈希为：%x,\n随机数为：%d\n",hash,nonce)
			break
		}else{
			nonce++
		}
	}
	return hash[:],nonce
}

func (pow *ProofOfWork)PrepareData(nonce uint64) []byte  {
	block := pow.Block
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PreviousHash,
		block.MerkleRoot,
		Uint64ToByte(block.Timestamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(nonce),
	}
	data := bytes.Join(tmp,[]byte{})
	return data
}


func (pow *ProofOfWork)IsValid()bool  {
	data := pow.PrepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	var tmp big.Int
	tmp.SetBytes(hash[:])
	return tmp.Cmp(pow.Target)==-1
}