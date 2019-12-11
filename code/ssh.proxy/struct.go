package ssh_proxy

import "golang.org/x/crypto/ssh"

type sockIP struct {
	A, B, C, D byte
	PORT       uint16
}

// 系统配置
type config struct {
	Addr []connect
}

type connect struct {
	Saddr   string      `json:"saddr"`   // 目标地址
	User    string      `json:"user"`    // 用户
	Stype   int         `json:"type"`    // 验证方式 1、密码验证 2、密钥验证
	Passwd  string      `json:"passwd"`  // 密码、密钥路径
	Remote  string      `json:"remote"`  // 远程地址
	Listen  string      `json:"listen"`  // 本地监听地址
	Connect string      `json:"connect"` // 连接类型 L 本地转发 R 远程转发 D 动态转发
	Son     []connect   `json:"son"`     // 子连接
	client  *ssh.Client // ssh 客户端
	sshConfig *ssh.ClientConfig // 连接ssh 的配置

}
