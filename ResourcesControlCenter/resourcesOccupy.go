package ResourcesControlCenter

import (
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

// 描述一个通用的客户端资源独占结构
type resourcesOccupy struct {
	// 用于阻塞其他请求该资源的互斥锁
	lockDown sync.Mutex
	// 描述资源的占用状态，为 1 时代表已被占用，否则反之
	lockStates uint32
	// 标识资源的占用者，为 UUID 的字符串形式
	holder string
}

/*
占用客户端的某个资源。

当 tryMode 为真时，将尝试占用资源并返回占用结果，此对应返回值 bool 部分。
若 tryMode 为假，则此项返回真。

返回的字符串指代资源的占用者，这用于资源释放函数
func (c *container) Release(holder string) bool
中的 holder 参数
*/
func (r *resourcesOccupy) Occupy(tryMode bool) (bool, string) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return r.Occupy(tryMode)
	}
	uniqueId := newUUID.String()
	// get new unique id
	if tryMode {
		success := r.lockDown.TryLock()
		if !success {
			return false, ""
		}
		// if test failed
		r.holder = uniqueId
		return true, uniqueId
		// if success to lock
	}
	// if is try mode
	atomic.StoreUint32(&r.lockStates, 1)
	r.lockDown.Lock()
	// lock down resources
	r.holder = uniqueId
	// set the holder of this resources
	return true, uniqueId
	// return
}

// 释放客户端的某个资源，返回值代表执行结果。
// holder 指代该资源的占用者，当且仅当填写的占用者
// 可以与内部记录的占用者对应时才可以成功释放该资源
func (r *resourcesOccupy) Release(holder string) bool {
	if r.holder == holder && r.holder != "" {
		r.holder = ""
		atomic.StoreUint32(&r.lockStates, 0)
		r.lockDown.Unlock()
		return true
	}
	return false
}

// 返回资源的占用状态，为真时代表已被占用，否则反正
func (r *resourcesOccupy) GetOccupyStates() bool {
	if atomic.LoadUint32(&r.lockStates) == 0 {
		return false
	} else {
		return true
	}
}
