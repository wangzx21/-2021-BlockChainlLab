package main

import (
	"github.com/boltdb/bolt"
	"log"
)

// BlockchainIterator 利用一个迭代器访问当前区块链
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next 根据迭代器中currentHash得到序列化后的区块，进行反序列化得到区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		//获取存放当前区块链的Bucket
		b := tx.Bucket([]byte(blocksBucket))
		//根据currentHash获得对应的序列化后的Block
		encodedBlock := b.Get(i.currentHash)
		//反序列化获得Block
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	//将迭代器的currentHash设置为上一块的哈希值
	i.currentHash = block.PrevBlockHash

	return block
}