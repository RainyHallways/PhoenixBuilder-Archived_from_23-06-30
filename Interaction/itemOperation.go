package GlobalAPI

import (
	"encoding/gob"
	"fmt"
	"phoenixbuilder/ResourcesControlCenter"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 描述一个空气物品
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

// 描述单个物品所在的位置
type ItemLocation struct {
	WindowID    int16 // 物品所在库存的窗口 ID
	ContainerID uint8 // 物品所在库存的库存类型 ID
	Slot        uint8 // 物品所在的槽位
}

// 描述铁砧操作的操作结果
type AnvilOperationResponce struct {
	// 指代操作结果，为真时代表成功，否则反之
	SuccessStates bool
	// 指代被操作物品的最终位置，可能不存在。
	// 如果不存在，则代表物品已被丢出
	Destination *ItemLocation
}

// 将库存编号为 source 所指代的物品移动到 destination 所指代的槽位
// 且只移动 moveCount 个物品。
// details 指代相应槽位的预期变动结果，它将作为更新本地库存数据的依据。
// 当且仅当物品操作得到租赁服的响应后，此函数才会返回物品操作结果。
func (g *GlobalAPI) MoveItem(
	source ItemLocation,
	destination ItemLocation,
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
//
// 返回值的第一项代表物品名称的修改结果。
// 当发生错误时此参数将始终为 nil
//
// 当遭遇改名失败时，将尝试撤销名称修改请求。
// 如果原有物品栏已被占用，则会尝试将铁砧中的失败品
// 返还到背包中另外一个可用的物品栏。
// 如果背包中所有物品栏都已被占用，则物品会被留在铁砧内。
//
// 部分情况下此函数可能会遇见无法处理的错误，届时程序将抛出严重错误(panic)
func (g *GlobalAPI) ChangeItemName(
	name string,
	slot uint8,
) (*AnvilOperationResponce, error) {
	containerOpenDatas := g.Resources.Container.GetContainerOpenDatas()
	// 取得已打开的容器的数据
	if containerOpenDatas == nil {
		return nil, fmt.Errorf("ChangeItemName: Anvil has not opened")
	}
	// 如果铁砧未被打开
	var itemDatas protocol.ItemInstance
	get, err := g.Resources.Inventory.GetItemStackInfo(
		uint32(containerOpenDatas.WindowID),
		1,
	)
	if err != nil {
		return nil, fmt.Errorf("ChangeItemName: %v", err)
	}
	ResourcesControlCenter.DeepCopy(
		&get,
		&itemDatas,
		func() {
			gob.Register(map[string]interface{}{})
		},
	)
	// 取得已放入铁砧的物品的物品数据
	var backup protocol.ItemInstance
	ResourcesControlCenter.DeepCopy(
		&get,
		&backup,
		func() {
			gob.Register(map[string]interface{}{})
		},
	)
	// 备份物品数据
	filterAns, err := g.Resources.Inventory.ListSlot(0, &[]int32{0})
	if err != nil {
		panic(fmt.Sprintf("ChangeItemName: %v", err))
	}
	// 筛选出背包中还未被占用实际物品的物品栏
	optionalSlot := []uint8{slot}
	optionalSlot = append(optionalSlot, filterAns...)
	if len(filterAns) <= 0 {
		optionalSlot = []uint8{}
	}
	// optionalSlot 指代被操作物品最终可能出现的位置
	revertFunc := func() (*AnvilOperationResponce, error) {
		for _, value := range optionalSlot {
			placeStackRequestAction := protocol.PlaceStackRequestAction{}
			placeStackRequestAction.Source = protocol.StackRequestSlotInfo{
				ContainerID:    0,
				Slot:           1,
				StackNetworkID: backup.StackNetworkID,
			}
			placeStackRequestAction.Destination = protocol.StackRequestSlotInfo{
				ContainerID:    0xc,
				Slot:           value,
				StackNetworkID: 0,
			}
			placeStackRequestAction.Count = byte(backup.Stack.Count)
			// 构造一个新的 placeStackRequestAction 结构体
			resp, err := g.SendItemStackRequestWithResponce(
				&packet.ItemStackRequest{
					Requests: []protocol.ItemStackRequest{
						{
							RequestID: g.Resources.ItemStackOperation.GetNewRequestID(),
							Actions: []protocol.StackRequestAction{
								&placeStackRequestAction,
							},
						},
					},
				},
				[]ItemChangeDetails{
					{
						details: map[ResourcesControlCenter.ContainerID]ResourcesControlCenter.StackRequestContainerInfo{
							0x0: {
								WindowID: uint32(containerOpenDatas.WindowID),
								ChangeResult: map[uint8]protocol.ItemInstance{
									1: AirItem,
								},
							},
							0xc: {
								WindowID: 0,
								ChangeResult: map[uint8]protocol.ItemInstance{
									value: backup,
								},
							},
						},
					},
				},
			)
			// 尝试将被槽位物品还原到背包中的 value 物品栏处
			if err != nil {
				return nil, err
			}
			if resp[0].Status == protocol.ItemStackResponseStatusOK {
				return &AnvilOperationResponce{
					SuccessStates: false,
					Destination: &ItemLocation{
						WindowID:    0,
						ContainerID: 0xc,
						Slot:        value,
					},
				}, nil
			}
			// 如果成功还原的话，那么返回值
		}
		// 尝试把被操作物品从铁砧放回背包中
		return &AnvilOperationResponce{
			SuccessStates: false,
			Destination: &ItemLocation{
				WindowID:    int16(containerOpenDatas.WindowID),
				ContainerID: 0,
				Slot:        1,
			},
		}, nil
		// 看起来背包已经满了，我们不得不把物品留在铁砧中
	}
	// 构造一个函数用于处理改名失败时的善后处理
	for _, value := range optionalSlot {
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
			Slot:           value,
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
			return revertFunc()
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
						value: itemDatas,
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
			return revertFunc()
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
		if err != nil {
			return revertFunc()
		}
		// 取得物品名称操作结果
		if ans.Status == 0x9 {
			return revertFunc()
		}
		// 此时改名失败，原因是物品的新名称与原始名称重名
		if ans.Status == protocol.ItemStackResponseStatusOK {
			return &AnvilOperationResponce{
				SuccessStates: true,
				Destination: &ItemLocation{
					WindowID:    0,
					ContainerID: 0xc,
					Slot:        value,
				},
			}, nil
		}
		// 当改名成功时
	}
	return revertFunc()
	// 返回值
}

// 将 source 所指代的槽位中的全部物品丢出。
// windowID 指代被丢出物品所在库存的窗口 ID 。
// 返回值第一项代表丢出结果，
// 为真时代表成功丢出，否则反之
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
