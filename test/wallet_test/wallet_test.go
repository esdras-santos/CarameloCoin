package wallet_test

import (
	"bytes"
	"gochain/wallet"
	"testing"
)

var w = wallet.MakeWallet()

func TestAddress(t *testing.T){
	pubHash := w.PublicKeyHash()
	versionedHash := append([]byte{byte(0x00)},pubHash...)
	checksum := wallet.CheckSum(versionedHash)
	fullHash := append(versionedHash,checksum...)
	address := wallet.Base58Encode(fullHash)

	if !bytes.Equal(w.Address(),address){
		t.Error("wrong address")
	}
}

func TestAddressToPublicKeyHash(t *testing.T){
	if !bytes.Equal(w.PublicKeyHash(),wallet.AddressToPKH(string(w.Address()))){
		t.Error("wrong public key hash")
	}
}

func TestValidateAddress(t * testing.T){
	if !wallet.ValidateAddress(string(w.Address())){
		t.Error("invalid address")
	}

	if wallet.ValidateAddress("1abcdef123456789kjfghdsdfg"){
		t.Error("valid address")
	}
}

func TestLoadWalletFile(t *testing.T){
	w.SaveFile("test","./tmp/wallet_test.data")
	lw := wallet.Wallet{}
	lw.LoadFile("test","./tmp/wallet_test.data")
	if  string(lw.Address()) != string(w.Address()){
		t.Error("load error")
	}
}