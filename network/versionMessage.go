package network

import (
	"gochain/blockchain"
	"gochain/utils"
	"math"
	"math/rand"
	"time"
)


type VersionMessage struct {
	Command			 []byte
	Version          []byte //4 bytes
	Services         []byte //8 bytes le
	Timestamp        []byte //8 bytes le
	ReceiverServices []byte	//8 bytes le
	ReceiverIP       []byte //4 bytes 
	ReceiverPort     []byte //2 bytes
	SenderServices   []byte //8 bytes le
	SenderIp         []byte //4 bytes
	SenderPort       []byte //2 bytes
	Nonce            []byte //8 bytes
	UserAgent        []byte //dinamic length
	LatestBlock      []byte	//4 bytes
	//relay must have 01 to true or 00 to false
	Relay []byte // 1 byte
}

func (vm *VersionMessage) Init(timestamp, receiverip, nonce, relay []byte) {
	var chain = blockchain.BlockChain{}
	var bestHeight = chain.GetBestHeight()

	vm.Command = []byte("version")
	vm.Version = []byte{0x00000001}
	vm.Services = []byte{0x00000000000000}
	if timestamp == nil {
		vm.Timestamp = utils.ToHex(time.Now().Unix())
	} else {
		vm.Timestamp = timestamp
	}
	vm.ReceiverServices = []byte{0x00000000000000}
	vm.ReceiverIP = receiverip
	vm.ReceiverPort = []byte(PORT)
	vm.SenderServices = []byte{0x00000000000000}
	vm.SenderIp = AddressToBytes(NODEIP)
	vm.SenderPort = []byte(PORT)
	if nonce == nil {
		vm.Nonce = utils.ToLittleEndian(utils.ToHex(int64(rand.Intn(int(math.Pow(2, 64))))), 8)
	} else {
		vm.Nonce = []byte{0x6e,0x6f,0x74,0x20,0x6d,0x65,0x21,0x21}
	}
	vm.UserAgent = []byte("/CarameloCoin:0.1/")
	vm.LatestBlock = utils.ToHex(bestHeight)
	vm.Relay = []byte{0x01}
}

//will bug if the is IPv6
func (vm *VersionMessage) Parse(data []byte) {
	vm.Version = utils.ToLittleEndian(data[:4],4)
	vm.Services = utils.ToLittleEndian(data[4:12],8)
	vm.Timestamp = utils.ToLittleEndian(data[12:20],8)
	vm.ReceiverServices = utils.ToLittleEndian(data[20:28],8)
	vm.ReceiverIP = data[40:44]
	vm.ReceiverPort = utils.ToLittleEndian(data[44:46],2)
	vm.SenderServices = utils.ToLittleEndian(data[46:54],8)
	vm.SenderIp = utils.ToLittleEndian(data[66:70],4)
	vm.SenderPort = utils.ToLittleEndian(data[70:72],2)
	vm.Nonce = data[72:80]
	var len int
	utils.ReadVarint(data[80:],&len)
	var startIn int
	if len <= 253{
		startIn = 81
	}else if len <= 254{
		startIn = 82
	}else if len <= 255{
		startIn = 83
	}
	sl := startIn+len
	vm.UserAgent = data[startIn:sl]
	vm.LatestBlock = utils.ToLittleEndian(data[sl:sl+4],4)
	vm.Relay = data[sl+4:]
}

func (vm VersionMessage) GetCommand() []byte{
	return vm.Command
}

func (vm VersionMessage) Serialize() []byte {
	result := utils.ToLittleEndian(vm.Version, 4)
	result = append(result, utils.ToLittleEndian(vm.Services, 8)...)
	result = append(result, utils.ToLittleEndian(vm.Timestamp, 8)...)
	result = append(result, utils.ToLittleEndian(vm.ReceiverServices, 8)...)
	result = append(result, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}...)
	result = append(result, vm.ReceiverIP...)
	result = append(result, utils.ToLittleEndian(vm.ReceiverPort, 2)...)
	result = append(result, utils.ToLittleEndian(vm.SenderServices, 8)...)
	result = append(result, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}...)
	result = append(result, vm.SenderIp...)
	result = append(result, utils.ToLittleEndian(vm.SenderPort, 2)...)
	result = append(result, vm.Nonce...)
	buf := []byte{}
	utils.EncodeVarint(int64(len(vm.UserAgent)), &buf)
	result = append(result, buf...)
	result = append(result, vm.UserAgent...)
	result = append(result, utils.ToLittleEndian(vm.LatestBlock, 4)...)
	result = append(result, vm.Relay...)
	return result
}
