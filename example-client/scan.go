package main

import (
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
	"github.com/loudbund/go-utils/utils_v1"
	"time"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	scanFileLog("/tmp/test-modsynlog-client", time.Now().Format("2006-01-02"), func(D *filelog_v1.UDataSend) error {
		fmt.Println(D.Date, D.Id, D.DataType, string(D.Data))
		return nil
	})
}

// 消费日志文件demo
func scanFileLog(folder string, date string, dataEvent func(D *filelog_v1.UDataSend) error) {
	var (
		KeyLogFolder = folder
		KeyDate      = date
		KeyHandle    *filelog_v1.CFileLog
		KeyId        int64
	)
	KeyHandle = filelog_v1.New(KeyLogFolder, KeyDate)
	KeyId = 0
	// 循环处理
	for {
		// 读取一条数据
		if D, err := KeyHandle.GetOne(KeyId); err != nil { // 读取出错
			fmt.Println(err)
		} else if D == nil { // 读取到nil数据，可能指定日期的日志已消费完成，或者是当天的还没有产生新数据
			if KeyHandle.GetFinish(true) {
				// 已经读完并且当前日志不是今天，则切换到下一天读取
				if KeyDate < utils_v1.Time().Date() {
					KeyHandle.Close()
					KeyDate = utils_v1.Time().DateAdd(KeyDate, 1)
					KeyHandle = filelog_v1.New(KeyLogFolder, KeyDate)
					KeyId = 0
					fmt.Println("准备读取:", KeyDate, KeyId)
				}
			}
			// 延时1秒
			time.Sleep(time.Second)
		} else { // 读取到nil数据
			if err := dataEvent(D); err != nil {
				// 延时5秒
				time.Sleep(time.Second * 5)
			} else {
				KeyId++
			}
		}
	}
}
