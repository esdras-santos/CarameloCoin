package blockchain

import (
	"crypto/sha256"
	"fmt"
	"gochain/utils"
	
	"math"
	"math/big"

)



type ProofOfWork struct{
	Block *Block
}


//each block will be mined in approximately in each 6 seconds and in the end of the 24 hours the difficulty will be adjusted
//number of seconds in one day
const ONE_DAY = 86400

//adjust the bits in the of a period of approximately 1 day.
//prevBits is the previous block bits before the adjust.
//timeDifferential is the timeStamp difference between the first block of the period and the last block of the period.
//return the new bits adjusted
func (pow *ProofOfWork) NewBits(prevBits []byte, timeDifferential *big.Int) []byte{
	if timeDifferential.Cmp(big.NewInt(0).Mul(big.NewInt(ONE_DAY),big.NewInt(4))) == 1{
		timeDifferential = big.NewInt(0).Mul(big.NewInt(ONE_DAY),big.NewInt(4))
	}
	if timeDifferential.Cmp(big.NewInt(0).Div(big.NewInt(ONE_DAY),big.NewInt(4))) == -1{
		timeDifferential = big.NewInt(0).Div(big.NewInt(ONE_DAY),big.NewInt(4))
	}

	newTarget := big.NewInt(0).Mul(BitsToTarget(prevBits),big.NewInt(0).Div(timeDifferential,big.NewInt(ONE_DAY)))
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

	return intHash.Cmp(pow.Block.BH.Target()) == -1
}

