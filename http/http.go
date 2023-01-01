package http

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Http struct {

}

var DHttp = &Http{}

// CreateHttpProxyClient 创建一个http代理客户端，代理类型由uri确定，支持http/https/socks5，默认为http
func (obj *Http)CreateHttpProxyClient(uri string, user string, pass string, timeout int) *http.Client {
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

// CreateFileServerByGin 创建一个HTTP文件服务器，通过gin框架
func (obj *Http)CreateFileServerByGin(localPath string, relativePath string, port string, OnlyListenLocalAddr bool)  {
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
func (obj *Http)CreateFileServer(localPath string, relativePath string, port string, OnlyListenLocalAddr bool)  {
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

// Get 请求一个URL，返回响应体
func (obj *Http)Get(url string) (result string) {
	http.DefaultClient.Timeout = time.Second * 3
	resp,err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b,err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(resp.StatusCode, err)
	}
	return string(b)
}
