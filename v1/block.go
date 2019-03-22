package v1

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//定义初始数据
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

//定义区块结构
type Block struct {
	Version uint64 //区块版本号

	PreviousHash []byte //前一区块哈希

	MerkleRoot []byte //

	Timestamp uint64 //从1970.1.1至今的秒数

	Difficulty uint64 // 挖矿的难度值

	Nonce uint64 // 随机数

	PresentHash []byte // 当前区块哈希

	data []byte // 数据
}

//生成当前区块hash的方法
func (b *Block)CalcHash()  {
	var data []byte
	/*data = append(data, Uint64ToByte(b.Version)...)
	data = append(data, b.PreviousHash...)
	data = append(data, b.MerkleRoot...)
	data = append(data, Uint64ToByte(b.Timestamp)...)
	data = append(data, Uint64ToByte(b.Difficulty)...)
	data = append(data, Uint64ToByte(b.Nonce)...)
	data = append(data, b.data...)*/

	//bytes.Join 优化代码
	tmp := [][]byte{
		Uint64ToByte(b.Version),
		b.PreviousHash,
		b.MerkleRoot,
		Uint64ToByte(b.Timestamp),
		Uint64ToByte(b.Difficulty),
		Uint64ToByte(b.Nonce),
		b.data,
	}
	data = bytes.Join(tmp,[]byte{})
	hash := sha256.Sum256(data)
	b.PresentHash=hash[:]
}

//生成区块
func GenerateBlocks(data string,PreviousHash []byte)*Block  {
	block := Block{
		Version:1,
		PreviousHash:PreviousHash,
		MerkleRoot:[]byte{},
		Timestamp:uint64(time.Now().Unix()),
		Difficulty:10,
		Nonce:0,
		PresentHash:[]byte{},
		data:[]byte(data),
	}

	block.CalcHash()

	return &block
}

