package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {

}

const Usage = `
	./blockchain.exe createBlockchain ACCOUNT				_添加数据到区块链
	./blockchain.exe printChain						_打印区块链
	./blockchain.exe getBalance	ACCOUNT				_查询余额
	./blockchain.exe send FROM TO AMOUNT MINER DATA	_转账
`

func (cli *CLI) Run() {

	cmds := os.Args

	if len(cmds) < 2 {
		fmt.Println(Usage)
		os.Exit(1)
	}
	switch cmds[1] {
	case "createBlockchain":
		if len(cmds) != 3{
			fmt.Printf("无效的命令，请核实后输入！\n")
			fmt.Println(Usage)
			os.Exit(1)
		}
		fmt.Printf("创建区块链命令被调用\n")
		addr := cmds[2] //TODO
		cli.CreateBlockchain(addr)
	case "printChain":
		fmt.Printf("打印区块链 %s 命令被调用\n", cmds[1])
		cli.printChain()
	case "getBalance":
		fmt.Printf("获取账户 %s 的余额\n", cmds[2])
		cli.GetBalance(cmds[2])
	case "send":
		//./blockchain.exe send FROM TO AMOUNT MINER DATA
		if len(cmds) != 7{
			fmt.Printf("无效的命令，请核实后输入！\n")
			fmt.Println(Usage)
			os.Exit(1)
		}
		from := cmds[2]
		to := cmds[3]
		amount,_ := strconv.ParseFloat(cmds[4],64)
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from,to,amount,miner,data)
	default:
		fmt.Println("无效的命令")
		fmt.Println(Usage)
		os.Exit(1)
	}

}
