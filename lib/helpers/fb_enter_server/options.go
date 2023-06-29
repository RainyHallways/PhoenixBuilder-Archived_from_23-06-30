package fb_enter_server

import "time"

type Options struct {
	AuthServer                 string
	FBServerConnectionTimeOut  time.Duration
	FBGetTokenTimeOut          time.Duration
	MCServerConnectionTimeOut  time.Duration
	MCServerConnectRetryTimes  int
	FBUserName                 string
	FBUserPassword             string
	FBUserToken                string
	ServerCode                 string
	ServerPassword             string
	WriteBackToken             bool
	ExpectedCmdFeedBack        bool
	PrintUQHolderDebugInfo     bool
	TransferTimeOut            time.Duration
	TransferCheckNumTimeOut    time.Duration
	MakeBotCreative            bool
	DisableCommandBlock        bool
	MaximumWaitTime            time.Duration
	DeadOnOpPrivilegeRemoved   bool
	OpPrivilegeRemovedCallBack func()
}

func MakeDefaultOption() *Options {
	return &Options{
		AuthServer:                 "",
		FBServerConnectionTimeOut:  time.Minute,
		FBGetTokenTimeOut:          time.Minute,
		MCServerConnectionTimeOut:  time.Minute,
		MCServerConnectRetryTimes:  0,
		FBUserName:                 "",
		FBUserPassword:             "",
		FBUserToken:                "",
		ServerCode:                 "",
		ServerPassword:             "",
		WriteBackToken:             true,
		ExpectedCmdFeedBack:        false,
		PrintUQHolderDebugInfo:     false,
		TransferTimeOut:            time.Minute,
		TransferCheckNumTimeOut:    time.Minute,
		MakeBotCreative:            true,
		DisableCommandBlock:        true,
		MaximumWaitTime:            time.Minute * 3,
		DeadOnOpPrivilegeRemoved:   true,
		OpPrivilegeRemovedCallBack: nil,
	}
}
