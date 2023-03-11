package blockNBT_API

import (
	"sync"
)

// ------------------------- sendwscmdwithresp.go -------------------------

// 存放命令请求的等待列
var CommandRequest sync.Map = sync.Map{}

// 存放命令请求的返回值
var CommandResponce sync.Map = sync.Map{}

// ------------------------- end -------------------------
