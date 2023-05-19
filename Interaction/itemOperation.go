package GlobalAPI

import (
	"fmt"
	"phoenixbuilder/ResourcesControlCenter"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 空物品数据使用它来描述
var AirItem protocol.ItemInstance = protocol.ItemInstance{
	StackNetworkID: 0,
	Stack: protocol.ItemStack{
		ItemType: protocol.ItemType{
			NetworkID:     0,
			MetadataValue: 0,
		},
		BlockRuntimeID: 0,
		Count:          0,
		NBTData:        map[string]interface{}(nil),
		CanBePlacedOn:  []string(nil),
		CanBreak:       []string(nil),
		HasNetworkID:   false,
	},
}

// 在描述物品移动操作时所使用的结构体
type MoveItemDatas struct {
	WindowID    int16 // 物品所在库存的窗口 ID
	ContainerID uint8 // 物品所在库存的库存类型 ID
	Slot        uint8 // 物品所在的槽位
}

// 将库存编号为 source 所指代的物品移动到 destination 所指代的槽位
// 且只移动 moveCount 个物品。
// details 指代相应槽位的预期变动结果，它将作为更新本地库存数据的依据。
// 当且仅当物品操作得到租赁服的响应后，此函数才会返回物品操作结果。
func (g *GlobalAPI) MoveItem(
	source MoveItemDatas,
	destination MoveItemDatas,
	details ItemChangeDetails,
	moveCount uint8,
) ([]protocol.ItemStackResponse, error) {
	placeStackRequestAction := protocol.PlaceStackRequestAction{}
	// 初始化
	itemOnSource, err := g.Resources.Inventory.GetItemStackInfo(uint32(source.WindowID), source.Slot)
	if err != nil {
		return []protocol.ItemStackResponse{}, fmt.Errorf("MoveItem: %v", err)
	}
	itemOnDestination, _ := g.Resources.Inventory.GetItemStackInfo(uint32(destination.WindowID), destination.Slot)
	// 取得 source 和 destination 处的物品信息
	if moveCount <= uint8(itemOnSource.Stack.Count) || source.WindowID == -1 {
		placeStackRequestAction.Count = moveCount
	} else {
		placeStackRequestAction.Count = uint8(itemOnSource.Stack.Count)
	}
	// 得到欲移动的物品数量
	placeStackRequestAction.Source = protocol.StackRequestSlotInfo{
		ContainerID:    source.ContainerID,
		Slot:           source.Slot,
		StackNetworkID: itemOnSource.StackNetworkID,
	}
	placeStackRequestAction.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    destination.ContainerID,
		Slot:           destination.Slot,
		StackNetworkID: itemOnDestination.StackNetworkID,
	}
	// 构造 placeStackRequestAction 结构体
	ans, err := g.SendItemStackRequestWithResponce(
		&packet.ItemStackRequest{
			Requests: []protocol.ItemStackRequest{
				{
					Actions: []protocol.StackRequestAction{
						&placeStackRequestAction,
					},
					FilterStrings: []string{},
				},
			},
		},
		[]ItemChangeDetails{details},
	)
	if err != nil {
		return []protocol.ItemStackResponse{}, fmt.Errorf("MoveItem: %v", err)
	}
	// 发送物品操作请求
	return ans, nil
	// 返回值
}

// 将已放入铁砧第一格(注意是第一格)的物品的物品名称修改为 name 并返还到背包中的 slot 处。
// 当且仅当租赁服回应操作结果后此函数再返回值。
// 部分情况下此函数可能会遇见无法处理的错误，届时程序将抛出严重错误(panic)
func (g *GlobalAPI) ChangeItemName(
	name string,
	slot uint8,
) (bool, error) {
	containerOpenDatas := g.Resources.Container.GetContainerOpenDatas()
	if containerOpenDatas == nil {
		return false, fmt.Errorf("ChangeItemName: Anvil has not opened")
	}
	// 如果铁砧未被打开
	itemDatas, err := g.Resources.Inventory.GetItemStackInfo(
		uint32(containerOpenDatas.WindowID),
		1,
	)
	if err != nil {
		return false, fmt.Errorf("ChangeItemName: %v", err)
	}
	// 取得已放入铁砧的物品的物品数据
	revertFunc := func() error {
		source := MoveItemDatas{
			WindowID:    int16(containerOpenDatas.WindowID),
			ContainerID: 0,
			Slot:        1,
		}
		destination := MoveItemDatas{
			WindowID:    0,
			ContainerID: 0xc,
			Slot:        slot,
		}
		ans, err := g.MoveItem(
			source,
			destination,
			ItemChangeDetails{
				map[ResourcesControlCenter.ContainerID]ResourcesControlCenter.StackRequestContainerInfo{
					0: {
						WindowID: uint32(containerOpenDatas.WindowID),
						ChangeResult: map[uint8]protocol.ItemInstance{
							1: AirItem,
						},
					},
					0xc: {
						WindowID: 0,
						ChangeResult: map[uint8]protocol.ItemInstance{
							slot: itemDatas,
						},
					},
				},
			},
			uint8(itemDatas.Stack.Count),
		)
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
		if ans[0].Status != protocol.ItemStackResponseStatusOK {
			return fmt.Errorf("ChangeItemName: Could not revert operation because of the new operation which numbered %v have been canceled by error code %v. This maybe is a BUG, please provide this logs to the developers!\nnewAns = %#v; source = %#v; destination = %#v; moveCount = %v", ans[0].RequestID, ans[0].Status, ans, source, destination, itemDatas.Stack.Count)
		}
		return nil
	}
	// 构造一个用于错误恢复的函数
	newRequestID := g.Resources.ItemStackOperation.GetNewRequestID()
	// 请求一个新的 RequestID 用于 ItemStackRequest
	placeStackRequestAction := protocol.PlaceStackRequestAction{}
	placeStackRequestAction.Count = byte(itemDatas.Stack.Count)
	placeStackRequestAction.Source = protocol.StackRequestSlotInfo{
		ContainerID:    0x3c,
		Slot:           0x32,
		StackNetworkID: newRequestID,
	}
	placeStackRequestAction.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    0xc,
		Slot:           slot,
		StackNetworkID: 0,
	}
	// 构造一个新的 PlaceStackRequestAction 结构体
	newItemStackRequest := packet.ItemStackRequest{
		Requests: []protocol.ItemStackRequest{
			{
				RequestID: newRequestID,
				Actions: []protocol.StackRequestAction{
					&protocol.CraftRecipeOptionalStackRequestAction{
						RecipeNetworkID:   0,
						FilterStringIndex: 0,
					},
					&protocol.ConsumeStackRequestAction{
						DestroyStackRequestAction: protocol.DestroyStackRequestAction{
							Count: byte(itemDatas.Stack.Count),
							Source: protocol.StackRequestSlotInfo{
								ContainerID:    0,
								Slot:           1,
								StackNetworkID: itemDatas.StackNetworkID,
							},
						},
					},
					&placeStackRequestAction,
				},
				FilterStrings: []string{name},
			},
		},
	}
	// 构造一个新的 ItemStackRequest 结构体
	err = g.Resources.ItemStackOperation.SetItemName(
		&itemDatas,
		name,
	)
	if err != nil {
		err = revertFunc()
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
		return false, nil
	}
	// 更新物品数据中的名称字段以用于更新本地库存数据
	_, ok := itemDatas.Stack.NBTData["RepairCost"]
	if !ok {
		itemDatas.Stack.NBTData["RepairCost"] = int32(0)
	}
	// 更新物品数据中的 RepairCost 字段用于更新本地库存数据
	err = g.Resources.ItemStackOperation.WriteRequest(
		newRequestID,
		map[ResourcesControlCenter.ContainerID]ResourcesControlCenter.StackRequestContainerInfo{
			0xc: {
				WindowID: 0,
				ChangeResult: map[uint8]protocol.ItemInstance{
					slot: itemDatas,
				},
			},
			0x0: {
				WindowID: uint32(containerOpenDatas.WindowID),
				ChangeResult: map[uint8]protocol.ItemInstance{
					1: AirItem,
				},
			},
			0x1: {
				WindowID: uint32(containerOpenDatas.WindowID),
				ChangeResult: map[uint8]protocol.ItemInstance{
					2: AirItem,
				},
			},
			0x3c: {
				WindowID: uint32(containerOpenDatas.WindowID),
				ChangeResult: map[uint8]protocol.ItemInstance{
					0x32: AirItem,
				},
			},
		},
	)
	if err != nil {
		err = revertFunc()
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
	}
	// 写入请求到等待队列
	err = g.WritePacket(&newItemStackRequest)
	if err != nil {
		panic(fmt.Sprintf("ChangeItemName: %v", err))
	}
	// 发送物品操作请求
	g.Resources.ItemStackOperation.AwaitResponce(newRequestID)
	// 等待租赁服响应物品操作请求
	ans, err := g.Resources.ItemStackOperation.LoadResponceAndDelete(newRequestID)
	if err != nil || ans.Status != protocol.ItemStackResponseStatusOK {
		err = revertFunc()
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
	}
	// 当改名失败时尝试将物品恢复到背包中对应的位置
	if ans.Status != protocol.ItemStackResponseStatusOK {
		return false, nil
	}
	// 如果名称未发生变化或者因为其他一些原因所导致的改名失败
	return true, nil
	// 返回值
}

// 将 source 所指代的槽位中的全部物品丢出。
// windowID 指代被丢出物品所在库存的窗口 ID
func (g *GlobalAPI) DropItemAll(
	source protocol.StackRequestSlotInfo,
	windowID uint32,
) (bool, error) {
	ans, err := g.SendItemStackRequestWithResponce(
		&packet.ItemStackRequest{
			Requests: []protocol.ItemStackRequest{
				{
					Actions: []protocol.StackRequestAction{
						&protocol.DropStackRequestAction{
							Count:    64,
							Source:   source,
							Randomly: false,
						},
					},
				},
			},
		},
		[]ItemChangeDetails{
			{
				map[ResourcesControlCenter.ContainerID]ResourcesControlCenter.StackRequestContainerInfo{
					ResourcesControlCenter.ContainerID(source.ContainerID): {
						WindowID: windowID,
						ChangeResult: map[uint8]protocol.ItemInstance{
							source.Slot: AirItem,
						},
					},
				},
			},
		},
	)
	if err != nil {
		return false, fmt.Errorf("DropItemAll: %v", err)
	}
	if ans[0].Status != protocol.ItemStackResponseStatusOK {
		return false, nil
	}
	// 发送物品丢掷请求
	return true, nil
	// 返回值
}
