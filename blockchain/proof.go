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
	Target *big.Int
}

//each block will be mined in approximately in each 6 seconds and in the end of the 24 hours the difficulty will be adjusted
//number of seconds in one day
var ONE_DAY int =  86400
var DIFFICULTYDB *badger.DB



func (pow *ProofOfWork) DifficultyAdjustment(){

}

//this will be used to update the fb(first block after adjust) and lb(last block after the adjust)
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

	nonce := 0

	for nonce < math.MaxInt64{
		data := pow.InitData(int64(nonce))
		hash = sha256.Sum256(data)	
		fmt.Printf("\r%x",hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1{
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

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}	

	return buff.Bytes()
}