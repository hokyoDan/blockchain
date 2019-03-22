package v1

//定义区块链
type BlockChain struct {
	blockChain []*Block
}

//添加创世区块至区块链
func GenerateBlockChains() *BlockChain  {
	genesisBlock := GenerateBlocks(genesisInfo,[]byte{0x0000000000000000})
	blockChain := BlockChain{[]*Block{genesisBlock}}
	return &blockChain
}

//添加区块至区块链
func (bc *BlockChain)AddBlocks(data string)  {
	lastBlock := bc.blockChain[len(bc.blockChain)-1]
	lastHash := lastBlock.PresentHash
	block := GenerateBlocks(data,lastHash)
	bc.blockChain = append(bc.blockChain,block)
}
