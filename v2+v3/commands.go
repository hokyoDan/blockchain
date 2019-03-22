package main

import (
	"bytes"
	"fmt"
	"time"
)

func (cli *CLI)AddBlock(data string)  {
	cli.bc.AddBlocks(data)
	fmt.Println("新区块添加成功！")
}

func (cli *CLI)printChain()  {

	it := cli.bc.NewIterator()
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
		fmt.Printf("Data：%s\n",block.Data)


		pow := GeneratePow(block)
		fmt.Printf("IsValid：%t\n",pow.IsValid())
		if bytes.Equal(block.PreviousHash,[]byte{}){
			fmt.Println("区块打印完毕")
			break
		}
	}
}