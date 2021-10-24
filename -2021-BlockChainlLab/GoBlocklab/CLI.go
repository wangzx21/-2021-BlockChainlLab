package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) Run() {
	cli.validateArgs()
	//调用NewFlagSet会返回一个带有指定名称和错误处理的空命令集
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	//os.Args[]存储命令行参数的字符串切片，第一个参数是执行文件的名称，因此命令行参数从下标为1开始
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	//解析命令行参数
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}


func(cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data  区块信息")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
}


//判断命令行参数，如果没有输入参数则显示提示信息
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//加入新区块
func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("成功加入区块...")
}

//打印当前区块链的信息
func (cli *CLI) printChain() {
	//创建一个当前区块链的迭代器
	bci := cli.bc.Iterator()

	for {
		//从顶端区块，即新加入的区块开始向前面的区块进行迭代
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}