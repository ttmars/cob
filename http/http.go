package http

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Http struct {

}

var DHttp = &Http{}

// TestMulHttpProxy 批量测试多个HTTP代理
func (obj *Http)TestMulHttpProxy(proxys []string, timeout int, maxG int, printLog bool) (success []string, fail []string){
	var wg sync.WaitGroup
	var mu sync.Mutex
	ch := make(chan bool, maxG)
	for _,proxy := range proxys {
		wg.Add(1)
		ch <- true
		go func(proxy string) {
			t := time.Now()
			_,err := obj.TestOneHttpProxy(proxy, timeout)
			if err != nil {
				if printLog {
					log.Println(proxy, "fail", err)
				}
				mu.Lock()
				fail = append(fail, proxy)
				mu.Unlock()
			}else{
				if printLog {
					log.Println(proxy, "success", time.Since(t))
				}
				mu.Lock()
				success = append(success, proxy)
				mu.Unlock()
			}
			wg.Done()
			<-ch
		}(proxy)
	}
	wg.Wait()
	return
}

// TestOneHttpProxy 测试单个http代理是否可用
func (obj *Http)TestOneHttpProxy(proxy string, timeout int) (string, error) {
	client := obj.CreateHttpProxyClient(proxy, timeout)
	resp,err := client.Get("https://ip.sb/")
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	b,err := io.ReadAll(resp.Body)
	if err != nil{
		return "",err
	}
	return string(b), nil
}

// CreateHttpProxyClient 创建一个http代理客户端，代理类型由uri确定，支持http/https/socks5，默认为http
// http://fans007:fans888@45.76.169.156:39000
func (obj *Http)CreateHttpProxyClient(proxy string, timeout int) *http.Client {
	return &http.Client{
		Transport:&http.Transport{Proxy: func(_ *http.Request) (*url.URL, error){
			return url.Parse(proxy)
		}},
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
