package main

import (
	"fmt"
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
	"github.com/loudbund/go-modsynclog/modsynclog_v1/grpc_proto_applog"
	"github.com/loudbund/go-utils/utils_v1"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
}

// 6、主函数 -------------------------------------------------------------------------
func main() {
	modsynclog_v1.SdkInitHttpSet("http://127.0.0.1:1234")
	modsynclog_v1.SdkInitGRpcAppLog("127.0.0.1:1235")

	// for i := 0; i < 11; i++ {
	// 	if err := modsynclog_v1.SendLog(123, []string{utils_v1.Time().DateTime()}); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	httpLog()
	httpAppLog()

	gRpcLog()
	for i := 0; i < 2; i++ {
		gRpcAppLog()
	}
}

// http写常规日志
func httpLog() {
}

// http写应用日志
func httpAppLog() {
	D := make([]*grpc_proto_applog.AppLogData, 0)
	D = append(D, &grpc_proto_applog.AppLogData{
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
	D = append(D, &grpc_proto_applog.AppLogData{
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
	if err := modsynclog_v1.SdkHttpAppLogAdd(D); err != nil {
		fmt.Println(err)
	}
	// 2、grpc方式发送日志
}

// grpc写常规日志
func gRpcLog() {
}

// grpc写应用日志
func gRpcAppLog() {
	D := &grpc_proto_applog.AppLogData{
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

	if err := modsynclog_v1.SdkGRpcAppLogAdd(D); err != nil {
		fmt.Println(err)
	}
}
