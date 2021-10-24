package main
//区块
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

// Block 区块的定义
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int // 用于验证
}

// SetHash 计算块的哈希
func (b *Block) SetHash() {
	//将时间戳转换为String类型并转换为byte[]数组
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	//将PrevBlockHash,Data,timestamp连到一起组成头部
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	//计算header的hash值
	hash := sha256.Sum256(headers)
	//赋值给b.Hash
	b.Hash = hash[:]
}

// NewBlock 创建一个新区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	//创建一个新区块
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{},0}
	//工作量证明
	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}


//数据持久化存储

// Serialize 将Block序列化一个字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// DeserializeBlock 反序列化
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}