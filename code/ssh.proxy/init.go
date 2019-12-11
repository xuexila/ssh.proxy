package ssh_proxy

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"workspace/many.program/code/common"
)

func init(){
	initParams()
	loadConfig()
}

func loadConfig(){
	by,err:=ioutil.ReadFile(common.Cpath)
	if err!=nil {
		common.Error("读取配置文件失败",err.Error())
		os.Exit(1)
	}
	if err=json.Unmarshal(by,conf);err!=nil {
		common.Error("解析配置文件失败",err.Error())
		os.Exit(1)
	}
	if len(conf.Addr)<1 {
		common.Error("无代理配置")
		os.Exit(1)
	}
}

func initParams(){
	flag.BoolVar(&common.H, "h", false, "帮助说明")
	flag.Int64Var(&retime,"r",10,"连接中断后重新连接的等待时间")
	flag.StringVar(&common.Cpath, "c", "conf.json", "系统配置文件")
	flag.Parse()
	if common.H {
		flag.Usage()
		os.Exit(0)
	}
	common.Fileabs(common.Cpath)
	common.Log("参数解析完成")
}