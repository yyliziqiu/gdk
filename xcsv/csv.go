package xcsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/yyliziqiu/gdk/xfile"
	"github.com/yyliziqiu/gdk/xutil"
)

func Save(filename string, models any) error {
	v := reflect.ValueOf(models)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("models type must be slice or array")
	}

	size := v.Len()
	if size == 0 {
		return nil
	}

	rows := make([][]string, 0, size+1)
	rows = append(rows, titles(v.Index(0).Interface()))
	for i := 0; i < size; i++ {
		rows = append(rows, xutil.ReflectValuesToString(v.Index(i).Interface()))
	}

	return SaveRows(filename, rows)
}

func SaveRows(filename string, rows [][]string) error {
	// 创建存储目录
	if err := xfile.MakeDir(filepath.Dir(filename)); err != nil {
		return fmt.Errorf("mkdir failed [%v]", err)
	}

	// 优化文件名
	if !strings.HasSuffix(filename, ".csv") {
		filename = filename + ".csv"
	}

	// 创建 CSV 文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create CSV file failed [%v]", err)
	}
	defer file.Close()

	// 写入 CSV 文件
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err = writer.WriteAll(rows); err != nil {
		return fmt.Errorf("write date to CSV failed [%v]", err)
	}

	return nil
}

func titles(s any) []string {
	mt := reflect.TypeOf(s)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		title := mt.Field(i).Tag.Get("csv")
		if title == "" {
			title = mt.Field(i).Name
		}
		fields = append(fields, title)
	}
	return fields
}
