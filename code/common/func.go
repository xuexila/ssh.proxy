package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 打印错误信息
func Error(i ...interface{}) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	log.Println(i...)
}

// 打印正确日志。
func Log(i ...interface{}) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[用户日志]")
	log.SetOutput(os.Stdout)
	log.Println(i...)
}

// 获取绝对路径
func Fileabs(cpath string) string {
	if !strings.HasPrefix(cpath, "/") {
		cpath = filepath.Join(Appath, cpath)
	}
	return cpath
}

// 随机数生成器
func Urand() string{
	b:=make([]byte,1024)
	_,_=rand.Read(b)
	return Md5byte(b)
}

// md5 函数
func Md5byte(s []byte) string {
	h := md5.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

// 给字符串Md5
func Md5string(s string) string {
	return Md5byte([]byte(s))
}

// 带有 密钥的 sha1 hash
func CreateSignature (s ,key string) string{
	h:=hmac.New(sha1.New,[]byte(key))
	h.Write([]byte(s))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 判断目录是否存在，否则创建目录
func Mkdir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// 快速简易写文件
func FilePutContents(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	_ = file.Close()
	return err
}
