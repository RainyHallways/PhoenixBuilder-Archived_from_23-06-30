package blockNBT_depends

import (
	"fmt"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/mcstructure"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/mirror/chunk"
	"strings"
	"time"
)

var containerOpenData interface{}
var container_Hotbar_0_StackNetworkID int32

type EnchItem struct {
	WantPutItem   types.ChestSlot
	ContainerInfo *types.Module
}

type EnchItemList []EnchItem

func openContainer(
	Environment *environment.PBEnvironment,
	MainHandItemInfo *types.ChestSlot,
	ContainerBlockName *string,
	ContainerBlockStates *string,
	ContainerPos [3]int32,
) error {
	var blockName string
	if !strings.Contains(*ContainerBlockName, "minecraft:") {
		blockName = fmt.Sprintf("minecraft:%v", *ContainerBlockName)
	} else {
		blockName = *ContainerBlockName
	}
	// prepare
	containerOpenData = nil
	InitProcessor()
	// prepare
	got, err := mcstructure.ParseStringNBT(*ContainerBlockStates, true)
	if err != nil {
		return fmt.Errorf("openContainer: Failed to get block states; states = %#v", ContainerBlockStates)
	}
	containerBlockStates, normal := got.(map[string]interface{})
	if !normal {
		return fmt.Errorf("openContainer: Failed to converse ContainerBlockStates to map[string]interface{}; ContainerBlockStates = %#v", ContainerBlockStates)
	}
	standardRuntimeID, found := chunk.StateToRuntimeID(blockName, containerBlockStates) // 我相信你一定找得到的
	if !found {
		return fmt.Errorf("openContainer: Failed to get the runtimeID of this container which named %v; ContainerBlockStates = %#v", *ContainerBlockName, containerBlockStates)
	}
	blockRuntimeID := chunk.StandardRuntimeIDToNEMCRuntimeID(standardRuntimeID)
	if blockRuntimeID == chunk.AirRID || blockRuntimeID == chunk.NEMCAirRID {
		return fmt.Errorf("openContainer: Failed to converse StandardRuntimeID to NEMCRuntimeID; StandardRuntimeID = %#v, ContainerBlockName = %#v, ContainerBlockStates = %#v", standardRuntimeID, *ContainerBlockName, containerBlockStates)
	}
	networkID, ok := ItemRunTimeID[MainHandItemInfo.Name]
	// get blockRunTimeId and networkId
	if ok {
		Environment.Connection.(*minecraft.Conn).WritePacket(&packet.InventoryTransaction{
			LegacyRequestID:    0,
			LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
			Actions:            []protocol.InventoryAction{},
			TransactionData: &protocol.UseItemTransactionData{
				LegacyRequestID:    0,
				LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
				Actions:            []protocol.InventoryAction(nil),
				ActionType:         0,
				BlockPosition:      protocol.BlockPos{ContainerPos[0], ContainerPos[1], ContainerPos[2]},
				HotBarSlot:         0,
				HeldItem: protocol.ItemInstance{
					StackNetworkID: 0,
					Stack: protocol.ItemStack{
						ItemType: protocol.ItemType{
							NetworkID:     int32(networkID),
							MetadataValue: uint32(MainHandItemInfo.Damage),
						},
						BlockRuntimeID: 0,
						Count:          uint16(MainHandItemInfo.Count),
						CanBePlacedOn:  MainHandItemInfo.CanPlaceOn,
						CanBreak:       MainHandItemInfo.CanDestroy,
						HasNetworkID:   false,
					},
				},
				BlockRuntimeID: blockRuntimeID,
			},
		})
	}
	// open container
	containerOpenData = PacketProcessor(Environment, true, packet.IDContainerOpen)[0].(*packet.ContainerOpen)
	// process packet
	return nil
}

func closeContainer(
	Environment *environment.PBEnvironment,
	WindowID byte,
) {
	Environment.Connection.(*minecraft.Conn).WritePacket(&packet.ContainerClose{
		WindowID:   WindowID,
		ServerSide: false,
	})
}

func requestStackNetworkID(
	Environment *environment.PBEnvironment,
) error {
	container_Hotbar_0_StackNetworkID = 0
	InitProcessor()
	// prepare
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	_, err := cmdsender.SendWSCommandWithResponce("replaceitem entity @s slot.hotbar 1 air")
	if err != nil {
		return fmt.Errorf("requestStackNetworkID: %v", err)
	}
	got := PacketProcessor(Environment, false, packet.IDInventoryContent)
	for _, value := range got {
		i := value.(*packet.InventoryContent)
		if i.WindowID == 0 {
			container_Hotbar_0_StackNetworkID = i.Content[0].StackNetworkID
			break
		}
	}
	return nil
}

func putItemIntoContainer(
	Environment *environment.PBEnvironment,
	ItemInfo *types.ChestSlot,
	ContainerID byte,
) error {
	err := requestStackNetworkID(Environment)
	if err != nil {
		return fmt.Errorf("putItemIntoContainer: %v", err)
	}
	PlaceStackRequestAction := protocol.PlaceStackRequestAction{}
	PlaceStackRequestAction.Count = ItemInfo.Count
	PlaceStackRequestAction.Source = protocol.StackRequestSlotInfo{
		ContainerID:    ContainerCombinedHotBarAndInventory,
		Slot:           0,
		StackNetworkID: container_Hotbar_0_StackNetworkID,
	}
	PlaceStackRequestAction.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    ContainerID,
		Slot:           ItemInfo.Slot,
		StackNetworkID: 0,
	}
	request := &packet.ItemStackRequest{
		Requests: []protocol.ItemStackRequest{
			{
				RequestID: -1,
				Actions: []protocol.StackRequestAction{
					&PlaceStackRequestAction,
				},
				FilterStrings: []string{},
			},
		},
	}
	Environment.Connection.(*minecraft.Conn).WritePacket(request)
	return nil
}

/*
将 EnchItemList 列表中的物品放入物品，主要是为了支持附魔性物品

需要占用槽位 slot.hotbar 0 和 slot.hotbar 1

支持的容器请见当前文件夹下 container.go 中的 ContainerIdIndexMap
*/
func PutItemIntoContainerRun(
	Environment *environment.PBEnvironment,
	Input EnchItemList,
) error {
	var cod packet.ContainerOpen
	err := ReplaceitemAndEnchant(Environment, &Input[0].WantPutItem)
	if err != nil {
		return fmt.Errorf("PutItemIntoContainerRun: %v", err)
	}
	err = openContainer(Environment, &Input[0].WantPutItem, Input[0].ContainerInfo.Block.Name, &Input[0].ContainerInfo.Block.BlockStates, [3]int32{int32(Input[0].ContainerInfo.Point.X), int32(Input[0].ContainerInfo.Point.Y), int32(Input[0].ContainerInfo.Point.Z)})
	if err != nil {
		return fmt.Errorf("PutItemIntoContainerRun: %v", err)
	}
	failedCount := 0
	for {
		if failedCount > 100 {
			return fmt.Errorf("PutItemIntoContainerRun: Failed to open the container, please check the target area is loaded")
		}
		if containerOpenData != nil {
			cod = *containerOpenData.(*packet.ContainerOpen)
			break
		}
		failedCount++
		time.Sleep(50 * time.Millisecond)
	}
	for key, value := range Input {
		if key > 0 {
			err := ReplaceitemAndEnchant(Environment, &value.WantPutItem)
			if err != nil {
				return fmt.Errorf("PutItemIntoContainerRun: %v", err)
			}
		}
		containerID := ContainerIdIndexMap[*value.ContainerInfo.Block.Name]
		err = putItemIntoContainer(Environment, &value.WantPutItem, byte(containerID))
		if err != nil {
			return fmt.Errorf("PutItemIntoContainerRun: %v", err)
		}
	}
	closeContainer(Environment, cod.WindowID)
	return nil
}
