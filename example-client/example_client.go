package main

import (
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
	"time"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	// 1、创建客户端连接
	// 日志服务的ip为127.0.0.1；端口号为2222；日志同步到目录"/tmp/test-modsynlog-client"
	_ = modsynclog_v1.NewClient(
		"127.0.0.1",
		2222,
		"/tmp/test-modsynlog-client",
		modsynclog_v1.ClientOptions{ // 省略掉该参数，InitHistoryDayNum为0，即只补充同步当天日志
			InitHistoryDayNum:   -1,     // 补充同步昨天的日志
			SendFlag:            123456, // 传输码
			RetainHistoryDayNum: -3,     // 日志保留3天
		},
	)

	// 2、模拟消费日志数据
	logHandle := filelog_v1.New("/tmp/test-modsynlog-client", time.Now().Format("2006-01-02"))
	index := int64(0)                                  // 读取位置
	if D, err := logHandle.GetOne(index); err != nil { // 1、读取出错
		fmt.Println("模拟读取数据：", err.Error())
	} else if D == nil { // 2、读到了结尾
		fmt.Println("未获取到第0条数据", D)
	} else { // 读到了数据
		fmt.Println("模拟读取数据：", D.DataType, string(D.Data))
	}
	logHandle.Close()

	// 处理其他逻辑
	select {}
}
