package cob

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"
)

type Nepal struct {

}
var DNepal = &Nepal{}

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

// GetRootPath 获取项目根路径，请将函数复制到自己的项目中使用，不能直接调用！
func (obj *Nepal)GetRootPath() (root string) {
	_,filename,_,_ := runtime.Caller(0)
	return path.Dir(filename)
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

// 将任意结构体切片转换为二维切片[][]interface{}
func (obj *Nepal)StructToSlice(slice interface{}, addFieldName bool) [][]interface{} {
	// 将传入的切片转换为 reflect.Value 类型
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		log.Fatalln("传入的不是切片")
	}

	// 创建结果切片
	result := make([][]interface{}, val.Len())

	// 遍历切片中的所有元素
	for i := 0; i < val.Len(); i++ {
		// 获取每个元素的类型
		elemType := val.Index(i).Type()
		// 创建一个新的切片，用来存储每个元素中的所有字段
		innerSlice := make([]interface{}, elemType.NumField())
		// 遍历元素中的所有字段
		for j := 0; j < elemType.NumField(); j++ {
			// 获取字段值，并将其转换为 interface{} 类型
			fieldVal := val.Index(i).Field(j).Interface()
			innerSlice[j] = fieldVal
		}
		result[i] = innerSlice
	}

	if addFieldName {
		elemType := val.Index(0).Type()
		fieldNameSlice := make([]interface{}, elemType.NumField())
		for j := 0; j < elemType.NumField(); j++ {
			fieldName := elemType.Field(j).Name
			fieldNameSlice[j] = fieldName
		}
		var tmp [][]interface{}
		tmp = append(tmp, fieldNameSlice)
		result = append(tmp, result...)
	}

	return result
}
