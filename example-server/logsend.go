package main

import (
	"fmt"
	"github.com/loudbund/go-modsynclog/modsynclog_v1"
	"github.com/loudbund/go-utils/utils_v1"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	modsynclog_v1.SetDomain("http://127.0.0.1:1234")
	// for i := 0; i < 11; i++ {
	// 	if err := modsynclog_v1.SendLog(123, []string{utils_v1.Time().DateTime()}); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
	Ds := make([]string, 0)
	for i := 0; i < 11; i++ {
		Ds = append(Ds, utils_v1.Time().DateTime())
	}
	if err := modsynclog_v1.SendLog(123, Ds); err != nil {
		fmt.Println(err)
	}
}
