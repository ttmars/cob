package cob

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"path"
	"runtime"
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

// CreateFileServerByGin 创建一个HTTP文件服务器，通过gin框架
func (obj *Nepal)CreateFileServerByGin(localPath string, relativePath string, port string, OnlyListenLocalAddr bool)  {
	router := gin.Default()
	router.StaticFS(relativePath, http.Dir(localPath))
	if OnlyListenLocalAddr {
		router.Run("127.0.0.1:"+port)
	}else{
		router.Run(":"+port)
	}
}

// CreateFileServer 创建一个HTTP文件服务器，注意relativePath路径后要加斜杠
// example:DNepal.CreateFileServer("C:\\lee\\project\\go\\auto", "/static/", "8888", true)
func (obj *Nepal)CreateFileServer(localPath string, relativePath string, port string, OnlyListenLocalAddr bool)  {
	fs := http.FileServer(http.Dir(localPath))
	http.Handle(relativePath, http.StripPrefix(relativePath, fs))
	if OnlyListenLocalAddr {
		log.Printf("Listening and serving HTTP on 127.0.0.1:%s\n", port)
		http.ListenAndServe("127.0.0.1:"+port, nil)
	}else{
		log.Printf("Listening and serving HTTP on 0.0.0.0:%s\n", port)
		http.ListenAndServe(":"+port, nil)
	}
}

// GetRootPath 获取项目根路径
func (obj *Nepal)GetRootPath() (root string) {
	_,filename,_,_ := runtime.Caller(0)
	return path.Dir(filename)
}
