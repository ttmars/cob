package http

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T)  {
	var proxys []string
	for i:=39000;i<=39500;i++{
		proxy := "http://fans007:fans888@45.76.169.156:" + strconv.Itoa(i)
		proxys = append(proxys, proxy)
	}

	tt := time.Now()
	fmt.Println("代理数量：", len(proxys))
	success,fail := DHttp.TestMulHttpProxy(proxys, 15, 100, false)
	fmt.Println("success:", len(success), "fail:", len(fail), "耗时:", time.Since(tt))
	fmt.Println(fail)
}
