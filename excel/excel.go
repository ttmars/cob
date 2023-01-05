package excel

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"path/filepath"
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
