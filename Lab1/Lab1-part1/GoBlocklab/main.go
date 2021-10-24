package main

func main() {
	//先创建一条区块链
	bc := NewBlockchain()
	//程序退出前关闭数据库
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()

	/*bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")*/



	/*for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate())) //验证工作证明
		fmt.Println()
	}*/
}