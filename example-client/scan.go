package main

import (
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
	"github.com/loudbund/go-modsynclog/modsysclog_read_logs"
	"time"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	readLog := modsysclog_read_logs.NewReadLogs("/tmp/test-modsynlog-client")
	if err := readLog.Read(time.Now().Format("2006-01-02"), 0, func(Date string, DataId int64, Data *filelog_v1.UDataSend, Finish bool) int {
		// 读取结束 或 读取到空数据
		if Finish || Data == nil {
			return 0
		}
		// 打印点数据
		fmt.Println(Date, DataId, string(Data.Data), Finish)
		return 1
	}); err != nil {
		fmt.Println(err)
	}
}
