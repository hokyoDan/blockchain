package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

type Wallets struct {
	WalletsMap map[string]*WalletKeyPair
}

//创建wallets 返回wallets实例
func NewWalltes()*Wallets  {

	var ws Wallets
	ws.WalletsMap = make(map[string]*WalletKeyPair)

	//加载本地钱包
	if !ws.LoadFromFile(){
		fmt.Println("加载钱包失败")
	}
	//把所有的钱包从本地加载出来
	return &ws
}

//创建新的钱包
func (ws *Wallets)CreateWallets()string  {

	wallet := NewWalletKeyPair()
	address := wallet.GetAddress()
	ws.WalletsMap[address]=wallet

	// 保存到本地文件
	res := ws.SaveToFile()
	if !res{
		return "生成地址失败"
	}
	return address
}

const walletName  = "wallet.dat"

func (ws *Wallets)SaveToFile() bool {

	var buffer bytes.Buffer

	//将接口类型明确注册一下，否则gob编码失败!
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil{
		fmt.Println("序列化钱包失败，err",err)
		return false
	}

	content := buffer.Bytes()

	err = ioutil.WriteFile(walletName,content,0600)
	if err != nil{
		fmt.Println("钱包创建失败，err",err)
		return false
	}
	return true
}

//将钱包文件从本地中读取出来
func (ws *Wallets)LoadFromFile()bool  {

	//钱包文件不存在的时候要主动创建
	if !IsFileExist(walletName){
		fmt.Println("钱包文件不存在，准备创建")
		return true
	}

	content,err := ioutil.ReadFile(walletName)
	if err != nil{
		fmt.Println("钱包读取失败：",err)
		return false
	}

	//gob解码
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wallets)
	if err != nil{
		fmt.Println("钱包解码失败：",err)
		return false
	}

	//赋值给ws
	ws.WalletsMap=wallets.WalletsMap

	return true
}

func (ws *Wallets)ListAddress()[]string  {

	var addresses []string
	for address,_ := range ws.WalletsMap{
		addresses = append(addresses,address )
	}
	return addresses
}
