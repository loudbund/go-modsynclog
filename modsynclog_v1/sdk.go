package modsynclog_v1

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/loudbund/go-json/json_v1"
	"github.com/loudbund/go-modsynclog/modsynclog_v1/grpc_proto_applog"
	"github.com/loudbund/go-request/request_v1"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"strings"
)

var httpDomain = ""                                 // http每次连接地址
var gRpcAppLogClient grpc_proto_applog.AppLogClient // grpc长连接句柄

// http地址设置
func SdkInitHttpSet(domain string) {
	httpDomain = strings.TrimRight(domain, "/") + "/"
}

// http数据写入
func SdkHttpDataAdd(DataType int, Data []string) error {
	if httpDomain == "" {
		return errors.New("未设置domain，请先执行SetDomain设置域名")
	}

	// 执行写日志接口请求
	D, _ := json_v1.JsonEncode(Data)
	if code, body, err := request_v1.PostForm(httpDomain+"log/write", map[string]string{
		"type": strconv.Itoa(DataType),
		"data": D,
	}); err != nil {
		return errors.New("日志发送失败:" + err.Error())
	} else if code != 200 {
		return errors.New("日志发送失败:状态码不为200，为" + strconv.Itoa(code))
	} else if body != `{"errcode": 0,"data":"ok"}` {
		return errors.New("日志发送失败:返回body:" + body)
	}
	return nil
}

// http程序日志写入
func SdkHttpAppLogAdd(inData []*grpc_proto_applog.AppLogData) error {
	if httpDomain == "" {
		return errors.New("未设置domain，请先执行SetDomain设置域名")
	}

	// 数据内容体现用proto加密
	sData := make([]string, 0)
	for _, v := range inData {
		data, err := proto.Marshal(v)
		if err != nil {
			log.Error("marshaling error: ", err)
		}
		sData = append(sData, string(data))
	}
	data, _ := json_v1.JsonEncode(sData)

	// 执行写日志接口请求
	if code, body, err := request_v1.PostForm(httpDomain+"log/write", map[string]string{
		"type": "1101",
		"data": data,
	}); err != nil {
		return errors.New("日志发送失败:" + err.Error())
	} else if code != 200 {
		return errors.New("日志发送失败:状态码不为200，为" + strconv.Itoa(code))
	} else if body != `{"errcode": 0,"data":"ok"}` {
		return errors.New("日志发送失败:返回body:" + body)
	}
	return nil
}

// grpc地址设置 127.0.0.1:1235
func SdkInitGRpcAppLog(Address string) {
	// 创建grpc连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}
	// defer func() { _ = conn.Close() }()

	// 执行调用
	gRpcAppLogClient = grpc_proto_applog.NewAppLogClient(conn)
}

// grpc程序日志写入
func SdkGRpcAppLogAdd(inData *grpc_proto_applog.AppLogData) error {
	// 数据内容体现用proto加密
	data, err := proto.Marshal(inData)
	if err != nil {
		log.Error("marshaling error: ", err)
	}

	// 调用grpc写日志
	r, err := gRpcAppLogClient.AppLogWrite(context.Background(), &grpc_proto_applog.AppLogRequest{
		DataType: 1101,
		Data:     data,
	})
	if err != nil {
		log.Fatal("could not greet: %v", err)
		return err
	}
	if r.ErrCode != 0 {
		return errors.New(r.ErrMessage)
	}
	return nil
}
