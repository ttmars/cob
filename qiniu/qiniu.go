package qiniu

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Qiniu struct {
	AccessKey string		// 秘钥对
	SecretKey string
	Bucket string			// 空间(桶)
	Zone *storage.Zone		// 地域：
	Domain string
}

var DQiniu = &Qiniu{
	AccessKey: "xxx",
	SecretKey: "xxx",
	Bucket: "images-mkd",			// 空间与地域一一对应
	Zone: &storage.ZoneHuanan,		// https://developer.qiniu.com/kodo/1238/go
	Domain: "http://image.youthsweet.com",		// 空间绑定的域名
}

// SyncLocalToQiniu 同步本地目录到七牛，若minute=30，即只同步最近半小时的新文件
func (obj *Qiniu)SyncLocalToQiniu(localDir string, qiniuDir string, minute int)  {
	filepath.Walk(localDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			filename,err := filepath.Rel(localDir, path)
			if err != nil {
				log.Fatalln(err)
			}
			filename = filepath.Join(qiniuDir, filename)
			filename = strings.ReplaceAll(filename, "\\", "/")
			t1,_ := time.ParseDuration("-" + strconv.Itoa(minute) + "m")
			t2 := time.Now().Add(t1)
			if info.ModTime().After(t2) {
				obj.PutFile(path, filename)
				//fmt.Printf("文件名：%v 文件修改时间：%v 当前时间：%v 后退30M:%v\n", info.Name(), info.ModTime(), time.Now(), t2)
			}
		}
		return nil
	})
}

// GetOutsideChain 根据前缀匹配获取外链
func (obj *Qiniu)GetOutsideChain(prefix string) (urls []string) {
	mac := auth.New(obj.AccessKey, obj.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	cfg.Zone=obj.Zone
	bucketManager := storage.NewBucketManager(mac, &cfg)

	limit := 1000
	delimiter := ""
	//初始列举marker为空
	marker := ""
	for {
		entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(obj.Bucket, prefix, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}
		//print entries
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Key, "/") {
				urls = append(urls, obj.Domain + "/" + entry.Key)
			}
		}
		if hashNext {
			marker = nextMarker
		} else {
			//list end
			break
		}
	}
	return
}

// PutFile 上传文件
func (obj *Qiniu)PutFile(filepath string, filename string)  {
	putPolicy := storage.PutPolicy{
		Scope: obj.Bucket,				// 空间与地域一一对应，要同步更改
	}
	mac := qbox.NewMac(obj.AccessKey, obj.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = obj.Zone
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, filename, filepath, &putExtra)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("上传成功：", filepath)
}

// DeleteFile 删除bucket中的文件
func (obj *Qiniu)DeleteFile(filename string) error {
	mac := qbox.NewMac(obj.AccessKey, obj.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	cfg.Zone=obj.Zone
	bucketManager := storage.NewBucketManager(mac, &cfg)

	bucket := obj.Bucket
	key := filename
	err := bucketManager.Delete(bucket, key)
	return err
}
