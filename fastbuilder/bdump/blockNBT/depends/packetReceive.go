package blockNBT_depends

import (
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft/protocol/packet"
)

var Islistening bool
var ReceivePacket []packet.Packet

func PacketProcessor(Environment *environment.PBEnvironment, NeedWaiting bool, RequestPacketId uint32) []packet.Packet {
	if NeedWaiting {
		cmdsender := Environment.CommandSender.(*commands.CommandSender)
		cmdsender.SendWSCommandWithResponce("list")
	}
	// waiting for the packet
	ans := []packet.Packet{}
	for _, value := range ReceivePacket {
		if value.ID() == RequestPacketId {
			ans = append(ans, value)
		}
	}
	// get datas
	closeProcessor()
	// close processor
	return ans
}

func InitProcessor() {
	Islistening = true
	ReceivePacket = []packet.Packet{}
}

func closeProcessor() {
	Islistening = false
	ReceivePacket = []packet.Packet{}
}
