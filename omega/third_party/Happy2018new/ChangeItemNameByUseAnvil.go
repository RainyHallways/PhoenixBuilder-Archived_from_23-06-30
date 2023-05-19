package Happy2018new

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	GlobalAPI "phoenixbuilder/Interaction"
	"phoenixbuilder/fastbuilder/mcstructure"
	"phoenixbuilder/minecraft/nbt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

type ChangeItemNameByUseAnvil struct {
	*defines.BasicComponent
	apis     GlobalAPI.GlobalAPI
	Triggers []string `json:"菜单触发词"`
	Usage    string   `json:"菜单项描述"`
	FilePath string   `json:"从何处提取物品的新名称(填写路径)"`
}

func (o *ChangeItemNameByUseAnvil) Init(settings *defines.ComponentConfig, storage defines.StorageAndLogProvider) {
	marshal, _ := json.Marshal(settings.Configs)
	if err := json.Unmarshal(marshal, o); err != nil {
		panic(err)
	}
}

func (o *ChangeItemNameByUseAnvil) Inject(frame defines.MainFrame) {
	o.Frame = frame
	o.apis = GlobalAPI.GlobalAPI{
		WritePacket: func(p packet.Packet) error {
			o.Frame.GetGameControl().SendMCPacket(p)
			return nil
		},
		BotInfo: GlobalAPI.BotInfo{
			BotName:      o.Frame.GetUQHolder().GetBotName(),
			BotIdentity:  "",
			BotUniqueID:  o.Frame.GetUQHolder().BotUniqueID,
			BotRunTimeID: o.Frame.GetUQHolder().BotRuntimeID,
		},
		Resources: o.Frame.GetResources(),
	}
	o.Frame.GetGameListener().SetGameMenuEntry(&defines.GameMenuEntry{
		MenuEntry: defines.MenuEntry{
			Triggers:     o.Triggers,
			ArgumentHint: "[快捷栏槽位: int] [x: int] [y: int] [z: int]",
			FinalTrigger: false,
			Usage:        o.Usage,
		},
		OptionalOnTriggerFn: o.ChangeItemName,
	})
}

func (o *ChangeItemNameByUseAnvil) ChangeItemName(chat *defines.GameChat) bool {
	go func() {
		var mode uint8 = 0
		var targetSlot uint8 = 0
		var readPos []int32 = []int32{}
		var itemName string = ""
		o.apis.BotInfo.BotName = o.Frame.GetUQHolder().GetBotName()
		// 初始化
		if len(chat.Msg) > 0 {
			got, err := strconv.ParseUint(chat.Msg[0], 10, 32)
			if err != nil {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c无法解析槽位数据§f，§c请确认你提供了正确的整数\n详细日志已发送到控制台")
				pterm.Error.Printf("修改物品名称: %v\n", err)
				return
			}
			targetSlot = uint8(got)
			if targetSlot > 8 {
				o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c你提供的槽位参数 §b%v §c已大于 §b8", targetSlot))
				return
			}
		} else {
			o.Frame.GetGameControl().SayTo(chat.Name, "§e你没有提供槽位参数§f，§e现在默认重定向为 §b0")
		}
		// 确定槽位位置
		if len(chat.Msg) > 1 {
			mode = 1
			if len(chat.Msg) < 4 {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c提供的参数不足§f，§c当前缺少一个或多个坐标")
				return
			}
			for i := 0; i < 3; i++ {
				got, err := strconv.ParseInt(chat.Msg[i+1], 10, 32)
				if err != nil {
					o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c无法解析坐标数据§f，§c错误发生在位置 §b%v §f，§c请确认你提供了正确的坐标数据\n详细日志已发送到控制台", i))
					pterm.Error.Printf("修改物品名称: %v\n", err)
					return
				}
				readPos = append(readPos, int32(got))
			}
		}
		// 如果用户希望在游戏内完成名称编辑操作
		if mode == 0 {
			datas, err := o.Frame.GetFileData(o.FilePath)
			if err != nil {
				o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c无法打开 §bomega_storage/data/%v §c处的文件\n详细日志已发送到控制台", o.FilePath))
				pterm.Error.Printf("修改物品名称: %v\n", err)
				return
			}
			if len(datas) <= 0 {
				o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§bomega_storage/data/%v §c处的文件没有填写物品名称§f，§c可能这个文件是个空文件§f，§c也可能是文件本身不存在", o.FilePath))
				return
			}
			itemName = strings.ReplaceAll(string(datas), "\r", "")
		} else if mode == 1 {
			_, holder := o.apis.Resources.Structure.Occupy(false)
			resp, err := o.apis.SendStructureRequestWithResponce(
				&packet.StructureTemplateDataRequest{
					StructureName: "Omega:ChangeItemNameByUseAnvil",
					Position:      protocol.BlockPos{readPos[0], readPos[1], readPos[2]},
					Settings: protocol.StructureSettings{
						PaletteName:               "default",
						IgnoreEntities:            true,
						IgnoreBlocks:              false,
						Size:                      protocol.BlockPos{1, 1, 1},
						Offset:                    protocol.BlockPos{0, 0, 0},
						LastEditingPlayerUniqueID: o.apis.BotInfo.BotUniqueID,
						Rotation:                  0,
						Mirror:                    0,
						Integrity:                 100,
						Seed:                      0,
						AllowNonTickingChunks:     false,
					},
					RequestType: packet.StructureTemplateRequestExportFromSave,
				},
			)
			o.apis.Resources.Structure.Release(holder)
			if err != nil {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c未能请求命令方块数据\n详细日志已发送到控制台")
				pterm.Error.Printf("修改物品名称: %v\n", err)
				return
			}
			// 请求结构数据
			_, reversedMap, _ := mcstructure.SplitArea(
				mcstructure.BlockPos{readPos[0], readPos[1], readPos[2]},
				mcstructure.BlockPos{readPos[0], readPos[1], readPos[2]},
				64,
				64,
				true,
			)
			got, err := mcstructure.GetMCStructureData(
				mcstructure.Area{
					BeginX: readPos[0],
					BeginY: readPos[1],
					BeginZ: readPos[2],
					SizeX:  1,
					SizeY:  1,
					SizeZ:  1,
				},
				resp.StructureTemplate,
			)
			if err != nil {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c未能请求命令方块数据\n详细日志已发送到控制台")
				pterm.Error.Printf("修改物品名称: %v\n", err)
				return
			}
			allAreas := []mcstructure.Mcstructure{got}
			processedData, err := mcstructure.DumpBlocks(
				allAreas,
				reversedMap,
				mcstructure.Area{
					BeginX: int32(readPos[0]),
					BeginY: int32(readPos[1]),
					BeginZ: int32(readPos[2]),
					SizeX:  1,
					SizeY:  1,
					SizeZ:  1,
				},
			)
			if err != nil {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c未能请求命令方块数据\n详细日志已发送到控制台")
				pterm.Error.Printf("修改物品名称: resp = %#v\n", resp)
				return
			}
			if len(processedData) <= 0 {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c未能请求命令方块数据\n详细日志已发送到控制台")
				pterm.Error.Printf("修改物品名称: resp = %#v\n", resp)
				return
			}
			// 从结构中提取命令方块数据
			newBuffer := bytes.NewBuffer(processedData[0].NBTData)
			var commandBlockNBT map[string]interface{}
			err = nbt.NewDecoderWithEncoding(newBuffer, nbt.LittleEndian).Decode(&commandBlockNBT)
			if err != nil {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c未能请求命令方块数据\n详细日志已发送到控制台")
				pterm.Error.Printf("修改物品名称: processedData[0].NBTData = %#v\n", processedData[0].NBTData)
				return
			}
			_, ok := commandBlockNBT["Command"]
			if !ok {
				o.Frame.GetGameControl().SayTo(chat.Name, "§c目标方块不是命令方块")
				return
			}
			itemName, _ = commandBlockNBT["Command"].(string)
		}
		// 获取物品的新名称
		itemDatas, err := o.apis.Resources.Inventory.GetItemStackInfo(0, targetSlot)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c在读取快捷栏 §b%v §c时发送了错误\n详细日志已发送到控制台", targetSlot))
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if itemDatas.Stack.ItemType.NetworkID == 0 {
			o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c请确保机器人在快捷栏 §b%v §c有一个物品\n详细日志已发送到控制台", targetSlot))
			pterm.Warning.Printf("修改物品名称: itemDatas = %#v\n", itemDatas)
			return
		}
		// 确定被改名物品存在
		cmdResp, err := o.apis.SendWSCommandWithResponce("querytarget @s")
		if err != nil {
			panic(pterm.Error.Sprintf("修改物品名称: %v", err))
		}
		parseAns, err := o.apis.ParseQuerytargetInfo(cmdResp)
		if err != nil {
			panic(pterm.Error.Sprintf("修改物品名称: %v", err))
		}
		if len(parseAns) <= 0 {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c机器人可能没有 §bOP §c权限")
			return
		}
		pos := [3]int32{
			int32(math.Floor(float64(parseAns[0].Position[0]))),
			int32(math.Floor(float64(parseAns[0].Position[1]))),
			int32(math.Floor(float64(parseAns[0].Position[2]))),
		}
		// 取得机器人当前的坐标
		successStates, err := o.apis.ChangeItemNameByUsingAnvil(
			pos,
			`["direction": 0, "damage": "undamaged"]`,
			[]GlobalAPI.AnvilChangeItemName{
				{
					Slot: targetSlot,
					Name: itemName,
				},
			},
		)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c物品名称修改失败\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if !successStates[0] {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c物品名称修改失败§f，§c请检查新的名称是否与原始名称相同或该物品是否可被移动")
			return
		}
		// 修改物品名称
		newItemDatas, err := o.apis.Resources.Inventory.GetItemStackInfo(0, targetSlot)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c在读取快捷栏 §b%v §c时发送了错误\n详细日志已发送到控制台", targetSlot))
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		// 读取新物品的数据
		dropResp, err := o.apis.DropItemAll(
			protocol.StackRequestSlotInfo{
				ContainerID:    28,
				Slot:           targetSlot,
				StackNetworkID: newItemDatas.StackNetworkID,
			},
			0,
		)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c尝试丢出新物品时失败\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if !dropResp {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c尝试丢出新物品时失败\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: dropResp = %#v\n", dropResp)
			return
		}
		// 丢出新物品
		o.Frame.GetGameControl().SayTo(chat.Name, "§a已成功修改物品名称")
		// 返回值
	}()
	return true
}
