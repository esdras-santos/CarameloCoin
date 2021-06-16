package network

import (
	"gochain/utils"
	"math"
	"math/big"
	"math/rand"
	"time"
)


type VersionMessage struct {
	Version          []byte
	Services         []byte
	Timestamp        []byte
	ReceiverServices []byte
	ReceiverIP       []byte
	ReceiverPort     []byte
	SenderServices   []byte
	SenderIp         []byte
	SenderPort       []byte
	Nonce            []byte
	UserAgent        []byte
	LatestBlock      []byte
	//relay must have 01 to true or 00 to false
	Relay []byte
}

func (vm *VersionMessage) Init(version, services, timestamp, receiverservices, receiverip, receiverport, senderservices, senderip, senderport, nonce, useragent, latestblock []byte, relay bool) {
	vm.Version = version
	vm.Services = services
	if timestamp == nil {
		vm.Timestamp = utils.ToHex(time.Now().Unix())
	} else {
		vm.Timestamp = timestamp
	}
	vm.ReceiverServices = receiverservices
	vm.ReceiverIP = receiverip
	vm.ReceiverPort = receiverport
	vm.SenderServices = senderservices
	vm.SenderIp = senderip
	vm.SenderPort = senderport
	if nonce == nil {
		vm.Nonce = utils.ToLittleEndian(utils.ToHex(int64(rand.Intn(int(math.Pow(2, 64))))), 8)
	} else {
		vm.Nonce = nonce
	}
	vm.UserAgent = useragent
	vm.LatestBlock = latestblock
	if relay {
		vm.Relay = []byte{0x01}
	} else {
		vm.Relay = []byte{0x00}
	}
}

func (vm *VersionMessage) Serialize() []byte {
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
	utils.EncodeVarint(*big.NewInt(int64(len(vm.UserAgent))), &buf)
	result = append(result, buf...)
	result = append(result, vm.UserAgent...)
	result = append(result, utils.ToLittleEndian(vm.LatestBlock, 4)...)
	result = append(result, vm.Relay...)
	return result
}
