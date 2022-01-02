package modsynclog_v1

import (
	"errors"
	"github.com/loudbund/go-json/json_v1"
	"github.com/loudbund/go-request/request_v1"
	"strconv"
	"strings"
)

var Domain = ""

// 设置域名
func SetDomain(domain string) {
	Domain = strings.TrimRight(domain, "/") + "/"
}

// 发送日志
func SendLog(DataType int, Data []string) error {
	if Domain == "" {
		return errors.New("未设置domain，请先执行SetDomain设置域名")
	}

	// 执行写日志接口请求
	D, _ := json_v1.JsonEncode(Data)
	if code, body, err := request_v1.PostForm(Domain+"log/write", map[string]string{
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
