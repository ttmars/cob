package cob

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Nepal struct {

}
var DNepal = Nepal{}

// GetRandomString 获取n位随机字符串，n为偶数
func (obj *Nepal)GetRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// GetRandomNum 获取N~M之间的随机整数，N<=M
func (obj *Nepal)GetRandomNum(N,M int) int  {
	if M<N {
		log.Fatalln("N必须小于等于M!")
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(M-N+1)
	return n+N
}
