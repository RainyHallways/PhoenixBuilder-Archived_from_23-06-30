package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"

	"github.com/google/uuid"
)

// 向租赁服发送 WS 命令且获取返回值
func (g *GlobalAPI) SendWSCommandWithResponce(command string) (packet.CommandOutput, error) {
	uniqueId, err := uuid.NewUUID()
	if err != nil || uniqueId == uuid.Nil {
		resp, err := g.SendWSCommandWithResponce(command)
		if err != nil {
			return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
		}
		return resp, nil
	}
	uniqueIdString := uniqueId.String()
	// 初始化
	CommandRequest.Store(uniqueIdString, uint8(0))
	// 写入请求到等待队列
	err = g.SendWSCommand(command, uniqueId)
	if err != nil {
		return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
	}
	// 发送命令
	for {
		got, success := CommandResponce.LoadAndDelete(uniqueIdString)
		if success {
			val, normal := got.(packet.CommandOutput)
			if !normal {
				return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: Responce is not a packet.CommandOutput struct")
			}
			return val, nil
		}
	}
}
