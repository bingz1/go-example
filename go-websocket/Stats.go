package main

type Stats struct {
	//反馈在线长连接的数量
	OnlineConnections int64 `json:"online_connections"`

	//反馈客户端的推送压力
	SendMessageTotal int64 `json:"send_message_total""`
	SendMessageFail int64 `json:"send_message_fail"`
	
	//反馈ConnMgr消息分发模块的压力
	DispatchPending int64 `json:"dispatch_pending"`
	PushJobPending int64 `json:"push_job_pending"`
	DispactchFail int64 `json:"dispactch_fail"`
	
	//返回出在线的房间总数，有利于分析内存暴涨的原因
	RoomCount int64 `json:"room_count"`

	//Merger模块
}