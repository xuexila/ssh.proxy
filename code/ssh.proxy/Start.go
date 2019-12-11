package ssh_proxy

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"time"
	"workspace/many.program/code/common"
)

func Start() {
	// fmt.Printf("%+v\n\n",conf.Addr)
	startAction(conf.Addr)
}

// 开始外部循环
func startAction(addrs []connect) {
	for _, item := range addrs {
		go func(item connect) {
			item.sshConfig=new(ssh.ClientConfig)
			item.config()
			common.Log("连接目标机器", item.Saddr)
			item.client = new(ssh.Client)
			item.client, err = ssh.Dial("tcp", item.Saddr, item.sshConfig)
			defer func() {
				if item.client != nil {
					_ = item.client.Close()
				}
			}()
			if err != nil {
				common.Error("连接远程服务器", item.Saddr, "失败", err.Error())
				os.Exit(1)
			}
			common.Log("服务器", item.Saddr, "连接成功")
			switch item.Connect {
			case "R": // 远程转发
				for {
					item.remoteForward()
					time.Sleep(time.Duration(retime) * time.Second)
				}
			case "L": // 本地转发
				for {
					item.socks5ProxyStart(item.localForward)
					time.Sleep(time.Duration(retime) * time.Second)
				}

			case "D": // 动态转发
				for {
					item.socks5ProxyStart(item.socks5Proxy)
					time.Sleep(time.Duration(retime) * time.Second)
				}

			default:

			}
		}(item)

	}
}

// 配置登陆配置
func (i *connect) config() {
	var auth ssh.AuthMethod
	if i.Stype == 2 {
		var pKey ssh.Signer
		b, err := ioutil.ReadFile(i.Passwd)
		if err != nil {
			common.Error(i.Saddr, "打开密钥文件失败", err.Error())
			return
		}
		pKey, err = ssh.ParsePrivateKey(b)
		if err != nil {
			common.Error("解析密钥文件失败", err.Error())
			return
		}
		auth = ssh.PublicKeys(pKey)
		i.sshConfig.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}
	} else {
		auth = ssh.Password(i.Passwd)
		i.sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	i.sshConfig.Auth = []ssh.AuthMethod{
		auth,
	}
	i.sshConfig.User = i.User
}
