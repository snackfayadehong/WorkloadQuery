package utity

import (
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"strings"
)

// RemoveAssignDir 工具包
// 清楚目录下指定后缀文件名
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
