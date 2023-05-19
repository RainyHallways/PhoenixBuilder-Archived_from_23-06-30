package ResourcesControlCenter

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
)

// 列出所有可访问库存的窗口 ID ，即 WindowID
func (i *inventoryContents) ListWindowID() []uint32 {
	i.lockDown.RLock()
	defer i.lockDown.RUnlock()
	// init
	ans := []uint32{}
	for key := range i.datas {
		ans = append(ans, key)
	}
	// get window id list
	return ans
	// return
}

// 列出指定窗口 ID 所对应库存中的所有已记录槽位
func (i *inventoryContents) ListSlot(windowID uint32) ([]uint8, error) {
	i.lockDown.RLock()
	defer i.lockDown.RUnlock()
	// init
	got, ok := i.datas[windowID]
	if !ok {
		return []uint8{}, fmt.Errorf("ListSlot: %v is not recorded in i.datas; i.datas = %#v", windowID, i.datas)
	}
	// if windowsID is not exist
	ans := []uint8{}
	for key := range got {
		ans = append(ans, key)
	}
	// get slots list
	return ans, nil
	// return
}

// 获取 windowID 所对应的库存数据
func (i *inventoryContents) GetInventoryInfo(windowID uint32) (map[uint8]protocol.ItemInstance, error) {
	i.lockDown.RLock()
	defer i.lockDown.RUnlock()
	// init
	ans, ok := i.datas[windowID]
	if !ok {
		return map[uint8]protocol.ItemInstance{}, fmt.Errorf("GetInventoryInfo: %v is not recorded in i.datas; i.datas = %#v", windowID, i.datas)
	}
	// if windowsID is not exist
	return ans, nil
	// return
}

// 从 windowID 库存中获取 slotLocation 槽位的物品数据
func (i *inventoryContents) GetItemStackInfo(windowID uint32, slotLocation uint8) (protocol.ItemInstance, error) {
	i.lockDown.RLock()
	defer i.lockDown.RUnlock()
	// init
	got, ok := i.datas[windowID]
	if !ok {
		return protocol.ItemInstance{}, fmt.Errorf("GetItemStackInfo: %v is not recorded in i.datas; i.datas = %#v", windowID, i.datas)
	}
	// if windowsID is not exist
	ans, ok := got[slotLocation]
	if !ok {
		return protocol.ItemInstance{}, fmt.Errorf("GetItemStackInfo: %v is not recorded in i.datas[%v]; i.datas[%v] = %#v", slotLocation, windowID, windowID, i.datas[windowID])
	}
	// if slot is not exist
	return ans, nil
	// return
}

// 创建窗口 ID 为 windowID 的库存，如果库存不存在的话
func (i *inventoryContents) createNewInventory(windowID uint32) {
	i.lockDown.Lock()
	defer i.lockDown.Unlock()
	// init
	if i.datas == nil {
		i.datas = make(map[uint32]map[uint8]protocol.ItemInstance)
	}
	if i.datas[windowID] == nil {
		i.datas[windowID] = make(map[uint8]protocol.ItemInstance)
	}
	// create new inventory
}

// 修改 windowID 库存中 slotLocation 槽位的物品数据，属于私有实现
func (i *inventoryContents) writeItemStackInfo(windowID uint32, slotLocation uint8, itemStackInfo protocol.ItemInstance) {
	i.createNewInventory(windowID)
	// create new inventory if needed
	i.lockDown.Lock()
	defer i.lockDown.Unlock()
	// lock down resources
	i.datas[windowID][slotLocation] = itemStackInfo
	// write datas
}

// 删除 windowID 所对应的库存。
// 例如，当容器被关闭后，那么可以通过此函数删除此容器的库存数据，属于私有实现
func (i *inventoryContents) deleteInventory(windowID uint32) error {
	i.lockDown.Lock()
	defer i.lockDown.Unlock()
	// init
	_, ok := i.datas[windowID]
	if !ok {
		return fmt.Errorf("deleteInventory: %v is not recorded in i.datas; i.datas = %#v", windowID, i.datas)
	}
	// if windowID is not exist
	delete(i.datas, windowID)
	newMap := map[uint32]map[uint8]protocol.ItemInstance{}
	for key, value := range i.datas {
		newMap[key] = value
	}
	i.datas = newMap
	// remove inventory from i.datas
	return nil
	// return
}
