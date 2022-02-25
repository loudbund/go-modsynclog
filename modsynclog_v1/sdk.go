package modsynclog_v1

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/loudbund/go-json/json_v1"
	"github.com/loudbund/go-modsynclog/modsynclog_v1/grpc_proto_log"
	"github.com/loudbund/go-request/request_v1"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
)

type AppLog struct {
	dataType      int32
	httpDomain    string                   // http每次连接地址
	gRpcLogClient grpc_proto_log.LogClient // grpc长连接句柄
}
type DbLog struct {
	dataType      int32
	httpDomain    string                   // http每次连接地址
	gRpcLogClient grpc_proto_log.LogClient // grpc长连接句柄
}

func NewSdkAppLog(httpDomain, gRpcAddress string) *AppLog {
	handle := &AppLog{
		dataType: 1101,
	}
	if httpDomain != "" {
		handle.httpDomain = httpDomain
	}
	if gRpcAddress != "" {
		// 创建grpc连接
		conn, err := grpc.Dial(gRpcAddress, grpc.WithInsecure())
		if err != nil {
			log.Error("日志grpc连接失败: %v", err)
		}
		// defer func() { _ = conn.Close() }()
		handle.gRpcLogClient = grpc_proto_log.NewLogClient(conn)
	}
	return handle
}

// http程序日志写入
func (Me *AppLog) SdkAppLogAddHttp(inData []*grpc_proto_log.AppLogData) error {
	if Me.httpDomain == "" {
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
	if code, body, err := request_v1.PostForm(Me.httpDomain+"/log/write", map[string]string{
		"type": strconv.Itoa(int(Me.dataType)),
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

// grpc程序日志写入
func (Me *AppLog) SdkAppLogAddGRpc(inData *grpc_proto_log.AppLogData) error {
	// 数据内容体现用proto加密
	data, err := proto.Marshal(inData)
	if err != nil {
		log.Error("marshaling error: ", err)
	}

	// 调用grpc写日志
	r, err := Me.gRpcLogClient.AppLogWrite(context.Background(), &grpc_proto_log.AppLogRequest{
		DataType: Me.dataType,
		Data:     data,
	})
	if err != nil {
		log.Error("could not greet: %v", err)
		return err
	}
	if r.ErrCode != 0 {
		return errors.New(r.ErrMessage)
	}
	return nil
}

// 数据表日志 ///////////////////////////////////////////

func NewSdkDbLog(httpDomain, gRpcAddress string) *DbLog {
	handle := &DbLog{
		dataType: 1001,
	}
	if httpDomain != "" {
		handle.httpDomain = httpDomain
	}
	if gRpcAddress != "" {
		// 创建grpc连接
		conn, err := grpc.Dial(gRpcAddress, grpc.WithInsecure())
		if err != nil {
			log.Error("日志grpc连接失败: %v", err)
		}
		// defer func() { _ = conn.Close() }()
		handle.gRpcLogClient = grpc_proto_log.NewLogClient(conn)
	}
	return handle
}

// http数据表日志写入
func (Me *DbLog) SdkDbLogAddHttp(inData []*grpc_proto_log.DbLogData) error {
	if Me.httpDomain == "" {
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
	if code, body, err := request_v1.PostForm(Me.httpDomain+"/log/write", map[string]string{
		"type": strconv.Itoa(int(Me.dataType)),
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

// grpc数据表日志写入
func (Me *DbLog) SdkDbLogAddGRpc(inData *grpc_proto_log.DbLogData) error {
	// 数据内容体现用proto加密
	data, err := proto.Marshal(inData)
	if err != nil {
		log.Error("marshaling error: ", err)
	}

	// 调用grpc写日志
	r, err := Me.gRpcLogClient.DbLogWrite(context.Background(), &grpc_proto_log.DbLogRequest{
		DataType: Me.dataType,
		Data:     data,
	})
	if err != nil {
		log.Error("could not greet: %v", err)
		return err
	}
	if r.ErrCode != 0 {
		return errors.New(r.ErrMessage)
	}
	return nil
}
