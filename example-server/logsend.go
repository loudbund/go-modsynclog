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

// 6、主函数 -------------------------------------------------------------------------
func main() {
	// 应用日志
	// sendLogApp()
	// 数据库日志
	sendLogDb()
}

// 2、日志发送：app日志
func sendLogApp() {
	handleAppLog := modsynclog_v1.NewSdkAppLog("http://127.0.0.1:1234", "127.0.0.1:1235")

	if true {
		D := make([]*grpc_proto_log.AppLogData, 0)
		D = append(D, &grpc_proto_log.AppLogData{
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
		})
		D = append(D, &grpc_proto_log.AppLogData{
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
		})
		// 1、http方式发送日志
		if err := handleAppLog.SdkAppLogAddHttp(D); err != nil {
			fmt.Println(err)
		}
	}

	// GRPC日志 //////////////////////
	if true {
		for i := 0; i < 2; i++ {
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

// 3、日志发送：数据库日志 //////////////////////////////////////
func sendLogDb() {
	handleDbLog := modsynclog_v1.NewSdkDbLog("http://127.0.0.1:1234", "127.0.0.1:1235")

	if false {
		D := make([]*grpc_proto_log.DbLogData, 0)
		D = append(D, &grpc_proto_log.DbLogData{
			DbInstance: "anoah",
			Type:       "table-create",
			Database:   "wawa",
			Table:      "haha",
			Ts:         "1234",
			Position:   "xxx",
			Xid:        1,
			Commit:     true,
			Data:       nil,
			Sql:        "create table",
		})
		D = append(D, &grpc_proto_log.DbLogData{
			DbInstance: "anoah",
			Type:       "insert",
			Database:   "wawa",
			Table:      "haha",
			Ts:         "1234",
			Position:   "xxx",
			Xid:        1,
			Commit:     true,
			Sql:        "",
			Data: map[string]string{
				"id": "123",
			},
		})
		// 1、http方式发送日志
		if err := handleDbLog.SdkDbLogAddHttp(D); err != nil {
			fmt.Println(err)
		}
	}

	// GRPC日志 //////////////////////
	if true {
		for i := 0; i < 2; i++ {
			D := &grpc_proto_log.DbLogData{
				DbInstance: "anoah",
				Type:       "insert",
				Database:   "wawa",
				Table:      "haha",
				Ts:         "1234",
				Position:   "xxx",
				Xid:        1,
				Commit:     true,
				Sql:        "",
				Data: map[string]string{
					"id": "123",
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
