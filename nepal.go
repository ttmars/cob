package cob

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"time"
)

type Nepal struct {

}
var DNepal = Nepal{}

// GetRandomString 获取n位随机字符串，n为偶数
func (obj *Nepal)GetRandomString(n int) (result string) {
	rand.Seed(time.Now().UnixNano())
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// GetRandomNum 获取N~M之间的随机整数，N<=M
func (obj *Nepal)GetRandomNum(N,M int) (result int)  {
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

// GetRootPath 获取项目根路径，请将函数复制到自己的项目中使用，不能直接调用！
func (obj *Nepal)GetRootPath() (root string) {
	_,filename,_,_ := runtime.Caller(0)
	return path.Dir(filename)
}

// HttpGet 请求一个URL，返回状态码和响应体
func (obj *Nepal)HttpGet(url string, client *http.Client) (result string, statusCode int) {
	if client == nil {
		client = http.DefaultClient
		client.Timeout = time.Second*3
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return "",0
	}

	resp,err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "",0
	}
	defer resp.Body.Close()

	b,err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "",resp.StatusCode
	}
	return string(b),resp.StatusCode
}

// ReadFileToSlice 将文件中的所有非空行读入切片
func (obj *Nepal)ReadFileToSlice(filepath string) (lines []string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if string(line) != "" {
			lines = append(lines, string(line))
		}
	}
	return
}

// CreateHttpProxyClient 创建一个http代理客户端，代理类型由uri确定，支持http/https/socks5，默认为http
func (obj *Nepal)CreateHttpProxyClient(uri string, user string, pass string, timeout int) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		u,err := url.Parse(uri)
		if user != "" && pass != "" {
			u.User = url.UserPassword(user, pass)
		}
		return u,err
	}

	return &http.Client{
		Transport:&http.Transport{Proxy: proxy},
		Timeout: time.Duration(timeout) * time.Second,
	}
}

// InitLog 初始化日志器，设置日志前缀、日志路径，默认添加行号显示
func (obj *Nepal)InitLog(logPrefix string, logFilePath string)  {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if logPrefix != "" {
		log.SetPrefix(logPrefix)
	}

	if logFilePath != "" {
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetOutput(logFile)
	}
}
