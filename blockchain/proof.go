package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"gochain/utils"
	"log"
	"math"
	"math/big"

	"github.com/dgraph-io/badger"
)



type ProofOfWork struct{
	Block *Block
}

//the db need to be changed later
//use LevelDB
var DIFFICULTYDB *badger.DB
const (
	DBPATH = "./tmp/difficultyadjust"
	//each block will be mined in approximately in each 6 seconds and in the end of the 24 hours the difficulty will be adjusted
	//number of seconds in one day
	ONE_DAY =  86400
)

//adjust the difficulty of the bits
func (pow *ProofOfWork) NewBits() []byte{
	lastBlock := pow.GetPosition("lb")
	firstBlock := pow.GetPosition("fb")
	fts := big.NewInt(0).SetBytes(firstBlock.TimeStamp)
	lts := big.NewInt(0).SetBytes(lastBlock.TimeStamp)
	timeDifferential := big.NewInt(0).Sub(lts,fts)
	if timeDifferential.Cmp(big.NewInt(0).Mul(big.NewInt(ONE_DAY),big.NewInt(4))) == 1{
		timeDifferential = big.NewInt(0).Mul(big.NewInt(ONE_DAY),big.NewInt(4))
	}
	if timeDifferential.Cmp(big.NewInt(0).Div(big.NewInt(ONE_DAY),big.NewInt(4))) == -1{
		timeDifferential = big.NewInt(0).Div(big.NewInt(ONE_DAY),big.NewInt(4))
	}

	newTarget := big.NewInt(0).Mul(lastBlock.Target(),timeDifferential.Div(timeDifferential,big.NewInt(ONE_DAY)))
	return TargetToBits(newTarget)
}

func TargetToBits(target *big.Int) []byte{
	var exponent int
	var coefficient []byte
	rawBytes := target.Bytes()
	if rawBytes[0] > 0x7f{
		exponent = len(rawBytes) + 1
		coefficient = append([]byte{0x00}, rawBytes[:2]...)
	}else{
		exponent = len(rawBytes)
		coefficient = rawBytes[:3]
	}
	newBits := append(utils.ToLittleEndian(coefficient,len(coefficient)), byte(exponent))
	return newBits
}


//will be used LevelDB
func (pow *ProofOfWork) GetPosition(position string) Block{
	return Block{}
}

//this will be used to update the fb(first block after adjust) and lb(last block after the adjust)
//the db need to be changed later
func (pow *ProofOfWork) UpdatePosition(block *Block,position string){
	err := DIFFICULTYDB.Update(func(txn *badger.Txn) error{
		if _,err := txn.Get(block.Hash()); err == nil{
			return nil
		}

		blockData := block.Serialize()
		err := txn.Set(block.Hash(), blockData)
		Handle(err)

		err = txn.Set([]byte(position), block.Hash())
		Handle(err)

		return nil
	})
	Handle(err)
}


func (pow *ProofOfWork) InitData(nonce int64) []byte{
	pow.Block.Nonce = ToHex(nonce)
	data := pow.Block.Serialize()
	return data
}

func (pow *ProofOfWork) Run()(int,[]byte){
	var intHash big.Int
	var hash [32]byte
	target := pow.Block.Target()

	nonce := 0

	for nonce < math.MaxInt64{
		data := pow.InitData(int64(nonce))
		hash = sha256.Sum256(data)	
		fmt.Printf("\r%x",hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(target) == -1{
			break
		}else{
			nonce++
		}
	}
	fmt.Println()
	return nonce,hash[:]
}

func BitsToTarget(bits []byte) *big.Int{
	exponent := big.NewInt(0)
	coefficient := big.NewInt(0)
	exponent.SetBytes([]byte{bits[4]})
	coefficient.SetBytes(utils.ToLittleEndian(bits[:3],3))
	return coefficient.Mul(coefficient,exponent.Exp(big.NewInt(256),exponent.Sub(exponent,big.NewInt(3)),nil))
}

func (pow *ProofOfWork) Validate()bool{
	var intHash big.Int

	data := pow.Block.Serialize()

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Block.Target()) == -1
}

func ToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}	

	return buff.Bytes()
}