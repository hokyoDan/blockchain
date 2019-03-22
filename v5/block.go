package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)



//定义区块结构
type Block struct {
	Version uint64 //区块版本号

	PreviousHash []byte //前一区块哈希

	MerkleRoot []byte //

	Timestamp uint64 //从1970.1.1至今的秒数

	Difficulty uint64 // 挖矿的难度值

	Nonce uint64 // 随机数

	PresentHash []byte // 当前区块哈希

	//Data []byte // 数据

	Transactions []*Transaction
}

////生成当前区块hash的方法
//func (b *Block)CalcHash()  {
//	var data []byte
//	/*data = append(data, Uint64ToByte(b.Version)...)
//	data = append(data, b.PreviousHash...)
//	data = append(data, b.MerkleRoot...)
//	data = append(data, Uint64ToByte(b.Timestamp)...)
//	data = append(data, Uint64ToByte(b.Difficulty)...)
//	data = append(data, Uint64ToByte(b.Nonce)...)
//	data = append(data, b.data...)*/
//
//	//bytes.Join 优化代码
//	tmp := [][]byte{
//		Uint64ToByte(b.Version),
//		b.PreviousHash,
//		b.MerkleRoot,
//		Uint64ToByte(b.Timestamp),
//		Uint64ToByte(b.Difficulty),
//		Uint64ToByte(b.Nonce),
//		b.data,
//	}
//	data = bytes.Join(tmp,[]byte{})
//	hash := sha256.Sum256(data)
//	b.PresentHash=hash[:]
//}
//

//拼接交易id进行哈希运算，模仿梅克尔根
func (block *Block)NewMerkle()  {
	var hash []byte
	for _,tx := range block.Transactions{
		txid := tx.TXId
		hash =  append(hash,txid...)
	}
	merkleRoot := sha256.Sum256(hash)
	block.MerkleRoot = merkleRoot[:]
}


//生成区块
func GenerateBlocks(txs []*Transaction,PreviousHash []byte)*Block  {
	block := Block{
		Version:1,
		PreviousHash:PreviousHash,
		MerkleRoot:[]byte{},
		Timestamp:uint64(time.Now().Unix()),
		Difficulty:10,
		//Nonce:0,
		PresentHash:[]byte{},
		Transactions:txs,
	}

	//把模拟的merkle根放入其中
	block.NewMerkle()

	pow := GeneratePow(&block)
	hash,nonce := pow.Run()
	block.PresentHash=hash
	block.Nonce=nonce

	return &block
}

//将区块进行序列号编码
func (b *Block)Serialize()[]byte  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}

//将区块进行反序列化解码
func Deserialize(data []byte) *Block  {
	//fmt.Printf("解码传入的数据：%x\n",data)
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic(err)
		fmt.Println("解码错误")
	}
	return &block
}