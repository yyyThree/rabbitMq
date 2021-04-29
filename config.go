package rabbitmq

import ginConfig "gin/config"

type Config struct {
	base
	exchange
	ttl
	admin
}

type base struct {
	host     string
	port     int
	user     string
	password string
	vhost    string
}

type exchange struct {
	direct      string // 基础直连交换机，业务系统默认使用
	topic       string // 主题交换机
	deathLetter string // 死信交换机
}

// 有效期
type ttl struct {
	queueMsg int // 队列中消息有效期，毫秒
	msg      int // 每条消息的有效期，毫秒
}

// vhost对应的管理员账号，用于交换机、队列的声明
type admin struct {
	user     string
	password string
}

// 获取普通业务系统的账号配置，用于正常的业务消息发布、订阅
func getConfig() *Config {
	return &Config{
		base: base{
			host:     ginConfig.Config.Rabbitmq.Host,
			port:     ginConfig.Config.Rabbitmq.Port,
			user:     ginConfig.Config.Rabbitmq.User,
			password: ginConfig.Config.Rabbitmq.Password,
			vhost:    ginConfig.Config.Rabbitmq.Vhost,
		},
		exchange: exchange{
			direct:      ginConfig.Config.Rabbitmq.ExDirect,
			topic:       ginConfig.Config.Rabbitmq.ExTopic,
			deathLetter: ginConfig.Config.Rabbitmq.ExDeathLetter,
		},
		ttl: ttl{
			queueMsg: ginConfig.Config.Rabbitmq.TtlQueueMsg * 1e3,
			msg:      ginConfig.Config.Rabbitmq.TtlMsg * 1e3,
		},
	}
}

// 获取管理员账号配置，用于交换机、队列的处理
func getAdminConfig() *Config {
	return &Config{
		base: base{
			host:     ginConfig.Config.Rabbitmq.Host,
			port:     ginConfig.Config.Rabbitmq.Port,
			user:     ginConfig.Config.Rabbitmq.AdminUser,
			password: ginConfig.Config.Rabbitmq.AdminPassword,
			vhost:    ginConfig.Config.Rabbitmq.Vhost,
		},
		exchange: exchange{
			direct:      ginConfig.Config.Rabbitmq.ExDirect,
			topic:       ginConfig.Config.Rabbitmq.ExTopic,
			deathLetter: ginConfig.Config.Rabbitmq.ExDeathLetter,
		},
		ttl: ttl{
			queueMsg: ginConfig.Config.Rabbitmq.TtlQueueMsg * 1e3,
			msg:      ginConfig.Config.Rabbitmq.TtlMsg * 1e3,
		},
	}
}
