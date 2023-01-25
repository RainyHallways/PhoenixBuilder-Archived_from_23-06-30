package blockNBT

import (
	"fmt"
	blockNBT_depends "phoenixbuilder/fastbuilder/bdump/blockNBT/depends"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/mirror/chunk"
	"strings"
)

var ContainerIndexList map[string]string = map[string]string{
	"blast_furnace":      "Items",
	"lit_blast_furnace":  "Items",
	"smoker":             "Items",
	"lit_smoker":         "Items",
	"furnace":            "Items",
	"lit_furnace":        "Items",
	"chest":              "Items",
	"barrel":             "Items",
	"trapped_chest":      "Items",
	"lectern":            "book",
	"hopper":             "Items",
	"dispenser":          "Items",
	"dropper":            "Items",
	"cauldron":           "Items",
	"lava_cauldron":      "Items",
	"jukebox":            "RecordItem",
	"brewing_stand":      "Items",
	"undyed_shulker_box": "Items",
	"shulker_box":        "Items",
	// 以上都是现阶段支持了的容器
	"frame":      "Item",
	"glow_frame": "Item",
	// 物品展示框依赖于容器解析相关
}

type ContainerInput struct {
	ContainerData *map[string]interface{}
	Environment   *environment.PBEnvironment
	Mainsettings  *types.MainConfig
	BlockInfo     *types.Module
}

func Container(input *ContainerInput) error {
	if input.BlockInfo.StringNBT != nil {
		containerdata, err := getContainerDataRun(*input.ContainerData, *input.BlockInfo.Block.Name)
		if err != nil {
			return fmt.Errorf("Container: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
		}
		input.BlockInfo.ChestData = &containerdata
	}
	err := placeContainer(input.Environment, input.Mainsettings, input.BlockInfo)
	if err != nil {
		return fmt.Errorf("Container: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
	}
	return nil
}

func checkIfIsEffectiveContainer(name string) (string, error) {
	value, ok := ContainerIndexList[name]
	if ok {
		return value, nil
	}
	return "", fmt.Errorf("checkIfIsEffectiveContainer: \"%v\" not found", name)
}

// 将 Interface NBT 转换为 types.ChestData
func getContainerData(container interface{}) (types.ChestData, error) {
	var correct []interface{} = make([]interface{}, 0)
	// 初始化
	got, normal := container.([]interface{})
	if !normal {
		got, normal := container.(map[string]interface{})
		if !normal {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container; container = %#v", container)
		}
		correct = append(correct, got)
	} else {
		correct = got
	}
	// 把物品丢入 correct 里面
	// 如果这个物品是一个唱片机或者讲台，那么传入的 container 是一个 map[string]interface{} 而非 []interface{}
	// 为了更好的兼容性(更加方便)，这里都会把 map[string]interface{} 处理成通常情况下的 []interface{}
	// correct 就是处理结果
	ans := make(types.ChestData, 0)
	for key, value := range correct {
		var count uint8 = 0
		var itemData uint16 = 0
		var name string = ""
		var slot uint8 = 0
		var enchList []interface{}
		var item_lock uint8 = 0
		var keep_on_death uint8 = 0
		var can_place_on []string = []string{}
		var can_destroy []string = []string{}
		// 初始化
		containerData, normal := value.(map[string]interface{})
		if !normal {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v]; container[%v] = %#v", key, key, value)
		}
		// correct 这个列表中的每一项都必须是一个复合标签，也就得是 map[string]interface{} 才行
		_, ok := containerData["Count"]
		if !ok {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Count\"]; container[%v] = %#v", key, key, containerData)
		}
		count_got, normal := containerData["Count"].(byte)
		if !normal {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Count\"]; container[%v] = %#v", key, key, containerData)
		}
		count = uint8(count_got)
		// 拿一下物品数量
		// 它(物品数量)是一定存在的
		_, ok = containerData["Name"]
		if !ok {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Name\"]; container[%v] = %#v", key, key, containerData)
		}
		got, normal := containerData["Name"].(string)
		if !normal {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Name\"]; container[%v] = %#v", key, key, containerData)
		}
		name = strings.Replace(got, "minecraft:", "", 1)
		// 拿一下这个物品的物品名称
		// 可以看到，我这里是把命名空间删了的
		// 它(物品名称)是一定存在的
		_, ok = containerData["Damage"]
		if !ok {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Damage\"]; container[%v] = %#v", key, key, containerData)
		}
		damage_got, normal := containerData["Damage"].(int16)
		if !normal {
			return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Damage\"]; container[%v] = %#v", key, key, containerData)
		}
		itemData = uint16(damage_got)
		// 获取物品的 Damage 值
		// 这里的 Damage 值不一定就是物品的数据值(附加值)
		// 此处的 Damage 字段是一定存在的
		_, ok = containerData["tag"]
		if ok {
			tag, normal := containerData["tag"].(map[string]interface{})
			if !normal {
				return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"tag\"]; container[%v] = %#v", key, key, containerData)
			}
			// 这个 container["tag"] 一定是一个复合标签
			_, ok = tag["Damage"]
			if ok {
				got, normal := tag["Damage"].(int32)
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"tag\"][\"Damage\"]; containe[%v][\"tag\"] = %#v", key, key, tag)
				}
				itemData = uint16(got)
			}
			// container["tag"]["Damage"]
			_, ok = tag["minecraft:item_lock"]
			if ok {
				item_lock, normal = tag["minecraft:item_lock"].(byte)
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"tag\"][\"minecraft:item_lock\"]; container[%v][\"tag\"] = %#v", key, key, tag)
				}
			}
			// container["tag"]["minecraft:item_lock"]
			// 此数据为 1 时使用 lock_in_slot ，为 2 时使用 lock_in_inventory
			_, ok = tag["minecraft:keep_on_death"]
			if ok {
				keep_on_death, normal = tag["minecraft:keep_on_death"].(byte)
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"tag\"][\"minecraft:keep_on_death\"]; container[%v][\"tag\"] = %#v", key, key, tag)
				}
			}
			// container["tag"]["minecraft:keep_on_death"]
			_, ok = tag["ench"]
			if ok {
				enchList, normal = tag["ench"].([]interface{})
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"tag\"][\"ench\"]; container[%v][\"tag\"] = %#v", key, key, tag)
				}
			}
			// // container["tag"]["ench"]
		}
		// 拿一下这个工具的耐久值（当然也可能是别的，甚至它都不是个工具）及其他一些数据
		// 这个 tag 里的 Damage 实际上也不一定就是物品的数据值(附加值)
		// 需要说明的是，tag 不一定存在，且 tag 存在，Damage 也不一定存在
		_, ok = containerData["Block"]
		if ok {
			Block, normal := containerData["Block"].(map[string]interface{})
			if !normal {
				return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Block\"]; container[%v] = %#v", key, key, containerData)
			}
			// 这个 container["Block"] 一定是一个复合标签，如果不是就必须报错哦
			// 如果 Block 找得到则说明这个物品是一个方块
			_, ok = Block["val"]
			if !ok {
				_, ok := Block["states"]
				if !ok {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Block\"][\"states\"]; container[%v][\"Block\"] = %#v", key, key, Block)
				}
				got, normal := Block["states"].(map[string]interface{})
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Block\"][\"states\"]; container[%v][\"Block\"] = %#v", key, key, Block)
				}
				standardRuntimeID, found := chunk.StateToRuntimeID(name, got)
				if !found {
					itemData = 0
				} else {
					legacyBlock, found := chunk.RuntimeIDToLegacyBlock(standardRuntimeID)
					if !found {
						itemData = 0
					} else {
						itemData = legacyBlock.Val
					}
				}
			} else {
				got, normal := Block["val"].(int16)
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Block\"][\"val\"]; container[%v][\"Block\"] = %#v", key, key, Block)
				}
				itemData = uint16(got)
			}
			// 如果这个物品是个方块，也就是 Block 找得到的话
			// 那在 Block 里面一定有一个 val 去声明这个方块的方块数据值(附加值) [仅 MEMCBE]
		}
		// 拿一下这个方块的方块数据值(附加值)
		// 这个 Block 里的 val 一定是这个物品对应的方块的方块数据值(附加值)
		// 需要说明的是，Block 不一定存在，但如果 Block 存在，则 val 一定存在 [仅 NEMCBE]
		/*
			以上三个都在拿物品数据值(附加值)
			需要说明的是，数据值的获取优先级是这样的
			Damage < tag["Damage"] < Block["val"]
			需要说明的是，以上列举的三个情况不能涵盖所有的物品数据值(附加值)的情况，所以我希望可以有个人看一下普世情况是长什么样的，请帮帮我！
		*/
		_, ok = containerData["Slot"]
		if ok {
			got, normal := containerData["Slot"].(byte)
			if !normal {
				return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"Slot\"]; container[%v] = %#v", key, key, containerData)
			}
			slot = uint8(got)
		}
		// 拿一下这个物品所在的栏位(槽位)
		// 这个栏位(槽位)不一定存在，例如唱片机和讲台这种就不存在了(这种方块就一个物品，就不需要这个数据了)
		_, ok = containerData["CanPlaceOn"]
		if ok {
			got, normal := containerData["CanPlaceOn"].([]interface{})
			if !normal {
				return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"CanPlaceOn\"]; container[%v] = %#v", key, key, containerData)
			}
			for singleBlockIndex, singleBlock := range got {
				singleBlockString, normal := singleBlock.(string)
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"CanPlaceOn\"][%v]; container[%v][\"CanPlaceOn\"] = %#v", key, singleBlockIndex, key, got)
				}
				can_place_on = append(can_place_on, strings.Replace(singleBlockString, "minecraft:", "", 1))
			}
		}
		// 物品组件 - can_place_on
		// 此字段不一定存在
		_, ok = containerData["CanDestroy"]
		if ok {
			got, normal := containerData["CanDestroy"].([]interface{})
			if !normal {
				return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"CanDestroy\"]; container[%v] = %#v", key, key, containerData)
			}
			for singleBlockIndex, singleBlock := range got {
				singleBlockString, normal := singleBlock.(string)
				if !normal {
					return types.ChestData{}, fmt.Errorf("getContainerData: Crashed in container[%v][\"CanDestroy\"][%v]; container[%v][\"CanDestroy\"] = %#v", key, singleBlockIndex, key, got)
				}
				can_destroy = append(can_destroy, strings.Replace(singleBlockString, "minecraft:", "", 1))
			}
		}
		// 物品组件 - can_destroy
		// 此字段不一定存在
		ans = append(ans, types.ChestSlot{
			Name:        name,
			Count:       count,
			Damage:      itemData,
			Slot:        slot,
			ItemLock:    item_lock,
			KeepOnDeath: keep_on_death,
			CanPlaceOn:  can_place_on,
			CanDestroy:  can_destroy,
			EnchList:    enchList,
		})
		// 提交数据
	}
	return ans, nil
}

// 取得容器数据的主函数
func getContainerDataRun(blockNBT map[string]interface{}, blockName string) (types.ChestData, error) {
	key, err := checkIfIsEffectiveContainer(blockName)
	if err != nil {
		return types.ChestData{}, fmt.Errorf("getContainerDataRun: Not a supported container")
	}
	got, ok := blockNBT[key]
	// 这里是确定一下这个容器是否是我们支持了的容器
	if ok {
		ans, err := getContainerData(got)
		if err != nil {
			return types.ChestData{}, fmt.Errorf("getContainerDataRun: %v", err)
		}
		return ans, nil
	}
	// 如果这是个容器且对应的 key 可以找到，那么就去拿一下对应的 ContainerData 结构体
	return types.ChestData{}, nil
	// 对于唱片机和讲台这种容器，如果它们没有被放物品的话，那么对应的 key 是找不到的
	// 但是这不是一个错误，所以我们返回一个空的 ContainerData 和一个空的 error
}

func placeContainer(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	BlockInfo *types.Module,
) error {
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	// prepare
	if BlockInfo.Block != nil {
		request := commands_generator.SetBlockRequest(BlockInfo, Mainsettings)
		_, err := cmdsender.SendWSCommandWithResponce(request)
		if err != nil {
			return fmt.Errorf("placeContainer: %v", err)
		}
	}
	// place block
	got := commands_generator.ReplaceItemRequest(BlockInfo, Mainsettings)
	// get replaceitem command
	if BlockInfo.ChestSlot != nil {
		if len(got) > 0 {
			err := cmdsender.SendDimensionalCommand(got[0])
			if err != nil {
				return fmt.Errorf("placeContainer: %v", err)
			}
		}
		return nil
	}
	// for old method
	if BlockInfo.ChestData != nil {
		enchItemList := blockNBT_depends.EnchItemList{}
		for key, value := range got {
			chestData := *BlockInfo.ChestData
			_, ok := blockNBT_depends.ContainerIdIndexMap[*BlockInfo.Block.Name]
			if len(chestData[key].EnchList) <= 0 || !ok || !blockNBT_depends.CheckVersion() {
				err := cmdsender.SendDimensionalCommand(value)
				if err != nil {
					return fmt.Errorf("placeContainer: %v", err)
				}
			} else {
				enchItemList = append(enchItemList, blockNBT_depends.EnchItem{
					WantPutItem:   chestData[key],
					ContainerInfo: BlockInfo,
				})
			}
		}
		if len(enchItemList) > 0 {
			err := blockNBT_depends.PutItemIntoContainerRun(Environment, enchItemList)
			if err != nil {
				return fmt.Errorf("placeContainer: %v", err)
			}
		}
	}
	// for new method
	return nil
}
