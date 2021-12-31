package modsynclog_v1

import (
	"container/list"
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
	"github.com/loudbund/go-json/json_v1"
	"github.com/loudbund/go-socket/socket_v1"
	"github.com/loudbund/go-utils/utils_v1"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 结构体1： 单个用户数据
type User struct {
	ClientId  string // 客户端id
	ClientIp  string // 客户端ip
	LoginTime string // App连接时间

	ReqDate   string // 请求日志的日期
	ReqDateId int64  // 请求日志的位置

	logReadHandles map[string]*filelog_v1.CFileLog // 日志处理实例map，键值为日期
}

// 结构体2： 服务端结构体
type Server struct {
	SocketServer *socket_v1.Server
	ListUser     *list.List       // 客户端链表
	Users        map[string]*User // 客户端clientid和User的map关系
	lockListUser sync.RWMutex     // 客户端链表同步锁

	date       string                          // 当前存入日志的日期
	logFolder  string                          // 日志文件目录
	logChan    chan *filelog_v1.UDataSend      // 并发转线性处理通道
	logHandles map[string]*filelog_v1.CFileLog // 日志处理实例map，键值为日期
}

// 对外函数：创建实例
func NewServer(Ip string, PortHttp, PortSocket int, logFolder string) *Server {
	Me := &Server{
		Users:      map[string]*User{},
		ListUser:   list.New(),
		date:       utils_v1.Time().Date(),
		logFolder:  logFolder,
		logChan:    make(chan *filelog_v1.UDataSend),
		logHandles: map[string]*filelog_v1.CFileLog{},
	}

	// 1、关闭前几天的日志
	Me.closePreDateLog()

	// 1、写日志协程
	go Me.messageWrite()

	// 2、 http服务,提供日志写入api
	http.HandleFunc("/log/write", Me.write)
	fmt.Println("http开始监听:" + Ip + ":" + strconv.Itoa(PortHttp))
	go func() {
		if err := http.ListenAndServe(Ip+":"+strconv.Itoa(PortHttp), nil); err != nil {
			log.Error(err)
		}
	}()

	// 3、socket服务器
	Me.SocketServer = socket_v1.NewServer(Ip, PortSocket, func(Event socket_v1.HookEvent) {
		Me.onHookEvent(Event)
	})

	return Me
}

// 关闭前几天的日志文件
func (Me *Server) closePreDateLog() {
	Today := utils_v1.Time().Date()
	iDate := utils_v1.Time().DateAdd(utils_v1.Time().Date(), -7)
	I := 0
	for {
		// 最多执行20次，避免死循环
		I++
		if I > 20 {
			break
		}

		// 只处理昨天前
		if iDate == Today {
			break
		}

		// 处理这天的
		handle := filelog_v1.New(Me.logFolder, iDate)
		handle.SetFinish()
		handle.Close()

		// 日期加1天
		iDate = utils_v1.Time().DateAdd(iDate, 1)
	}
}

// 处理http写日志请求
func (Me *Server) write(writer http.ResponseWriter, request *http.Request) {
	// 参数接收,【type必须小于20000】
	KeyType := int16(0)
	KeyData := ""
	if true {
		vType := request.PostFormValue("type")
		vData := request.PostFormValue("data")
		if vType == "" || vData == "" {
			query := request.URL.Query()
			if vType == "" {
				if _, ok := query["type"]; ok {
					vType = query["type"][0]
				}
			}
			if vData == "" {
				if _, ok := query["data"]; ok {
					vData = query["data"][0]
				}
			}
		}
		if vType != "" {
			if d, err := strconv.Atoi(vType); err == nil {
				if d < 20000 {
					KeyType = int16(d)
				}
			}
		}
		KeyData = vData
	}
	// 参数判断,【type和data必须有值】
	if KeyType == 0 || KeyData == "" {
		_, _ = writer.Write([]byte(`{"errcode": 101,"errmsg":"参数错误","err":""}`))
	} else {
		// 日志写入管道
		Me.logChan <- &filelog_v1.UDataSend{
			DataType: KeyType,
			Data:     []byte(KeyData),
		}
		_, _ = writer.Write([]byte(`{"errcode": 0,"data":"ok"}`))
	}
}

// 管道接收日志并写文件
func (Me *Server) messageWrite() {
	T := time.NewTicker(time.Second)
	for {
		select {
		case <-T.C:
			// 判断是否跨天了
			Time := time.Now().Unix()
			if true {
				Date := utils_v1.Time().Date(time.Unix(Time, 0))
				if Date != Me.date {
					// 关闭
					if _, ok := Me.logHandles[Me.date]; ok {
						Me.logHandles[Me.date].Close()
						delete(Me.logHandles, Me.date)
					}
					Me.date = Date
				}
			}

		case D, ok := <-Me.logChan:
			if !ok {
				return
			}
			// 判断是否跨天了
			Time := time.Now().Unix()
			if true {
				Date := utils_v1.Time().Date(time.Unix(Time, 0))
				if Date != Me.date {
					// 关闭
					if _, ok := Me.logHandles[Me.date]; ok {
						Me.logHandles[Me.date].Close()
						delete(Me.logHandles, Me.date)
					}
					Me.date = Date
				}
			}

			// 准备写入日志
			if _, ok := Me.logHandles[Me.date]; !ok {
				Me.logHandles[Me.date] = filelog_v1.New(Me.logFolder, Me.date)
			}
			if _, err := Me.logHandles[Me.date].Add(int32(Time), D.DataType, D.Data); err != nil {
				log.Error(err)
			}
		}
	}
}

// 1、处理数据,多线程转单线程处理
func (Me *Server) onHookEvent(Event socket_v1.HookEvent) {
	switch Event.EventType {
	case "message": // 1、消息事件
		fmt.Println("message:", utils_v1.Time().DateTime(), Event.Message.CType, string(Event.Message.Content))
		// 客户端请求日期和开始位置日志
		if Event.Message.CType == 301 {
			if jData, err := json_v1.JsonDecode(string(Event.Message.Content)); err != nil {
				log.Error(err)
			} else {
				Date, err1 := json_v1.GetJsonString(jData, "date")
				id, err2 := json_v1.GetJsonInt64Force(jData, "id")
				if err1 != nil {
					log.Error(err1)
				} else if err2 != nil {
					log.Error(err2)
				} else {
					Me.lockListUser.Lock()
					if _, ok := Me.Users[Event.User.ClientId]; ok {
						Me.Users[Event.User.ClientId].ReqDate = Date
						Me.Users[Event.User.ClientId].ReqDateId = id
					}
					Me.lockListUser.Unlock()
				}
			}
		}

	case "offline": // 2、下线事件
		Me.removeUser(Event.User.ClientId)

	case "online": // 3、上线消息
		U := Me.addUser(Event.User.ClientId, Event.User.Addr)
		go Me.sendLog(U)
	}
}

// 发送日志给客户端
func (Me *Server) sendLog(U *User) {
	fmt.Println("start send log:", U.ClientId)
	for {
		if _, ok := Me.Users[U.ClientId]; !ok {
			return
		}

		if U.ReqDate != "" {
			if _, ok := U.logReadHandles[U.ReqDate]; !ok {
				U.logReadHandles[U.ReqDate] = filelog_v1.New(Me.logFolder, U.ReqDate)
			} else {
				if D, err := U.logReadHandles[U.ReqDate].GetOne(U.ReqDateId); err != nil {
					log.WithFields(log.Fields{"n": "取数据失败"}).Error(err)
					time.Sleep(time.Second)
				} else if D == nil {
					// 如果日期内的日志已经发送完成，则发送标记
					if finished := U.logReadHandles[U.ReqDate].GetFinish(); finished {
						if err := Me.SocketServer.SendMsg(&U.ClientId, socket_v1.UDataSocket{
							Zlib:    1,
							CType:   304,
							Content: []byte(U.ReqDate),
						}); err != nil {
							log.WithFields(log.Fields{"n": "消息发送失败"}).Error(err)
							return
						}
					}
					time.Sleep(time.Second)
				} else {
					sendData := make([]*filelog_v1.UDataSend, 0)
					sendData = append(sendData, D)

					if err := Me.SocketServer.SendMsg(&U.ClientId, socket_v1.UDataSocket{
						Zlib:    1,
						CType:   302,
						Content: utilsEncodeUData(sendData),
					}); err != nil {
						log.WithFields(log.Fields{"n": "消息发送失败"}).Error(err)
						return
					} else {
						if _, ok := Me.Users[U.ClientId]; ok {
							Me.Users[U.ClientId].ReqDateId++
						}
					}
				}
			}
		} else {
			time.Sleep(time.Second)
		}
	}
}

// 添加用户
func (Me *Server) addUser(ClientId, Addr string) *User {
	IpPort := strings.Split(Addr, ":")
	U := &User{
		ClientId:       ClientId,
		ClientIp:       IpPort[0],
		LoginTime:      utils_v1.Time().DateTime(),
		logReadHandles: map[string]*filelog_v1.CFileLog{},
	}
	Me.lockListUser.Lock()
	Me.ListUser.PushBack(U)
	Me.Users[U.ClientId] = U
	Me.lockListUser.Unlock()

	return U
}

// 移除用户
func (Me *Server) removeUser(ClientId string) {
	Me.lockListUser.Lock()
	for e := Me.ListUser.Front(); e != nil; e = e.Next() {
		if e.Value.(*User).ClientId == ClientId {
			Me.ListUser.Remove(e)
			delete(Me.Users, ClientId)
		}
	}
	Me.lockListUser.Unlock()
}
