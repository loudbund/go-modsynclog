package main

import (
	"fmt"
	"github.com/loudbund/go-filelog/filelog_v1"
)

// 6、主函数 -------------------------------------------------------------------------
func main() {
	handle := filelog_v1.New("/tmp/test-synlog-client", "2021-12-29")
	D, _ := handle.GetOne(1)
	fmt.Println(D)
}
