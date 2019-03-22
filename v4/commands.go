package main

import (
	"bytes"
	"fmt"
	"time"
)

func (cli *CLI)CreateBlockchain(addr string)  {

	bc := CreateBlockChains(addr)
	if bc == nil{
		return
	}
	defer bc.db.Close()
	fmt.Println("区块链创建成功！")
}

func (cli *CLI)GetBalance(addr string)  {
	bc := GenerateBlockChains()

	if bc == nil{
		return
	}

	defer bc.db.Close()
	bc.GetBalance(addr)
}

func (cli *CLI)printChain()  {
	bc := GenerateBlockChains()

	if bc == nil{
		return
	}

	defer bc.db.Close()
	it := bc.NewIterator()
	for   {
		block := it.Next()
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		fmt.Printf("Version：%d\n",block.Version)
		fmt.Printf("PreviousHash：%x\n",block.PreviousHash)
		fmt.Printf("MerkleRoot：%x\n",block.MerkleRoot)
		timeFormat:= time.Unix(int64(block.Timestamp),0).Format("2006-01-02 15:04:05")
		fmt.Printf("Timestamp：%s\n",timeFormat)
		fmt.Printf("Difficulty：%d\n",block.Difficulty)
		fmt.Printf("Nonce：%d\n",block.Nonce)
		fmt.Printf("PresentHash：%x\n",block.PresentHash)
		fmt.Printf("Data：%s\n",block.Transactions[0].TXInputs[0].Address)


		pow := GeneratePow(block)
		fmt.Printf("IsValid：%t\n",pow.IsValid())
		if bytes.Equal(block.PreviousHash,[]byte{}){
			fmt.Println("区块打印完毕")
			break
		}
	}
}

func (cli *CLI)Send(from,to string,amount float64,miner,data string)  {

	bc := GenerateBlockChains()

	if bc == nil{
		return
	}

	defer bc.db.Close()
	// 1.创造挖矿区块
	coinbase := NewCoinbaseTx(miner,data)
	txs := []*Transaction{coinbase}

	// 2.创造普通区块
	newGeneralTx := NewGeneralTx(from,to,amount,*bc)

	if newGeneralTx != nil {
		txs = append(txs, newGeneralTx)
	}else{
		fmt.Printf("发现无效交易，过滤\n")
	}

	// 3.添加至区块链
	bc.AddBlocks(txs)
	fmt.Println("挖矿成功！")

}