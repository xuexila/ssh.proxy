package common

import (
	"os"
	"path/filepath"
)

func init(){
	Appath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err!=nil {
		Error("当前路径获取失败...",err.Error())
		os.Exit(1)
	}
}
