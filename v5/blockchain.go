package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcutil/base58"
	"log"
	"os"
)

//定义初始数据
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const block_bucket_01 = "block_bucket_01"
const lastHash = "lastHash"
const database = "blockchain.db"

//定义区块链
type BlockChain struct {
	db   *bolt.DB
	tail []byte
}

//添加创世区块至区块链//TODO
func CreateBlockChains(miner string) *BlockChain {
	//blockChain := BlockChain{[]*Block{genesisBlock}}

	if IsFileExist(database) {
		fmt.Println("文件已存在，不需要重复创建！")
		return nil
	}

	var tail []byte

	db, err := bolt.Open(database, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))

		bucket, err = tx.CreateBucket([]byte(block_bucket_01))
		if err != nil {
			log.Panic(err)
		}
		//生成创世区块
		//创世块中只有一个挖矿交易
		coinbase := NewCoinbaseTx(miner, genesisInfo)
		genesisBlock := GenerateBlocks([]*Transaction{coinbase}, []byte{})
		bucket.Put(genesisBlock.PresentHash, genesisBlock.Serialize())
		bucket.Put([]byte(lastHash), genesisBlock.PresentHash)

		/*//测试是否能读取到解码后的数据
		data := bucket.Get(genesisBlock.PresentHash)
		block := *Deserialize(data)
		fmt.Println(block)*/

		tail = genesisBlock.PresentHash

		return nil
	})

	return &BlockChain{db, tail}
}

//显示区块链 返回区块链实例
func GenerateBlockChains() *BlockChain {
	//blockChain := BlockChain{[]*Block{genesisBlock}}

	if !IsFileExist(database) {
		fmt.Println("文件不已存在，请先创建文件！")
		return nil
	}

	var tail []byte

	db, err := bolt.Open(database, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))
		if bucket == nil {
			fmt.Println("bucket不存在，请创建!")
			os.Exit(1)

		}
		tail = bucket.Get([]byte(lastHash))

		return nil
	})

	return &BlockChain{db, tail}
}

//添加区块至区块链
func (bc *BlockChain) AddBlocks(txs []*Transaction) {

	// 将有效的交易添加至区块
	validTXs :=  []*Transaction{}
	for _,validTX := range txs{
		if bc.VerifyTx(validTX){
			validTXs =append(validTXs,validTX)
			fmt.Printf("找到有效交易，id：%x\n",validTX.TXId)
		}else{
			fmt.Printf("该笔交易无效，id：%x\n",validTX.TXId)
		}
	}


	bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))
		if bucket == nil {
			fmt.Println("bucket不存在，请检查！")
			os.Exit(1)
		}

		//生成区块
		block := GenerateBlocks(txs, bc.tail)
		bucket.Put(block.PresentHash, block.Serialize())
		bucket.Put([]byte(lastHash), block.PresentHash)

		bc.tail = block.PresentHash

		return nil
	})
}

//定义一个迭代器
type BlockchainIterator struct {
	db      *bolt.DB
	current []byte
}

//生成一个迭代器
func (bc *BlockChain) NewIterator() *BlockchainIterator {
	return &BlockchainIterator{bc.db, bc.tail}
}

//迭代器Next()函数的实现，返回值为当前的区块，查找当前的区块
func (it *BlockchainIterator) Next() *Block {
	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(block_bucket_01))
		if bucket == nil {
			fmt.Println("没有此bucket")
			os.Exit(1)
		}

		blockInfo := bucket.Get(it.current) //it.current 是当前区块的哈希值
		block = *Deserialize(blockInfo)
		it.current = block.PreviousHash
		return nil
	})

	return &block
}

// 定义一个UTXOInfo的结构体，
//UTXOInfo
//
//1. TXID
//2. index
//3. output

type UTXOInfo struct {
	TXID   []byte
	Index  int64
	Output TXOutput
}

// 统计utxo
func (bc *BlockChain) FindAllUtxos(pubKeyHash []byte) []UTXOInfo {

	var UTXOInfos []UTXOInfo

	spentInputs := make(map[string][]int64)
	it := bc.NewIterator()

	//遍历区块链账本
	for {
		block := it.Next()

		//遍历每一笔交易
		for _, tx := range block.Transactions {

			//遍历每个交易中的每个输入
			//遍历输入之前判断该笔输入是否为挖矿交易，若是挖矿交易则跳过
			if !tx.IsCoinbase() {
				for _, input := range tx.TXInputs {

					//检查遍历到的输入地址是否和输入地址相同，若相同，则说明该input对应的output的钱被花过了，不在统计范围内
					if bytes.Equal(hashPubKey(input.PubKey), pubKeyHash) {
						fmt.Printf("找到了消耗过的output，id：%s 号交易中下标为 %d 的交易\n", input.TXId, input.Index)
						key := string(input.TXId)
						spentInputs[key] = append(spentInputs[key], input.Index)
					}
				}
			}

			key := string(tx.TXId)
			indexes := spentInputs[key]

		OUTPUT:
			//遍历每个交易中的每个输出
			for i, output := range tx.TXOutputs {

				if len(indexes) != 0 {
					for _, index := range indexes {
						if int64(i) == index {
							fmt.Printf("i==j,说明当前output已经被花掉了，跳过，不计入统计当中\n")
							continue OUTPUT
						}
					}
				}

				//找到属于我的output
				if bytes.Equal(output.PubKeyHash, pubKeyHash) {
					//fmt.Printf("账户 %s 中的第 %d 笔交易\n", address, i)
					utxoinfo := UTXOInfo{tx.TXId, int64(i), *output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}
			}

		}
		if len(block.PreviousHash) == 0 {
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}
	return UTXOInfos
}

func (bc *BlockChain) GetBalance(address string) {
	//var balance float64

	decodeInfo := base58.Decode(address)
	PubKeyHash := decodeInfo[1 : len(decodeInfo)-4]

	utxoinfos := bc.FindAllUtxos(PubKeyHash)

	var balance = 0.0

	for _, utxoinfo := range utxoinfos {
		balance += utxoinfo.Output.Value
	}
	fmt.Printf("账户：‘%s’的余额为：‘%f’\n", address, balance)
}

//// 1.遍历账本， 找到属于付款人合适的金额，把这个outputs找到
//	utxos,resValue = FindNeedUXTOs(from,amount)

func (bc *BlockChain) FindNeedUXTOs(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {

	NeedUtxos := make(map[string][]int64)
	var resValue float64

	//  复用 FindAllUtxos 代码
	UTXOInfos := bc.FindAllUtxos(pubKeyHash)
	for _, utxoinfo := range UTXOInfos {
		key := string(utxoinfo.TXID)
		NeedUtxos[key] = append(NeedUtxos[key], utxoinfo.Index)

		// 2. 判断金额是否足够，足够：返回，不足：继续遍历
		resValue += utxoinfo.Output.Value
		if resValue >= amount {
			break
		}
	}
	return NeedUtxos, resValue

}

func (bc *BlockChain) SignTx(tx *Transaction, privateKey *ecdsa.PrivateKey) {

	// 如果是挖矿交易，直接返回
	if tx.IsCoinbase(){
		return
	}

	prevTXs := make(map[string]Transaction)

	//设置包含交易id的切片数组
	TXIds := make([][]byte, 0)
	TXID1 := make([]byte, 0)
	for _, TXInput := range tx.TXInputs {

		i := 0
		for {
			TXID1 = append(TXID1, TXInput.TXId...)
			TXIds = append(TXIds, TXID1)
			i++
			if i == len(tx.TXInputs) {
				break
			}
		}

		txs := bc.FindTransaction(TXIds)
		if txs == nil {
			fmt.Println("没有找到该input对应的output")
			return
		}

		//如果该input对应的output存在，不执行，进入下一次循环，主要为了优化代码，若没有这段直接覆盖map的效率低
		_, ok := prevTXs[string(TXInput.TXId)]
		if ok {
			continue
		}

		//当确定该如果该input对应的output是第一次遍历到时，才放入map
		for _, tx := range txs {
			prevTXs[string(tx.TXId)] = *tx
		}

	}

	tx.Sign(privateKey, prevTXs)

}

func (bc *BlockChain) VerifyTx(tx *Transaction) bool{

	// 如果是挖矿交易，返回true
	if tx.IsCoinbase(){
		return true
	}

	prevTXs := make(map[string]Transaction)

	//设置包含交易id的切片数组
	TXIds := make([][]byte, 0)
	TXID1 := make([]byte, 0)
	for _, TXInput := range tx.TXInputs {

		i := 0
		for {
			TXID1 = append(TXID1, TXInput.TXId...)
			TXIds = append(TXIds, TXID1)
			i++
			if i == len(tx.TXInputs) {
				break
			}
		}

		txs := bc.FindTransaction(TXIds)
		if txs == nil {
			fmt.Println("没有找到该input对应的output")
			return false
		}

		//如果该input对应的output存在，不执行，进入下一次循环，主要为了优化代码，若没有这段直接覆盖map的效率低
		_, ok := prevTXs[string(TXInput.TXId)]
		if ok {
			continue
		}

		//当确定该如果该input对应的output是第一次遍历到时，才放入map
		for _, tx := range txs {
			prevTXs[string(tx.TXId)] = *tx
		}

	}

	return tx.Verify(prevTXs)

}


func (bc *BlockChain) FindTransaction(TXIds [][]byte) []*Transaction {

	it := bc.NewIterator()
	TXs := make([]*Transaction, 0)
	for {
		block := it.Next()

		for _, tx := range block.Transactions {

			for _, txid := range TXIds {
				if bytes.Equal(txid, tx.TXId) {
					fmt.Printf("交易id：%x\n", tx.TXId)
					TXs = append(TXs, tx)
				}
			}
		}

		if len(block.PreviousHash) == 0 {
			return TXs
		}
	}
	return nil
}
