package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

//定义初始数据
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const block_bucket_01  = "block_bucket_01"
const lastHash  = "lastHash"
const database  = "blockchain.db"

//定义区块链
type BlockChain struct {
	db   *bolt.DB
	tail []byte
}

//添加创世区块至区块链
func GenerateBlockChains() *BlockChain {
	//blockChain := BlockChain{[]*Block{genesisBlock}}

	var tail []byte

	db, err := bolt.Open(database, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))
		if bucket == nil {
			fmt.Println("bucket不存在，准备创建!")
			bucket, err = tx.CreateBucket([]byte(block_bucket_01))
			if err != nil {
				log.Panic(err)
			}
			//生成创世区块
			genesisBlock := GenerateBlocks(genesisInfo, []byte{})
			bucket.Put(genesisBlock.PresentHash, genesisBlock.Serialize())
			bucket.Put([]byte(lastHash), genesisBlock.PresentHash)

			/*//测试是否能读取到解码后的数据
			data := bucket.Get(genesisBlock.PresentHash)
			block := *Deserialize(data)
			fmt.Println(block)*/

			tail = genesisBlock.PresentHash

		} else {
			tail = bucket.Get([]byte(lastHash))
		}

		return nil
	})

	return &BlockChain{db, tail}
}

//添加区块至区块链
func (bc *BlockChain) AddBlocks(data string) {
	bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))
		if bucket == nil {
			fmt.Println("bucket不存在，请检查！")
			os.Exit(1)
		}

		//生成区块
		block := GenerateBlocks(data, bc.tail)
		bucket.Put(block.PresentHash, block.Serialize())
		bucket.Put([]byte(lastHash), block.PresentHash)

		bc.tail = block.PresentHash

		return nil
	})
}

//定义一个迭代器
type BlockchainIterator struct {
	db *bolt.DB
	current []byte
}

//生成一个迭代器
func (bc *BlockChain)NewIterator() *BlockchainIterator {
	return &BlockchainIterator{bc.db,bc.tail}
}

//迭代器Next()函数的实现
func (it *BlockchainIterator)Next() *Block  {
	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))
		if bucket == nil{
			fmt.Println("没有此bucket")
			os.Exit(1)
		}

		blockInfo := bucket.Get(it.current)
		block = *Deserialize(blockInfo)
		it.current = block.PreviousHash
		return nil
	})

	return &block
}