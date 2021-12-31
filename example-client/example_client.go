package main

import (
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	serverIp := "127.0.0.1"
	serverPort := 2222

	// 创建客户端连接
	_ = modsynclog_v1.NewClient(serverIp, serverPort, "/tmp/test-synlog-client")

	// 处理其他逻辑
	select {}
}
