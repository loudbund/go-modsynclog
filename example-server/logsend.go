package main

import (
	"fmt"
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
	"github.com/loudbund/go-modsynclog/modsynclog_v1/grpc_proto_log"
	"github.com/loudbund/go-utils/utils_v1"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetReportCaller(true)
}

// 主函数 -------------------------------------------------------------------------
func main() {
	// 数据库日志
	sendLogDb()
	// 应用日志
	sendLogApp()
	// json字串格式日志
	sendLogJson()
}

// 1、日志发送：数据库日志 //////////////////////////////////////
func sendLogDb() {
	// handleDbLog := modsynclog_v1.NewSdkDbLog("http://127.0.0.1:1234", "127.0.0.1:1235")
	handleDbLog := modsynclog_v1.NewSdkDbLog("127.0.0.1:1235")

	// GRPC日志 //////////////////////
	if true {
		for i := 0; i < 1; i++ {

			D := &grpc_proto_log.DbLogData{
				DbInstance: "dblog",
				Type:       "insert",
				Database:   "test",
				Table:      "haha",
				Ts:         "1234",
				Position:   "xxx",
				Xid:        1,
				Commit:     true,
				Sql:        "",
				Data: map[string]string{
					"id": "2",
				},
			}

			// grpc发送日志
			if err := handleDbLog.SdkDbLogAddGRpc(D); err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 5)
			} else {
				fmt.Println("Ok")
			}
		}
	}
}

// 2、日志发送：app日志
func sendLogApp() {
	// handleAppLog := modsynclog_v1.NewSdkAppLog("http://127.0.0.1:1234", "127.0.0.1:1235")
	handleAppLog := modsynclog_v1.NewSdkAppLog("127.0.0.1:1235")

	// GRPC日志 //////////////////////
	if true {
		for i := 0; i < 1; i++ {
			D := &grpc_proto_log.AppLogData{
				Env:       "dev",
				Sys:       "haha",
				Level:     "info",
				File:      "abc.go",
				Func:      "haha()",
				Time:      utils_v1.Time().DateTime(),
				TimeInt64: 0,
				Message:   "hahaha",
				Data: map[string]string{
					"techerId": "500",
				},
			}

			// grpc发送日志
			if err := handleAppLog.SdkAppLogAddGRpc(D); err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 5)
			} else {
				fmt.Println("Ok")
			}
		}
	}

}

// 3、日志发送：json字串格式日志 //////////////////////////////////////
func sendLogJson() {
	// handleDbLog := modsynclog_v1.NewSdkDbLog("http://127.0.0.1:1234", "127.0.0.1:1235")
	handleLog := modsynclog_v1.NewSdkJsonLog("127.0.0.1:1235")

	// GRPC日志 //////////////////////
	if true {
		for i := 0; i < 1; i++ {
			D := map[string]interface{}{
				"wawa": "haha",
			}
			// grpc发送日志
			if err := handleLog.SdkJsonLogAddGRpc(D); err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 5)
			} else {
				fmt.Println("Ok")
			}
		}
	}
}
