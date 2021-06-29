package script


func P2pkhScript(pubKeyHash []byte) []byte {
	script := []byte{0x76}
	script = append(script, 0xa9)
	script = append(script, []byte{byte(len(pubKeyHash))}...)
	script = append(script, pubKeyHash...)
	script = append(script, 0x88)
	script = append(script, 0xac)
	return script
}