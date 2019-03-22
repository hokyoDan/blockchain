package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	./blockchain.exe addBlock "xxxxxx"   添加数据到区块链
	./blockchain.exe printChain          打印区块链
`

func (cli *CLI)Run()  {

	cmds := os.Args

	if len(cmds)<2{
		fmt.Println(Usage)
		os.Exit(1)
	}
	switch cmds[1]{
	case "addBlock":
		fmt.Printf("添加区块链 %s 命令被调用，添加的数据为：%s\n",cmds[1],cmds[2])
		data := cmds[2]
		cli.AddBlock(data)
	case "printChain":
		fmt.Printf("打印区块链 %s 命令被调用\n",cmds[1])
		cli.printChain()
	default:
		fmt.Println("无效的命令")
		fmt.Println(Usage)
	}

}