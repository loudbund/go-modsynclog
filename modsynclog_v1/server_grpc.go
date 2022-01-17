package modsynclog_v1

import (
	"context"
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
	"github.com/loudbund/go-modsynclog/modsynclog_v1/grpc_proto_applog"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

//
type AppLogServer struct {
	server *Server
}

// grpc服务实例化 Addr:0.0.0.0:1235
func NewAppLog(Addr string, server *Server) {
	// 创建监听服务
	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// 服务注册
	grpc_proto_applog.RegisterAppLogServer(s, &AppLogServer{
		server: server,
	})

	// 启动监听
	fmt.Println("grpc开始监听:" + Addr)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// grpc远程函数
func (x *AppLogServer) AppLogWrite(ctx context.Context, in *grpc_proto_applog.AppLogRequest) (*grpc_proto_applog.AppLogResponse, error) { // 进行编码
	// 日志写入管道
	x.server.logChan <- &filelog_v1.UDataSend{
		DataType: 1101,
		Data:     in.Data,
	}

	// 3、写入日志
	return &grpc_proto_applog.AppLogResponse{ErrCode: 0, ErrMessage: "ok", Data: ""}, nil
}
