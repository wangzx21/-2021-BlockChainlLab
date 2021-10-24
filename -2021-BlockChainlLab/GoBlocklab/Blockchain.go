package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//区块链
const dbFile = "blockchain.db"
const blocksBucket = "blocks"


// Blockchain 区块链的定义
// 存储最后一个块的哈希
type Blockchain struct {
	tip []byte
	db *bolt.DB
}

// Iterator 访问当前区块链时首先创建一个当前区块链的迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// AddBlock 向区块链添加一个新区块
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	//获取最后一个块的哈希，用来生成新的哈希
	err := bc.db.View(func(tx *bolt.Tx) error {
		//获取当前区块链
		b := tx.Bucket([]byte(blocksBucket))
		//获取最后一块区块哈希值，用于生成新区块
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}


	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//将新区块存入数据库中
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		//将顶部的区块哈希更新为新插入的区块
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

// NewGenesisBlock 创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain 以创世区块为头部创建新的区块链
func NewBlockchain() *Blockchain {
	var tip []byte
	//打开一个BoltDB文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		//先获取存储区块的bucket
		b := tx.Bucket([]byte(blocksBucket))

		//blocksBucket不存在则说明数据库中不存在区块链，创建一个，否则直接读取最后一个块的哈希
		if b == nil {
			fmt.Println("No existing blockchain found.Creating a new one...")
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			//写入键值对，区块哈希对应序列化后的区块
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//"l"键对应区块链顶端区块的哈希
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			//指向最后一个区块，这里也就是创世区块
			tip = genesis.Hash
		} else {
			//如果存在blocksBucket桶，也就是存在区块链
			//通过键"l"映射出顶端区块的Hash值
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}