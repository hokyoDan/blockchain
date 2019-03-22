package v1

import (
	"fmt"
)

func main() {

	blockChain:= GenerateBlockChains()
	blockChain.AddBlocks("这是第二个区块")
	blockChain.AddBlocks("这是第三个区块")

	for i,block := range blockChain.blockChain{
		fmt.Printf("#%d+++++++++++++++++++++++++\n",i)
		fmt.Printf("前一区块hash：%x\n",block.PreviousHash)
		fmt.Printf("当前区块hash：%x\n",block.PresentHash)
		fmt.Printf("区块数据：%s\n",block.data)
	}

}
