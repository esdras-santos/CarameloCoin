package blockchain

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"gochain/utils"

	"math"
	"math/big"
)



type ProofOfWork struct{
	Block *Block
}


//each block will be mined in approximately in each 1 minute and in the end of the 24 hours the difficulty will be adjusted
//number of seconds in one day
const ONE_DAY = 86400

//number of block that has to be mined in one day
const BLOCKSPERDAY = 1440

//adjust the bits in the of a period of approximately 1 day.
//prevBits is the previous block bits before the adjust.
//timeDifferential is the timeStamp difference between the first block of the period and the last block of the period.
//return the new bits adjusted
func (pow *ProofOfWork) NewBits(prevBits []byte, timeDifferential int) []byte{
	if timeDifferential > (ONE_DAY * 4){
		timeDifferential = (ONE_DAY * 4)
	}
	if timeDifferential < (ONE_DAY / 4){
		timeDifferential = (ONE_DAY / 4)
	}

	newTarget := big.NewInt(0).Mul(BitsToTarget(prevBits),big.NewInt(int64(timeDifferential /ONE_DAY)))
	return TargetToBits(newTarget)
}

func GetBits(height int64) []byte{
	var chain BlockChain
	var pow ProofOfWork
	
	
	if height == 0{
		return []byte{0x10,0xff,0xff,0xff}
	}else if height % BLOCKSPERDAY == 0{
		lastBlock,err := chain.GetBlock(chain.GetLastHash())
		Handle(err)
		return pow.NewBits(lastBlock.BH.Bits,int(GetTimeDifference()))
		
	}else{
		lastBlock,err := chain.GetBlock(chain.GetLastHash())
		Handle(err)
		return lastBlock.BH.Bits
		
	}
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
	newBits := append(utils.ToLittleEndian(coefficient), byte(exponent))
	return newBits
}

func (pow *ProofOfWork) InitData(nonce int64) []byte{
	pow.Block.BH.Nonce = utils.ToHex(nonce)
	data := pow.Block.Serialize()
	return data
}

func (pow *ProofOfWork) Run()(int){
	var intHash big.Int
	var hash [32]byte
	target := pow.Block.BH.Target()

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
	return nonce
}

func BitsToTarget(bits []byte) *big.Int{
	exponent := *big.NewInt(0)
	coefficient := *big.NewInt(0)
	exponent.SetBytes([]byte{bits[3]})
	coefficient.SetBytes(utils.ToLittleEndian(bits[:3]))
	
	return big.NewInt(0).Mul(&coefficient,big.NewInt(0).Exp(big.NewInt(256),big.NewInt(0).Sub(&exponent,big.NewInt(3)),nil))
}

func  GetTimeDifference() (int64){
	var chain BlockChain
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	for i := 0;i < BLOCKSPERDAY;i++{
		iter.Next()
	}
	firstblock,err := chain.GetBlock(iter.CurrentHash)
	Handle(err)
	lastblock,err := chain.GetBlock(iter.CurrentHash)
	Handle(err)
	tf := binary.BigEndian.Uint64(lastblock.BH.TimeStamp) - binary.BigEndian.Uint64(firstblock.BH.TimeStamp)
	return int64(tf)
}

func (pow *ProofOfWork) Validate()bool{
	var intHash big.Int

	data := pow.Block.Serialize()

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Block.BH.Target()) == -1
}

