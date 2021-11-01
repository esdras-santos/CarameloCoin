package wallet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"os"

	"golang.org/x/crypto/ripemd160"
)

const(
	checksumLength = 4
	version = byte(0x00)
	//walletFile = "./tmp/wallet.data"
)

type Wallet struct{
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w Wallet) Address() []byte{
	pubHash := w.PublicKeyHash()
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

func PktoPKH(pubkey []byte) []byte{
	pubHash := sha256.Sum256(pubkey)
	
	hasher := ripemd160.New()
	_,err := hasher.Write(pubHash[:])
	Handle(err)
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

func PKHtoAddress(pkh []byte) string{
	versionedHash := append([]byte{version},pkh...)
	checksum := CheckSum(versionedHash)
	fullHash := append(versionedHash,checksum...)
	address := Base58Encode(fullHash)
	return string(address)
}

func ValidateAddress(address string) bool{
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-checksumLength]
	targetChecksum := CheckSum(append([]byte{version},pubKeyHash...))
	return bytes.Compare(actualChecksum,targetChecksum) == 0
}

func NewKeyPair() (ecdsa.PrivateKey,[]byte){
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve,rand.Reader)
	Handle(err)
	
	pub := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pub
}



func MakeWallet() *Wallet{
	private,public := NewKeyPair()
	wallet := Wallet{private,public}

	return &wallet
}

func (w *Wallet) PublicKeyHash() []byte{
	pubHash := sha256.Sum256(w.PublicKey)
	
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

func CheckSum(payload []byte) []byte{
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w *Wallet) LoadFile(walletFile string) error{
	if _,err := os.Stat(walletFile);os.IsNotExist(err){
		return err
	}

	var wall Wallet

	fileContent, err := ioutil.ReadFile(walletFile)
	Handle(err)
	
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wall)
	Handle(err)

	w.PrivateKey = wall.PrivateKey
	w.PublicKey = wall.PublicKey
	return nil
}

func (w *Wallet) SaveFile(walletFile string) {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	Handle(err)
	

	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	Handle(err)
}

func keyAdjust(password string) []byte{
	key := []byte(password)
	if len(key) < 16{
		for i := len(key);i<16;i++{
			key = append(key, 0x00)
		}
	}else if len(key) < 24{
		for i := len(key);i<24;i++{
			key = append(key, 0x00)
		}
	}else if len(key) < 32{
		for i := len(key);i<32;i++{
			key = append(key, 0x00)
		}
	}
	return key
}

//encrypt the wallet
func encrypt(key []byte,data []byte) []byte{


    // generate a new aes cipher using our 32 byte long key
    c, err := aes.NewCipher(key)
    Handle(err)

    // gcm or Galois/Counter Mode, is a mode of operation
    // for symmetric key cryptographic block ciphers
    // - https://en.wikipedia.org/wiki/Galois/Counter_Mode
    gcm, err := cipher.NewGCM(c)
    // if any error generating new GCM
    // handle them
    Handle(err)

    // creates a new byte array the size of the nonce
    // which must be passed to Seal
    nonce := make([]byte, gcm.NonceSize())
    // populates our nonce with a cryptographically secure
    // random sequence
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        fmt.Println(err)
    }

    
    return gcm.Seal(nonce, nonce, data, nil)

	
}

//decrypt the wallet
func decrypt(key []byte,data []byte) []byte{    

    c, err := aes.NewCipher(key)
    Handle(err)

    gcm, err := cipher.NewGCM(c)
    Handle(err)

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        fmt.Println(err)
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    Handle(err)
    return plaintext
}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}