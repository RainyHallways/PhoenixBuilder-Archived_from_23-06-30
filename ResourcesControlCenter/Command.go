package ResourcesControlCenter

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"

	"github.com/google/uuid"
)

// 测定 key 是否在 c.commandRequest.datas 中。如果存在，那么返回真，否则返回假
func (c *commandRequestWithResponce) TestRequest(key uuid.UUID) bool {
	c.commandRequest.lockDown.RLock()
	defer c.commandRequest.lockDown.RUnlock()
	// init
	_, ok := c.commandRequest.datas[key]
	return ok
	// return
}

// 测定 key 是否在 c.commandResponce.datas 中。如果存在，那么返回真，否则返回假
func (c *commandRequestWithResponce) TestResponce(key uuid.UUID) bool {
	c.commandResponce.lockDown.RLock()
	defer c.commandResponce.lockDown.RUnlock()
	// init
	_, ok := c.commandResponce.datas[key]
	return ok
	// return
}

// 将名为 key 的命令请求放入 c.commandRequest.datas 并占用(锁定)此请求对应的互斥锁
func (c *commandRequestWithResponce) WriteRequest(key uuid.UUID) error {
	if c.TestRequest(key) {
		return fmt.Errorf("WriteRequest: %v is already exist in c.commandRequest.datas", key.String())
	}
	// if key is already exist
	c.commandRequest.lockDown.Lock()
	defer c.commandRequest.lockDown.Unlock()
	// lock down resources
	c.commandRequest.datas[key] = &sync.Mutex{}
	c.commandRequest.datas[key].Lock()
	// lock down command request
	return nil
	// return
}

// 将名为 key 的命令请求从 c.commandRequest.datas 中移除并释放此请求对应的互斥锁，
// 属于私有实现
func (c *commandRequestWithResponce) deleteRequest(key uuid.UUID) error {
	if !c.TestRequest(key) {
		return fmt.Errorf("deleteRequest: %v is not recorded in c.commandRequest.datas", key.String())
	}
	// if key is not exist
	c.commandRequest.lockDown.Lock()
	defer c.commandRequest.lockDown.Unlock()
	// lock down resources
	tmp := c.commandRequest.datas[key]
	// get tmp of the current resources
	delete(c.commandRequest.datas, key)
	newMap := map[uuid.UUID]*sync.Mutex{}
	for k, value := range c.commandRequest.datas {
		newMap[k] = value
	}
	c.commandRequest.datas = newMap
	// remove the key and values from c.commandRequest.datas
	tmp.Unlock()
	// unlock command request
	return nil
	// return
}

// 将命令请求的返回值写入 c.commandResponce.datas
// 并释放 c.commandRequest.datas 中对应的互斥锁，属于私有实现
func (c *commandRequestWithResponce) writeResponce(key uuid.UUID, resp packet.CommandOutput) error {
	c.commandResponce.lockDown.Lock()
	defer c.commandResponce.lockDown.Unlock()
	// init
	c.commandResponce.datas[key] = resp
	// send command responce
	err := c.deleteRequest(key)
	if err != nil {
		return fmt.Errorf("writeResponce: %v", err)
	}
	// remove command reuqest from c.commandRequest.datas
	return nil
	// return
}

// 从 c.commandResponce.datas 读取名为 key 的命令请求的返回值
// 并将此返回值从 c.commandResponce.datas 移除
func (c *commandRequestWithResponce) LoadResponceAndDelete(key uuid.UUID) (packet.CommandOutput, error) {
	if !c.TestResponce(key) {
		return packet.CommandOutput{}, fmt.Errorf("loadResponceAndDelete: %v is not recorded in c.commandResponce.datas", key.String())
	}
	// if key is not exist
	c.commandResponce.lockDown.Lock()
	defer c.commandResponce.lockDown.Unlock()
	// lock down resources
	ans := c.commandResponce.datas[key]
	newMap := map[uuid.UUID]packet.CommandOutput{}
	for k, value := range c.commandResponce.datas {
		newMap[k] = value
	}
	c.commandResponce.datas = newMap
	// get responce and remove the key and values from c.commandResponce.datas
	return ans, nil
	// return
}

// 等待租赁服响应命令请求 key 。
// 在调用此函数后，会持续阻塞直到此命令请求所对应的互斥锁被释放
func (c *commandRequestWithResponce) AwaitResponce(key uuid.UUID) {
	if !c.TestRequest(key) {
		return
	}
	// if key is not exist
	c.commandRequest.lockDown.RLock()
	// lock down resources
	tmp := c.commandRequest.datas[key]
	// get tmp of the current resources
	c.commandRequest.lockDown.RUnlock()
	// unlock resources
	tmp.Lock()
	tmp.Unlock()
	// await responce
}
