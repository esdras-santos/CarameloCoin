package network

type VerAckMessage struct{
	Command []byte
}

func (vam VerAckMessage) GetCommand() []byte{
	return vam.Command
}

func (vam VerAckMessage) Serialize() []byte{
	return nil
}