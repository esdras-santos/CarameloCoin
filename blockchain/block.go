package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"gochain/utils"
	"log"
	"math/big"
	"time"
)

type Block struct{
	Version []byte
	PrevBlock []byte
	MerkleRoot []byte
	TimeStamp []byte
	Bits []byte
	Nonce []byte
}



func (b *Block) HashTransactions() []byte{
	var txHashes [][]byte
	

	for _,tx := range b.Transactions{
		txHashes = append(txHashes,tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func CreateBlock(txs []*Transaction,prevHash []byte, height int) *Block{
	block := &Block{time.Now().Unix(),[]byte{},txs,prevHash,0,height}

	
	pow := NewProof(block)
	nonce,hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func Genesis(coinbase *Transaction) *Block{
	return CreateBlock([]*Transaction{coinbase},[]byte{},0)
}

func (b *Block) Parse(s []byte) Block {
	version := utils.ToLittleEndian(s[:4],4)
	prevBlock := utils.ToLittleEndian(s[5:36],32)
	merkleRoot := utils.ToLittleEndian(s[37:68],32)
	timeStamp := utils.ToLittleEndian(s[69:72],4)
	bits := s[73:76]
	nonce := s[77:80]
	return Block{version,prevBlock,merkleRoot,timeStamp,bits,nonce}
}

func (b *Block) Serialize() []byte {
	result := utils.ToLittleEndian(b.Version,4)
	result = append(result, utils.ToLittleEndian(b.PrevBlock,32)...)
	result = append(result, utils.ToLittleEndian(b.MerkleRoot,32)...)
	result = append(result, utils.ToLittleEndian(b.TimeStamp,4)...)
	result = append(result, b.Bits...)
	result = append(result, b.Nonce...)
	return result
}

//return the hash of the block in little endian
func (b *Block) Hash() []byte{
	s := b.Serialize()
	sha := sha256.Sum256(s)
	return utils.ToLittleEndian(sha[:],32)
}

func NewProof(b *Block) *ProofOfWork{

	pow := &ProofOfWork{b,b.Target()}

	return pow
}

func (b *Block) Difficulty() *big.Int{
	lowest := big.NewInt(0)
	lowest.Mul(big.NewInt(0xffff),(big.NewInt(0)).Exp(big.NewInt(256),(big.NewInt(0)).Sub(big.NewInt(0x1d),big.NewInt(3)),nil))
	return lowest.Div(lowest,b.Target())
}

func (b *Block) Target() *big.Int{
	exponent := big.NewInt(0)
	coefficient := big.NewInt(0)
	target := big.NewInt(0)
	exponent.SetBytes([]byte{b.Bits[4]})
	coefficient.SetBytes(utils.ToLittleEndian(b.Bits[:3],3))
	target = coefficient.Mul(coefficient,exponent.Exp(big.NewInt(256),exponent.Sub(exponent,big.NewInt(3)),nil))
	return target
}

func (b *Block) Cip9()bool{
	return binary.BigEndian.Uint64(b.Version) >> 29 == 0b001
}

func (b *Block) Cip91() bool{
	return binary.BigEndian.Uint64(b.Version) >> 4 & 1 == 1
}

func (b *Block) Cip141()bool{
	return binary.BigEndian.Uint64(b.Version) >> 1 & 1 == 1
}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}