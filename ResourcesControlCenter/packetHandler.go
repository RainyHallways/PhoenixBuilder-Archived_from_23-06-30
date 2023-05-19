package ResourcesControlCenter

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 根据收到的数据包更新客户端的资源数据
func (r *Resources) handlePacket(pk *packet.Packet) {
	switch p := (*pk).(type) {
	case *packet.CommandOutput:
		uniqueId := p.CommandOrigin.UUID
		ok := r.Command.TestRequest(uniqueId)
		if !ok {
			return
		}
		r.Command.writeResponce(uniqueId, *p)
		// send ws command with responce
	case *packet.InventoryContent:
		for key, value := range p.Content {
			if value.Stack.ItemType.NetworkID != -1 {
				r.Inventory.writeItemStackInfo(p.WindowID, uint8(key), value)
			}
		}
		// inventory contents(global)
	case *packet.InventoryTransaction:
		for _, value := range p.Actions {
			if value.SourceType == protocol.InventoryActionSourceCreative {
				continue
			}
			r.Inventory.writeItemStackInfo(uint32(value.WindowID), uint8(value.InventorySlot), value.NewItem)
		}
		// inventory contents(for enchant command...)
	case *packet.ItemStackResponse:
		for _, value := range p.Responses {
			r.ItemStackOperation.updateItemData(value, &r.Inventory)
			// update local inventory datas
			err := r.ItemStackOperation.writeResponce(value.RequestID, value)
			if err != nil {
				panic("handlePacket: Attempt to send packet.ItemStackRequest without using ResourcesControlCenter")
			}
			// write responce
		}
		// item stack request
	case *packet.ContainerOpen:
		unsuccess, _ := r.Container.Occupy(true)
		if unsuccess {
			panic("handlePacket: Attempt to send packet.ContainerOpen without using ResourcesControlCenter")
		}
		r.Container.writeContainerCloseDatas(nil)
		r.Container.writeContainerOpenDatas(p)
		r.Inventory.createNewInventory(uint32(p.WindowID))
		r.Container.responceContainerOperation()
		// while open a container
	case *packet.ContainerClose:
		if p.WindowID != 0 && p.WindowID != 119 && p.WindowID != 120 && p.WindowID != 124 {
			err := r.Inventory.deleteInventory(uint32(p.WindowID))
			if err != nil {
				panic(fmt.Sprintf("handlePacket: Try to removed an inventory which not existed; p.WindowID = %v", p.WindowID))
			}
		}
		if !p.ServerSide && !r.Container.GetOccupyStates() {
			panic("handlePacket: Attempt to send packet.ContainerClose without using ResourcesControlCenter")
		}
		r.Container.writeContainerOpenDatas(nil)
		r.Container.writeContainerCloseDatas(p)
		r.Container.responceContainerOperation()
		// while a container is closed
	case *packet.StructureTemplateDataResponse:
		if r.Structure.GetOccupyStates() {
			r.Structure.writeStructureResponce(*p)
		}
		// packet.StructureTemplateDataRequest
	}
}
