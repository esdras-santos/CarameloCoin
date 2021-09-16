package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"gochain/utils"
	"log"
	"math/big"
	"time"
)

type BlockHeader struct{
	Version []byte
	PrevBlock []byte
	MerkleRoot []byte
	TimeStamp []byte
	Bits []byte
	Nonce []byte
}

type Block struct{
	Height []byte
	BH BlockHeader
	Transactions []*Transaction
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

func CreateBlock(txs []*Transaction,prevHash []byte, height int64) *Block{
	
	block := Block{Height: utils.ToHex(height),Transactions: txs}
	ht := block.HashTransactions()
	bits := GetBits(height)
	//currentTarget() must return the current target of the network and check if the difficult has changed
	blockheader := &BlockHeader{[]byte{0x00000001},prevHash, ht,utils.ToHex(time.Now().Unix()),bits,[]byte{0x00000000}}
	block.BH = *blockheader
	
	pow := NewProof(&block)
	nonce := pow.Run()
	block.BH.Nonce = utils.ToHex(int64(nonce))
	return &block
}

//create genesis block
func Genesis(coinbase *Transaction) *Block{
	
	return CreateBlock([]*Transaction{coinbase},nil,0)
}


func (b *Block) Parse(s []byte) {
	var block Block
    by := bytes.Buffer{}
    by.Write(s)
    d := gob.NewDecoder(&by)
    err := d.Decode(&block)
	Handle(err)
    b = &block
}

func (b *Block) Serialize() []byte{
	by := bytes.Buffer{}
    e := gob.NewEncoder(&by)
    err := e.Encode(b)
	Handle(err)
    return by.Bytes()
}

func (b *BlockHeader) Parse(s []byte) BlockHeader {
	var bh BlockHeader
    by := bytes.Buffer{}
    by.Write(s)
    d := gob.NewDecoder(&by)
    err := d.Decode(&bh)
	Handle(err)
    return bh
}

func (b *BlockHeader) Serialize() []byte {
	by := bytes.Buffer{}
    e := gob.NewEncoder(&by)
    err := e.Encode(b)
	Handle(err)
    return by.Bytes()
}

//return the hash of the block in little endian
func (b *BlockHeader) Hash() []byte{
	s := b.Serialize()
	sha := sha256.Sum256(s)
	return utils.ToLittleEndian(sha[:])
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

func (b *BlockHeader) Difficulty() *big.Int{
	
	lowest := big.NewInt(0).Mul(big.NewInt(int64(0xffff)),big.NewInt(0).Exp(big.NewInt(256),big.NewInt(0x1d-3),nil))
	//lowest = 0xffff * int(math.Pow(256,(0x1d - 3)))
	
	return big.NewInt(0).Div(lowest,b.Target()) 
}

func (b *BlockHeader) Target() *big.Int{
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