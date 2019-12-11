package ssh_proxy

var (
	conf    = new(config) // 系统配置
	retime  int64         // 连接异常次数后重新打开的等待时间
	err     error         // 全局错误
)
