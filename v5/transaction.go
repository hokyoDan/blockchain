package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"log"
	"math/big"
	"strings"
)

type TXInput struct {
	TXId      []byte
	Index     int64
	Signature []byte
	PubKey    []byte //公钥本身
}

type TXOutput struct {
	Value      float64
	PubKeyHash []byte //是公钥的哈希，而非本身
}

type Transaction struct {
	TXId      []byte
	TXInputs  []*TXInput
	TXOutputs []*TXOutput
}

//给定转账地址，得到这个地址的公钥哈希，完成对output的锁定
func (output *TXOutput) Lock(address string) {
	decodeInfo := base58.Decode(address)
	PubKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	output.PubKeyHash = PubKeyHash
}

//创建一个output方法

func NewTXOutput(address string, value float64) TXOutput {
	output := TXOutput{Value: value}
	output.Lock(address)
	return output
}

func (tx *Transaction) SetTxHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXId = hash[:]
}

const reward = 12.5

func NewCoinbaseTx(miner, data string) *Transaction {

	inputs := []*TXInput{&TXInput{nil, -1, nil, []byte(data)}}
	output := NewTXOutput(miner, reward)
	outputs := []*TXOutput{&output}
	tx := Transaction{nil, inputs, outputs}
	tx.SetTxHash()
	return &tx
}

// 判断该笔交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {

	//coinbase 满足的条件
	if len(tx.TXInputs) == 1 && tx.TXInputs[0].TXId == nil && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

func NewGeneralTx(from, to string, amount float64, bc BlockChain) *Transaction {

	// 1 打开钱包
	ws := NewWalltes()
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("%s 的私钥不存在，交易创建失败!\n", from)
		return nil
	}
	PrivateKey := wallet.PrivateKey
	PublicKey := wallet.PublicKey

	pubKeyHash := hashPubKey(PublicKey)

	utxos := make(map[string][]int64)
	var resValue float64

	// 1.遍历账本， 找到属于付款人合适的金额，把这个outputs找到
	utxos, resValue = bc.FindNeedUXTOs(pubKeyHash, amount)

	// 2.如果钱不足以转账，则创建交易失败
	if resValue < amount {
		fmt.Println("余额不足，转装失败！")
		return nil
	}

	// 3.将发款人的outputs转成inputs
	var inputs []*TXInput

	for txid, indexes := range utxos {
		for _, index := range indexes {
			input := TXInput{[]byte(txid), int64(index), nil, PublicKey}
			inputs = append(inputs, &input)
		}
	}

	// 4.创建输出，创建一个属于收款人的output
	var outputs []*TXOutput
	output := NewTXOutput(to, amount)
	outputs = append(outputs, &output)

	// 5.如果有找零，创建属于付款人的output
	if amount < resValue {
		output1 := NewTXOutput(from, resValue-amount)
		outputs = append(outputs, &output1)
	}

	// 6.设置交易id
	tx := Transaction{nil, inputs, outputs}
	tx.SetTxHash()

	// 7.对交易进行签名
	bc.SignTx(&tx, PrivateKey)

	// 8.返回交易id
	return &tx
}

//对交易进行签名
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	//1. 拷贝一份交易txCopy，
	// >做相应裁剪：把每一个input的Sig和pubkey设置为nil
	// > output不做改变
	TXCopy := tx.TrimmedCopyTX()

	//2. 遍历txCopy.inputs，
	// > 把这个input所引用的output的公钥哈希拿过来，赋值给pubkey
	for i, input := range TXCopy.TXInputs {
		//找到引用的交易
		prevTX := prevTXs[string(input.TXId)]
		output := prevTX.TXOutputs[input.Index]

		//for循环迭代出来的数据是一个副本，对这个input进行修改，不会影响到原始数据
		//所以我们这里需要使用下标方式修改

		//签名要对数据的hash进行签名
		//我们的数据都在交易中，我们要求交易的哈希
		//Transaction的SetTXID函数就是对交易的哈希
		//所以我们可以使用交易id作为我们的签名的内容
		//所以不要写成： input.PubKey=output.PubKeyHash ，要用下面这种

		TXCopy.TXInputs[i].PubKey = output.PubKeyHash

		//3. 生成要签名的数据（哈希）
		TXCopy.SetTxHash()
		signData := TXCopy.TXId

		//清理,原理同上
		TXCopy.TXInputs[i].PubKey = nil

		fmt.Printf("signData:%x\n", signData)

		//4. 对数据进行签名r, s
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signData)
		if err != nil {
			fmt.Printf("交易签名失败,err:", err)
		}

		//5. 拼接r,s为字节流
		signature := append(r.Bytes(), s.Bytes()...)

		//6. 赋值给原始的交易的Signature字段
		tx.TXInputs[i].Signature = signature

	}

	fmt.Println("对交易进行签名")
}

//对交易进行验证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {

	//1. 拷贝修剪的副本
	TXCopy := tx.TrimmedCopyTX()

	//2. 遍历原始交易（注意，不是txCopy)
	for i,input := range tx.TXInputs{

		// 3. 根据原始交易找到索引的output
		prevTx := prevTXs[string(input.TXId)]
		outputs := prevTx.TXOutputs

		//4. 找到对应output的 publicKeyHash 赋值给拷贝出来副本的 publicKey 字段
		TXCopy.TXInputs[i].PubKey=outputs[i].PubKeyHash

		//5.生成交易的hash（id）,得到待验证的数据
		TXCopy.SetTxHash()
		verifyData := TXCopy.TXId
		fmt.Printf("verifyData:%x\n", verifyData)

		// 清理数据
		TXCopy.TXInputs[i].PubKey=nil

		//6.得到pubKey 和 Sig
		sig := input.Signature
		pubKey := input.PubKey

		//7. 拼出原始的pubkey和sig

		var r big.Int
		var s big.Int
		rData := sig[:len(sig)/2]
		sData := sig[len(sig)/2:]
		r.SetBytes(rData)
		s.SetBytes(sData)

		var x big.Int
		var y big.Int
		xData := pubKey[:len(pubKey)/2]
		yData := pubKey[len(pubKey)/2:]
		x.SetBytes(xData)
		y.SetBytes(yData)
		curve := elliptic.P256()
		pubKeyRaw := ecdsa.PublicKey{curve,&x,&y}

		//8. 进行验证
		if !ecdsa.Verify(&pubKeyRaw,verifyData,&r,&s){
			return false
		}
		return true

	}

	fmt.Println("对交易进行签名")
	return true
}

func (tx *Transaction) TrimmedCopyTX() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput
	for _, input := range tx.TXInputs {
		input1 := TXInput{input.TXId, input.Index, nil, nil}
		inputs = append(inputs, &input1)
	}
	outputs = tx.TXOutputs

	TX1 := Transaction{tx.TXId, inputs, outputs}
	return TX1
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("	--- Transaction %x:", tx.TXId))

	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("		Input %d:", i))
		lines = append(lines, fmt.Sprintf("		TXID:      %x", input.TXId))
		lines = append(lines, fmt.Sprintf("		Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("		Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("		PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs {
		lines = append(lines, fmt.Sprintf("		Output %d:", i))
		lines = append(lines, fmt.Sprintf("		Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("		Script: %x", output.PubKeyHash))
	}
	lines = append(lines,"\n")
	lines = append(lines,"\n")
	//11111, 2222, 3333, 44444, 5555

	//`11111
	//2222
	//3333
	//44444
	//5555`

	return strings.Join(lines, "\n")
}

