package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const(
	checksumLength = 4
	version = byte(0x00)
)

type Wallet struct{
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w Wallet) Address() []byte{
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version},pubHash...)
	checksum := CheckSum(versionedHash)
	fullHash := append(versionedHash,checksum...)
	address := Base58Encode(fullHash)
	return address
}

func AddressToPKH(address string) []byte{
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-checksumLength]
	return pubKeyHash
}

func ValidateAddress(address string) bool{
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-checksumLength]
	targetChecksum := CheckSum(append([]byte{version},pubKeyHash...))
	return bytes.Compare(actualChecksum,targetChecksum) == 0
}

func NewKeyPair(compressed bool) (ecdsa.PrivateKey,[]byte){
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve,rand.Reader)
	Handle(err)

	var pub []byte
	
	pub = append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pub
}



func MakeWallet() *Wallet{
	private,public := NewKeyPair(true)
	wallet := Wallet{private,public}
	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte{
	pubHash := sha256.Sum256(pubKey)
	
	hasher := ripemd160.New()
	_,err := hasher.Write(pubHash[:])
	Handle(err)
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

/*before you pass the argument "transaction" you have to convert the Transaction struct 
to string like that "dataToVerify := fmt.Sprintf("%x\n", transaction)" and then cast 
to array of bytes and pass as argument like that "script.Script.Evaluate([]byte(dataToVerify))"
*/
func VerifySignature(transaction ,pubkey, sig []byte) bool{
	curve := elliptic.P256()

	r := big.Int{}
	s := big.Int{}

	sigLen := len(sig)
	r.SetBytes(sig[:(sigLen / 2)])
	s.SetBytes(sig[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	dataToVerify := fmt.Sprintf("%x\n", transaction)

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
		return false
	}
	return true
}


func CheckSum(payload []byte) []byte{
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}