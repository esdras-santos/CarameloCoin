package wallet_test

import(
	"testing"
	"gochain/wallet"
)

var w = wallet.MakeWallet()

func TestPubKeySEC(t *testing.T){
	pub := w.PublicKey[:1]
	pubint := int(pub[0])
	if pubint != 2 {
		if pubint != 3{
			t.Error("invalid")
		}	
	} 
}