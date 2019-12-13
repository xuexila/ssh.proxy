package ssh_proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
	"workspace/many.program/code/common"
)

func (i *connect) socks5Proxy(conn net.Conn) {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	var b [1024]byte
	n, err := conn.Read(b[:])
	if err != nil {
		if err!=io.EOF {
			common.Error(i.Saddr,"目标数据读取失败", err.Error())
		}

		return
	}
	// if len(b[:])>0 {
	// 	Log("数据：",string(b[:]))
	// }
	// log.Printf("% x", b[:n])
	_, _ = conn.Write([]byte{0x05, 0x00})
	n, err = conn.Read(b[:])
	if err != nil{
		if err!=io.EOF  {
			common.Error("第一次添加数据",i.Saddr,err.Error())
		}

		return
	}
	// log.Printf("% x", b[:n])

	var addr string
	switch b[3] {
	case 0x01:
		sip := sockIP{}
		if err := binary.Read(bytes.NewReader(b[4:n]), binary.BigEndian, &sip); err != nil {
			common.Error(i.Saddr,"请求解析错误", err.Error())
			return
		}
		addr = sip.toAddr()
	case 0x03:
		host := string(b[5 : n-2])
		var port uint16
		err = binary.Read(bytes.NewReader(b[n-2:n]), binary.BigEndian, &port)
		if err != nil {
			common.Error(err)
			return
		}
		addr = fmt.Sprintf("%s:%d", host, port)
	}
	if i.client==nil {
		i.login(true)
	}
	server, err := i.client.Dial("tcp", addr)
	defer func() {
		if server!=nil {
			_=server.Close()
		}
	}()

	if err != nil {
		common.Error("动态转发连接目标失败",err.Error())
		return
	}
	_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	go func() {
		_, _ = io.Copy(server, conn)
	}()
	_, _ = io.Copy(conn, server)
}
// 本地转发
func (i *connect) localForward(conn net.Conn){
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	if i.client==nil {
		i.login(true)
	}
	server, err := i.client.Dial("tcp", i.Remote)
	defer func() {
		if server!=nil {
			_=server.Close()
		}
	}()
	if err != nil {
		common.Error("本地转发，远程数据传输异常",err.Error())
		return
	}
	go func() {
		_, _ = io.Copy(server, conn)
	}()
	_, _ = io.Copy(conn, server)
}

// 本地转发，和动态转发通用函数
func (i *connect) socks5ProxyStart(fun func(conn net.Conn)) {
	common.Log("本地端口监听...",i.Listen)
	localServer, err := net.Listen("tcp", i.Listen)
	defer func() {
		if localServer != nil {
			_ = localServer.Close()
		}
	}()
	if err != nil {
		common.Error(i.Saddr,"代理端口监听失败", err.Error())
		return
	}
	common.Log("本地端口监听成功...",i.Listen)
	if len(i.Son)>0{
		go func() {
			time.Sleep(1 * time.Second)
			common.Log("开始创建子系统")
			startAction(i.Son)
		}()
	}

	for {
		client, err := localServer.Accept()
		if err != nil {
			common.Error(i.Saddr,"tcp 数据获取失败", err.Error())
			return
		}
		go fun(client)

	}

}

func (ip sockIP) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip.A, ip.B, ip.C, ip.D, ip.PORT)
}


// 远程转发
func (i *connect) remoteForward(){
	common.Log("远程转发端口监听",i.Remote)
	if i.client==nil {
		i.login(true)
	}
	server, err :=i.client.Listen("tcp",i.Remote)
	if err != nil {
		common.Error("建立远程端口失败",i.Remote,err.Error())
		return
	}
	common.Log("远程转发端口监听成功",i.Remote)

	for{
		client, err := server.Accept()
		if err != nil {
			common.Error(i.Remote,"TCP 远程数据接受失败", err.Error())
			return
		}
		go func(conn net.Conn) {
			defer func() {
				if conn != nil {
					_ = conn.Close()
				}
			}()
			server, err := net.Dial("tcp", i.Listen)
			if err != nil {
				common.Error(err.Error())
				return
			}
			go func() {
				_, _ = io.Copy(server, conn)
			}()
			_, _ = io.Copy(conn, server)
		}(client)


	}

}
