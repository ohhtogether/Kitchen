package controller

// RespData 响应给客户端的数据结构
type RespData struct {
	//Code  int64       `json:"code"`  // 错误消息
	Error string      `json:"error"` // 错误消息
	Data  interface{} `json:"data"`  // 正常的数据
}
