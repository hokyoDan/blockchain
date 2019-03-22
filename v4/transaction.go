package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type TXInput struct {
	TXId []byte
	Index int64
	Address string
}

type TXOutput struct {
	Value float64
	Address string
}

type Transaction struct {
	TXId []byte
	TXInputs []*TXInput
	TXOutputs []*TXOutput
}


func (tx *Transaction)SetTxHash()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil{
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXId = hash[:]
}

func NewCoinbaseTx(miner,data string)*Transaction  {

	inputs := []*TXInput{&TXInput{nil,-1,data}}
	outputs := []*TXOutput{&TXOutput{12.5,miner}}
	tx := Transaction{nil,inputs,outputs}
	tx.SetTxHash()
	return &tx
}

// 判断该笔交易是否为挖矿交易
func (tx *Transaction)IsCoinbase() bool {

	//coinbase 满足的条件
	if len(tx.TXInputs)==1 && tx.TXInputs[0].TXId==nil && tx.TXInputs[0].Index==-1{
		return true
	}
	return false
}


func NewGeneralTx(from, to string,amount float64,bc BlockChain)*Transaction{

	utxos := make(map[string][]int64)
	var resValue float64

	// 1.遍历账本， 找到属于付款人合适的金额，把这个outputs找到
	utxos,resValue = bc.FindNeedUXTOs(from,amount)

	// 2.如果钱不足以转账，则创建交易失败
	if resValue<amount{
		fmt.Println("余额不足，转装失败！")
		return nil
	}

	// 3.将发款人的outputs转成inputs
	var inputs []*TXInput

	for txid,indexes := range utxos{
		for _,index := range indexes{
			input := TXInput{[]byte(txid),int64(index),from}
			inputs = append(inputs, &input)
		}
	}

	// 4.创建输出，创建一个属于收款人的output
	var outputs []*TXOutput
	output := TXOutput{amount,to}
	outputs = append(outputs, &output)

	// 5.如果有找零，创建属于付款人的output
	if amount < resValue{
		output1 := TXOutput{resValue-amount,from}
		outputs = append(outputs,&output1)
	}

	// 6.设置交易id
	tx := Transaction{nil,inputs,outputs}
	tx.SetTxHash()

	// 7.返回交易id
	return &tx
}