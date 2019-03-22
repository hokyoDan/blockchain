package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
	"os"
)

type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWalletKeyPair() *WalletKeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		os.Exit(1)
	}
	publicKeyRaw := privateKey.PublicKey
	publicKey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{privateKey, publicKey}

}

func (w *WalletKeyPair) GetAddress() string {

	publicHash := hashPubKey(w.PublicKey)

	version := 0x00

	//21字节的数据
	payload := append([]byte{byte(version)}, publicHash...)

	//调用四字节
	checksum := checkSum(payload)
	//25字节
	payload = append(payload, checksum...)

	address := base58.Encode(payload)
	return address
}

func IsValidAddress(address string)bool  {
	decodeInfo := base58.Decode(address)

	if len(decodeInfo) != 25 {
		return false
	}

	checksum := checkSum(decodeInfo[:len(decodeInfo)-4])
	checksum1 := decodeInfo[len(decodeInfo)-4:]

	return bytes.Equal(checksum,checksum1)
}

func hashPubKey(publicKey []byte) []byte {
	hash := sha256.Sum256(publicKey)
	rip160Hasher := ripemd160.New()
	_, err := rip160Hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	publicHash := rip160Hasher.Sum(nil)
	return publicHash
}

func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	//4字节校验码
	checksum := second[:4]
	return checksum
}
