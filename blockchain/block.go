package blockchain

import (
	"crypto/sha256"
	"encoding/binary"
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
	b.Height = utils.ToHex(int64(binary.LittleEndian.Uint64(s[:8])))
	b.BH = b.BH.Parse(s[8:88])
	txnlen := int64(binary.LittleEndian.Uint64(s[88:96]))
	var lenIn int
	utils.ReadVarint(s[96:],&lenIn)
	var start int
	if lenIn <= 253{
		start = 97
	}else if lenIn <= 254{
		start = 98
	}else if lenIn <= 255{
		start = 99
	}
	var i int64
	end := start+int(lenIn)
	tx := Transaction{}
	for i = 0;i < txnlen;i++{
		b.Transactions = append(b.Transactions, tx.Parse(s[start:end]))
		utils.ReadVarint(s[end:],&lenIn) 
		if lenIn <= 253{
			start = end+1
		}else if lenIn <= 254{
			start = end+2
		}else if lenIn <= 255{
			start = end+3
		}
		end = start+int(lenIn)
	}
}

func (b *Block) Serialize() []byte{
	result := b.Height
	result = append(result, b.BH.Serialize()...)
	result = append(result,  utils.ToHex(int64(len(b.Transactions)))...)
	for _,tx := range b.Transactions{
		txs := tx.Serialize()
		utils.EncodeVarint(int64(len(txs)),&result)
		result = append(result, txs...)
	}
	return result
}

func (b *BlockHeader) Parse(s []byte) BlockHeader {
	version := utils.ToLittleEndian(s[:4])
	prevBlock := utils.ToLittleEndian(s[5:36])
	merkleRoot := utils.ToLittleEndian(s[37:68])
	timeStamp := utils.ToLittleEndian(s[69:72])
	bits := s[72:76]
	nonce := s[76:80]
	return BlockHeader{version,prevBlock,merkleRoot,timeStamp,bits,nonce}
}

func (b *BlockHeader) Serialize() []byte {
	result := utils.ToLittleEndian(b.Version)
	result = append(result, utils.ToLittleEndian(b.PrevBlock)...)
	result = append(result, utils.ToLittleEndian(b.MerkleRoot)...)
	result = append(result, utils.ToLittleEndian(b.TimeStamp)...)
	result = append(result, b.Bits...)
	result = append(result, b.Nonce...)
	return result
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