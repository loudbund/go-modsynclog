package main

import (
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	serverIp := "127.0.0.1"
	portSocket := 2222
	portHttp := 1234
	portGRpc := 1235

	// 创建客户端连接
	_ = modsynclog_v1.NewServer(serverIp, portSocket, portHttp, portGRpc, "/tmp/test-modsynlog-server")

	// 处理其他逻辑
	select {}
}
