# blockchain

## 命令介绍

### run.sh 
run.sh  文件会移除现有的.db区块链文件以及生成的blockchain.exe文件，然后新生成一个blockchain.exe的可执行程序。

### init.sh
执行 init.sh 会创造创世区块，默认由地址`1N3Ycacu4xAbVaUP18sGFcRGKXjk2tVfTA`生成

### send.sh
执行 send.sh 发送转账，默认由`1N3Ycacu4xAbVaUP18sGFcRGKXjk2tVfTA`向账户`1N2mh8woqmUTD6qRexcTcnNfXRuVxoR6pf`转2个比特币。由`1mm4J91v9hQhsNadE5cq4BQpMEAMqzka2`进行打包挖矿，附带的信息为 `test`

### blockchain.exe
```./blockchain.exe createBlockchain ACCOUNT    _添加创世区块链,ACCOUNT指定创建人
./blockchain.exe printChain    _打印区块链
./blockchain.exe getBalance ACCOUNT    _查询余额，ACCOUNT指定被查询者
./blockchain.exe send FROM TO AMOUNT MINER DATA     _转账，FROM:转账发起者。 TO：转账接受者。 AMOUNT:转账单位，单位为BTC。 MINER:指定矿工地址。  DATA:指定附加的交易数据
./blockchain.exe createWallet    _创建新地址（不会覆盖原有钱包）
./blockchain.exe listAddress    _打印钱包地址
./blockchain.exe showTx    _打印交易
```

## 文件介绍

### block.go
有关区块的相关方法

### blockchain.go
有关区块链的方法

### cli.go
解析命令提示行

### commands.go
cli.go中的命令以此文件为路由，调用相关的方法

### main.go
入口函数

### proofofwork.go
工作量证明的相关方法

### transcation.go
关于交易的相关方法

### utils.go
一些工具函数

### wallet.go
生成钱包对应的公私钥，地址等等

### wallet pub.go
对外开放的文件，调用wallet中的方法

### wallet.dat
wallet.dat 是存储公私钥的文件，默认存储五个地址对应的公私钥：
```address 0: 1JXPakHFHiLZAWDRpetYKMf6gpVK39aVEY
address 1: 1NGWo9mPRingjfhwXSr3E25FxhWgBWtoK7
address 2: 1NKnYSoLn982bzhGKKP34CTtNsC1gFziEk
address 3: 1mm4J91v9hQhsNadE5cq4BQpMEAMqzka2
address 4: 1N2mh8woqmUTD6qRexcTcnNfXRuVxoR6pf
address 5: 1N3Ycacu4xAbVaUP18sGFcRGKXjk2tVfTA
```

### blockchain.db
区块链文件
