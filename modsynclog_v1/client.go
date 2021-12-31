package modsynclog_v1

import (
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
	"github.com/loudbund/go-json/json_v1"
	"github.com/loudbund/go-socket/socket_v1"
	"github.com/loudbund/go-utils/utils_v1"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	logFolder string // 日志文件目录

	ReqDate   string // 请求日志的日期
	ReqDateId int64  // 请求日志的位置

	logHandles map[string]*filelog_v1.CFileLog // 日志处理实例map，键值为日期
}

// 对外函数：创建实例
func NewClient(serverIp string, serverPort int, logFolder string) *Client {
	// 1、实例化客户端
	Me := &Client{
		logFolder:  logFolder,
		logHandles: map[string]*filelog_v1.CFileLog{},
	}

	// 2、创建客户端socket连接
	SocketClient := socket_v1.NewClient(serverIp, serverPort, Me.onMessage, Me.onConnectFail, Me.onConnect, Me.onDisConnect)
	go SocketClient.Connect()

	return Me
}

// 1.1、收到了消息回调函数，这里处理消息
func (Me *Client) onMessage(Msg socket_v1.UDataSocket, C *socket_v1.Client) {
	Me.onMsg(Msg, C)
}

// 1.2、连接失败回调函数
func (Me *Client) onConnectFail(C *socket_v1.Client) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "连接失败！5秒后重连")
	go C.ReConnect(5) // 延时5秒后重连
}

// 1.3、连接成功回调函数
func (Me *Client) onConnect(C *socket_v1.Client) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "连接成功！")
	Me.initClientLogDate()
	Content, _ := json_v1.JsonEncode(map[string]interface{}{
		"date": Me.ReqDate,
		"id":   Me.ReqDateId,
	})
	_ = C.SendMsg(socket_v1.UDataSocket{
		Zlib:    0,
		CType:   301,
		Content: []byte(Content),
	})
}

// 1.4、掉线回调函数
func (Me *Client) onDisConnect(C *socket_v1.Client) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "掉线了,5秒后重连")
	for k, v := range Me.logHandles {
		v.Close()
		delete(Me.logHandles, k)
	}
	go C.ReConnect(5) // 延时5秒后重连
}

// 1.5、初始胡客户端日志位置
func (Me *Client) initClientLogDate() {
	// 初始胡日期为7天前,位置为-1
	Me.ReqDate = utils_v1.Time().DateAdd(utils_v1.Time().Date(), -7)

	if _, ok := Me.logHandles[Me.ReqDate]; !ok {
		Me.logHandles[Me.ReqDate] = filelog_v1.New(Me.logFolder, Me.ReqDate)
		if id, err := Me.logHandles[Me.ReqDate].GetAutoId(); err != nil {
			log.Panic(err)
		} else {
			Me.ReqDateId = id
		}
	}
}

// 2、消息处理
func (Me *Client) onMsg(Msg socket_v1.UDataSocket, C *socket_v1.Client) {

	if Msg.CType == 304 { // 收到发送结束消息

		// 当前处理的日期是昨天前
		if Me.ReqDate < utils_v1.Time().Date() {

			// 请求日期加1天
			oldDate := Me.ReqDate
			Me.ReqDate = utils_v1.Time().DateAdd(Me.ReqDate, 1)

			// 设置结束标记
			Me.logHandles[oldDate].SetFinish()
			Me.logHandles[oldDate].Close()
			delete(Me.logHandles, oldDate)

			// 初始化新的一天日志句柄
			if _, ok := Me.logHandles[Me.ReqDate]; !ok {
				Me.logHandles[Me.ReqDate] = filelog_v1.New(Me.logFolder, Me.ReqDate)
				if id, err := Me.logHandles[Me.ReqDate].GetAutoId(); err != nil {
					log.Panic(err)
				} else {
					Me.ReqDateId = id
				}
			}

			// 告诉服务器新的接收数据日期和id
			Content, _ := json_v1.JsonEncode(map[string]interface{}{
				"date": Me.ReqDate,
				"id":   Me.ReqDateId,
			})
			_ = C.SendMsg(socket_v1.UDataSocket{
				Zlib:    0,
				CType:   301,
				Content: []byte(Content),
			})
		}
	} else if Msg.CType == 302 { // 收到日志消息
		// 日志消息解密
		Rows := utilsDecodeUData(Msg.Content)
		for _, D := range Rows {
			// fmt.Println(Msg.CType, D.Id, D.Date, D.Time, D.DataType, D.DataLength, D.DataOffset, string(D.Data))

			// 日期不符
			if D.Date != Me.ReqDate {
				log.Error("日期不符")
				C.DisConnect()
			}

			// 初始化
			if _, ok := Me.logHandles[D.Date]; !ok {
				Me.logHandles[D.Date] = filelog_v1.New(Me.logFolder, D.Date)
			}

			// Id校验
			if D.Id != Me.logHandles[D.Date].AutoId {
				log.Error("数据ID和客户端id不一致", D.Id, Me.logHandles[D.Date].AutoId)
				C.DisConnect()
				break
			}

			// 数据长度校验
			if len(D.Data) != int(D.DataLength) {
				log.Error("数据长度字段和计算长度不一致")
				C.DisConnect()
			}

			// 写入数据，（写数据后 Me.logHandles[D.Date]的AutoId,DataOffset都会变化）
			if _, err := Me.logHandles[D.Date].Add(D.Time, D.DataType, D.Data); err != nil {
				log.Error(err)
				C.DisConnect()
			}
		}
	}
}
