package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strings"

	
	"log"
	"math/big"
	"time"
)



type Block struct{
	//blockheader
	Version []byte
	PrevBlock []byte
	MerkleRoot []byte
	TimeStamp []byte
	Bits []byte
	Nonce []byte

	Height []byte
	
	Transactions []*Transaction
}



func (b *Block) ToString() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("     Version:     %d", b.Version))
	lines = append(lines, fmt.Sprintf("          pb:     %x", b.PrevBlock))
	lines = append(lines, fmt.Sprintf("          mr:     %x", b.MerkleRoot))
	lines = append(lines, fmt.Sprintf("          ts:     %x", b.TimeStamp))
	lines = append(lines, fmt.Sprintf("        bits:     %x", b.Bits))
	lines = append(lines, fmt.Sprintf("       Nonce:     %x", b.Nonce))
	lines = append(lines, fmt.Sprintf("           h:     %x", b.Height))
	lines = append(lines, fmt.Sprintf("    tx nonce:     %d", b.Transactions[0].Nonce))
	lines = append(lines, fmt.Sprintf("   tx pubkey:     %x", b.Transactions[0].Pubkey))
	lines = append(lines, fmt.Sprintf("tx receipent:     %x", b.Transactions[0].Receipent))
	lines = append(lines, fmt.Sprintf("      tx sig:     %x", b.Transactions[0].Sig))
	lines = append(lines, fmt.Sprintf("    tx value:     %d", b.Transactions[0].Value))

	return strings.Join(lines, "\n")
}

//get the current block height from the network
func GetBlockHeight() uint{
	return 0
}

//add block height and spread it through the network
func AddBlockHeight(){

}


func (b *Block) HashTransactions() []byte{
	var txHashes [][]byte
	for _,tx := range b.Transactions{
		txb := tx.Serialize()
		txHashes = append(txHashes,txb)
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func CreateBlock(txs []*Transaction,prevHash []byte, height uint64) *Block{
	
	block := Block{Height: ToBytes(height),Transactions: txs}
	ht := block.HashTransactions()
	
	bits := GetBits(int64(height), prevHash)
	//currentTarget() must return the current target of the network and check if the difficult has changed
	block.Version = []byte{0x00000001}
	block.PrevBlock = prevHash
	block.MerkleRoot = ht
	block.TimeStamp = ToBytes(uint64(time.Now().Unix()))
	block.Bits = bits
	block.Nonce = []byte{0x00000000}

	
	
	pow := NewProof(&block)
	
	nonce := pow.Run()
	block.Nonce = ToBytes(uint64(nonce))
	return &block
}

//create genesis block
func Genesis(coinbase *Transaction) *Block{
	
	return CreateBlock([]*Transaction{coinbase},nil,0)
}


func (b *Block) Parse(s []byte) *Block {
	var block Block
    by := bytes.Buffer{}
    by.Write(s)
    d := gob.NewDecoder(&by)
    err := d.Decode(&block)
	Handle(err)

    return &block
}

func (b *Block) Serialize() []byte{
	by := bytes.Buffer{}
    e := gob.NewEncoder(&by)
    err := e.Encode(b)
	Handle(err)
    return by.Bytes()
}

//return the hash of the block in little endian
func (b *Block) Hash() []byte{
	var data []byte

	data = append(data, b.Version...)
	data = append(data, b.PrevBlock...)
	data = append(data, b.MerkleRoot...)
	data = append(data, b.Bits...)
	data = append(data, b.Nonce...)
	data = append(data, b.Height...)
	
	sha := sha256.Sum256(data)
	return sha[:]
}

func NewProof(b *Block) *ProofOfWork{
	pow := ProofOfWork{}
	//check if is the end of the 1 day period
	// if GetBlockHeight() % 8640 == 0{
	// 	b.BH.Bits = pow.NewBits()
	// }
	pow.Block = b

	return &pow
}

func (b *Block) Difficulty() *big.Int{
	
	lowest := big.NewInt(0).Mul(big.NewInt(int64(0xffff)),big.NewInt(0).Exp(big.NewInt(256),big.NewInt(0x1d-3),nil))
	//lowest = 0xffff * int(math.Pow(256,(0x1d - 3)))
	
	return big.NewInt(0).Div(lowest,b.Target()) 
}

func (b *Block) Target() *big.Int{
	target := BitsToTarget(b.Bits)
	return target
}

// func (b *BlockHeader) Cip9()bool{
// 	return binary.BigEndian.Uint64(b.Version) >> 29 == 0b001
// }

// func (b *BlockHeader) Cip91() bool{
// 	return binary.BigEndian.Uint64(b.Version) >> 4 & 1 == 1
// }

// func (b *Block) Cip141()bool{
// 	return binary.BigEndian.Uint64(b.Version) >> 1 & 1 == 1
// }

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}