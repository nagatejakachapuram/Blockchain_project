package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
	balances     map[string]int
}

func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
		map[string]int{},
	}
}

func (b *Block) calculateHash() string {
	blockData, err := json.Marshal(b.data)
	if err != nil {
		panic(err)
	}
	x := fmt.Sprintf("%x%s%x%d", blockData, b.previousHash, b.timestamp, b.pow)
	h := sha256.New()
	h.Write([]byte(x))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

func (b *Blockchain) addBlock(from, to string, amount int) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
	b.updateBalances(from, to, amount)
}

func (b *Blockchain) updateBalances(from, to string, amount int) {
	if _, ok := b.balances[from]; !ok {
		b.balances[from] = 0
	}
	if _, ok := b.balances[to]; !ok {
		b.balances[to] = 0
	}
	b.balances[from] -= amount
	b.balances[to] += amount
}

func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func main() {
	blockchain := CreateBlockchain(2)
	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)
	fmt.Println(blockchain.isValid())
	fmt.Println(blockchain.balances)
}
