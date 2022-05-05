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

	logHandles        map[string]*filelog_v1.CFileLog // 日志处理实例map，键值为日期
	initHistoryDayNum int                             // 启动时同步前几天日志数据的天数，0,当天，-1，昨天……；默认为0
}
type ClientOptions struct {
	InitHistoryDayNum int // NewClient的更多参数项
	SendFlag          int // socket的传输码
}

// 对外函数：创建实例
func NewClient(serverIp string, serverPort int, logFolder string, opt ...ClientOptions) *Client {
	// 1、实例化客户端
	Me := &Client{
		logFolder:         logFolder,
		logHandles:        map[string]*filelog_v1.CFileLog{},
		initHistoryDayNum: 0,
	}
	// 同步历史数据天数设置
	if len(opt) > 0 && opt[0].InitHistoryDayNum < 0 {
		Me.initHistoryDayNum = opt[0].InitHistoryDayNum
	}

	// 2、创建客户端socket连接
	SocketClient := socket_v1.NewClient(serverIp, serverPort, Me.onMessage, Me.onConnectFail, Me.onConnect, Me.onDisConnect)
	// socket传输码设置
	if len(opt) > 0 && opt[0].SendFlag > 0 {
		SocketClient.Set("SendFlag", opt[0].SendFlag)
	}

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
	Me.ReqDate = utils_v1.Time().DateAdd(utils_v1.Time().Date(), Me.initHistoryDayNum)

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
		// 需要跨天处理
		// 理论上收到这条消息时，Me.ReqDate 日期肯定不会是当前服务器的日期，所以这个if应该始终是true
		if Me.ReqDate < utils_v1.Time().Date() {

			// 1、请求日期加1天
			oldDate := Me.ReqDate
			Me.ReqDate = utils_v1.Time().DateAdd(Me.ReqDate, 1)

			// 2、设置结束标记、将会在日志日期文件夹里创建一个结束标记文件
			Me.logHandles[oldDate].SetFinish()
			Me.logHandles[oldDate].Close()
			delete(Me.logHandles, oldDate)

			// 3、初始化新的一天日志句柄
			if _, ok := Me.logHandles[Me.ReqDate]; !ok {
				Me.logHandles[Me.ReqDate] = filelog_v1.New(Me.logFolder, Me.ReqDate)
				if id, err := Me.logHandles[Me.ReqDate].GetAutoId(); err != nil {
					log.Panic(err)
				} else {
					Me.ReqDateId = id
				}
			}

			// 4、告诉服务器新的接收数据日期和id
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

			// 1、日期不符
			if D.Date != Me.ReqDate {
				log.Error("日期不符")
				C.DisConnect()
			}

			// 2、当日期的日志handle不存在时，创建一个日志handle
			if _, ok := Me.logHandles[D.Date]; !ok {
				Me.logHandles[D.Date] = filelog_v1.New(Me.logFolder, D.Date)
			}

			// 3、Id校验，由于id是自增的，理论上这条日志数据的ID，应该就是日志句柄里当前日期的自增id
			if D.Id != Me.logHandles[D.Date].AutoId {
				log.Error("数据ID和客户端id不一致", D.Id, Me.logHandles[D.Date].AutoId)
				C.DisConnect()
				break
			}

			// 4、数据长度校验
			if len(D.Data) != int(D.DataLength) {
				log.Error("数据长度字段和计算长度不一致")
				C.DisConnect()
			}

			// 5、写入日志数据，（写数据后 Me.logHandles[D.Date]的AutoId,DataOffset都会变化）
			if _, err := Me.logHandles[D.Date].Add(D.Time, D.DataType, D.Data); err != nil {
				log.Error(err)
				C.DisConnect()
			}
		}
	}
}
