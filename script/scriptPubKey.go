package script

import "gochain/wallet"

func P2pkhScript(w *wallet.Wallet) []byte {
	script := []byte{0x76}
	script = append(script, 0xa9)
	hash := wallet.PublicKeyHash(w.PublicKey)
	script = append(script, []byte{byte(len(hash))}...)
	script = append(script, hash...)
	script = append(script, 0x88)
	script = append(script, 0xac)
	return script
}