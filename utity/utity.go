package utity

import (
	"WorkloadQuery/conf"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// 工具包

// IsWithinWorkingTime 判断是否在工作时间;决定是否执行作业
func IsWithinWorkingTime() bool {
	hour := time.Now().Hour()
	return conf.Configs.CustomTaskTime.Run == 0 && hour >= conf.Configs.CustomTaskTime.StartTime && hour < conf.Configs.CustomTaskTime.EndTime
}

// RemoveAssignDir 清楚目录下指定后缀文件名
func RemoveAssignDir(dirPath, suffix string) (err error) {
	file := make([]string, 0, 500)
	//suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), suffix) {
			file = append(file, path)
		}
		return nil
	})
	done := make(chan int, len(file))
	group := new(errgroup.Group)
	for _, v := range file {
		filename := v
		group.Go(func() error {
			err = os.Remove(filename)
			done <- 1
			return err
		})
	}
	for i := 0; i < len(file); i++ {
		<-done
	}
	if err = group.Wait(); err != nil {
		return err
	}
	return nil
}

// RemoveTabsFromStruct 清理结构体中的所有字符串字段中的制表符
func RemoveTabsFromStruct(s interface{}) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}
	removeTabs(val)
}

// removeTabs 递归清理值中的制表符
func removeTabs(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		// 如果是字符串，清理制表符
		if v.CanSet() {
			str := v.String()
			str = strings.ReplaceAll(str, "\t", "")
			v.SetString(str)
		}
	case reflect.Ptr:
		// 如果是指针，递归处理指向的值
		if !v.IsNil() {
			removeTabs(v.Elem())
		}
	case reflect.Struct:
		// 如果是结构体，递归处理每个字段
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			removeTabs(field)
		}
	case reflect.Slice, reflect.Array:
		// 如果是切片或数组，处理每个元素
		for i := 0; i < v.Len(); i++ {
			removeTabs(v.Index(i))
		}
	case reflect.Map:
		// 如果是map，处理每个值
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			// 创建一个新的可设置的值来处理
			newValue := reflect.New(value.Type()).Elem()
			newValue.Set(value)
			removeTabs(newValue)
			v.SetMapIndex(key, newValue)
		}
	}
	// 其他类型不做处理
}
