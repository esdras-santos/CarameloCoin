package network


//when you start a new node this message will be sended
//when that message is sended your node will receive a "blockchain" message
type GetBlockMessage struct{
	Command []byte
}

func (gb *GetBlockMessage) Init(){
	gb.Command = []byte("getblock")
}

func (gb GetBlockMessage) GetCommand() []byte{
	return gb.Command
}

func (gb GetBlockMessage) Serialize() []byte{
	
	return nil
}

func (gb *GetBlockMessage) Parse(data []byte){
	
}

