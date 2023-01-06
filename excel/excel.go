package excel

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"reflect"
)

type Excel struct {

}

var DExcel = &Excel{}

// SliceToExcelBuff 将切片转换为Excel格式，并返回bytes.Buffer，可以保存为文件或写入http响应体
func (obj *Excel)SliceToExcelBuff(s [][]interface{}) (buf *bytes.Buffer) {
	// 创建一个Excel文件实例
	f := excelize.NewFile()
	defer f.Close()
	// 添加表格
	index := f.NewSheet("Sheet1")
	// 插入值
	for i,v := range s {
		for j,vv := range v {
			f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 65+j, i+1), vv)
		}
	}
	// 将Sheet1设置为当前表格
	f.SetActiveSheet(index)
	buf,_ = f.WriteToBuffer()
	return
}

// SliceToExcelFile 将切片写入Excel文件
func (obj *Excel)SliceToExcelFile(s [][]interface{}, path string, filename string) (err error) {
	// 创建一个Excel文件实例
	f := excelize.NewFile()
	defer f.Close()
	// 添加表格
	index := f.NewSheet("Sheet1")
	// 插入值
	for i,v := range s {
		for j,vv := range v {
			f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 65+j, i+1), vv)
		}
	}
	// 将Sheet1设置为当前表格
	f.SetActiveSheet(index)
	fp := filepath.Join(path, filename)
	fmt.Println(fp)
	if err = f.SaveAs(fp); err != nil {
		return err
	}
	return nil
}

// StructSliceToGinResp 将结构体切片写入Excel文件，写入http响应体，实现Excel文件下载
func (obj *Excel)StructSliceToGinResp(slice interface{}, filename string, context *gin.Context)  {
	obj.SliceToGinResp(obj.StructToSlice(slice, true), filename, context)
}

// SliceToGinResp 将切片写入Excel文件，写入http响应体，实现Excel文件下载
func (obj *Excel)SliceToGinResp(s [][]interface{}, filename string, context *gin.Context)  {
	// 创建一个Excel文件实例
	f := excelize.NewFile()
	defer f.Close()
	// 添加表格
	index := f.NewSheet("Sheet1")
	// 插入值
	for i,v := range s {
		for j,vv := range v {
			f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 65+j, i+1), vv)
		}
	}
	// 将Sheet1设置为当前表格
	f.SetActiveSheet(index)

	context.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	context.Writer.Header().Set("Content-Type", "application/vnd.ms-excel")
	f.WriteTo(context.Writer)
}

// StructToSlice 将任意结构体切片转换为二维切片[][]interface{}，addFieldName是否添加字段名
func (obj *Excel)StructToSlice(slice interface{}, addFieldName bool) [][]interface{} {
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
			var fieldVal interface{}
			if val.Index(i).Field(j).Kind() == reflect.Ptr {			// 判断字段是否为指针类型
				fieldVal = val.Index(i).Field(j).Elem().Interface()
			}else{
				fieldVal = val.Index(i).Field(j).Interface()
			}
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
