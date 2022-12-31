package cob

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	//DNepal.CreateFileServerByGin("C:\\lee\\project\\go\\auto", "/static/", "8888", true)
	DNepal.CreateFileServer("C:\\lee\\project\\go\\auto", "/static/", "8888", true)
}

func TestNepal_GetRandomString(t *testing.T) {
	fmt.Println(DNepal.GetRandomString(10))
}

func TestNepal_GetRandomNum(t *testing.T) {
	for i:=0;i<20;i++{
		time.Sleep(10*time.Millisecond)
		fmt.Printf("%v ", DNepal.GetRandomNum(1,1000))
	}
}
