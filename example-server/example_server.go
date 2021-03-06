package main

import (
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	// 1、创建客户端连接
	s := modsynclog_v1.NewServer("127.0.0.1", 2222, 1235, "/tmp/test-modsynlog-server")
	s.SocketServer.Set("SendFlag", 123456)

	// 2、直接发一条数据过去
	s.CommitData(1234, []byte("haha12345678"))

	// 处理其他逻辑
	select {}
}
