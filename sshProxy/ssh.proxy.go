package main

import (
	"time"
	ssh_proxy "workspace/many.program/code/ssh.proxy"
)

func main()  {
	ssh_proxy.Start()
	for{
		time.Sleep(1*time.Hour)
	}
}
