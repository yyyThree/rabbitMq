package rabbitmq

import (
	"errors"
	"github.com/yyyThree/rabbitmq/helper"
	"github.com/yyyThree/rabbitmq/library/redis"
)

const defaultLogDir = "rabbitmqLog" // 默认文件存储文件夹地址

var config Config
var adminConfig Config

type Config struct {
	Base
	Exchange
	Ttl
	Admin
	Log
	Redis
}

// 基础设置
type Base struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
}

// 交换机类型
type Exchange struct {
	Direct      string // 基础直连交换机，业务系统默认使用
	Topic       string // 主题交换机
	DeathLetter string // 死信交换机
}

// 有效期管理
type Ttl struct {
	QueueMsg int // 队列中消息有效期，毫秒，默认为 86400 * 1e3
	Msg      int // 每条消息的有效期，毫秒，默认为 86400 * 1e3
}

// vhost对应的管理员账号，用于交换机、队列的声明
type Admin struct {
	User     string
	Password string
}

// 日志配置
type Log struct {
	Dir string // 日志存储文件夹地址，默认为 rabbitmqLog
}

// redis配置，用于消息消费幂等性
type Redis struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       int
}

// 初始化配置
func InitConfig(c Config) error {
	if helper.IsEmpty(c) {
		return errors.New("config error")
	}

	// 默认值处理
	initDefault(&c)

	config, adminConfig = c, c

	// 设置管理员账号
	adminConfig.Base.User = adminConfig.Admin.User
	adminConfig.Base.Password = adminConfig.Admin.Password

	// 初始化日志配置
	initLog(c.Log.Dir)

	// 初始化redis配置
	err := initRedis(c.Redis)
	if err != nil {
		return err
	}

	return nil
}

// 默认值设置
func initDefault(c *Config) {
	if c.Ttl.QueueMsg == 0 {
		c.Ttl.QueueMsg = 86400 * 1e3
	}
	if c.Ttl.Msg == 0 {
		c.Ttl.Msg = 86400 * 1e3
	}
	if c.Log.Dir == "" {
		c.Log.Dir = defaultLogDir
	}
}

func initRedis(r Redis) error {
	if helper.IsEmpty(r) {
		return nil
	}
	_, err := redis.GetConn(redis.Config{
		Host:     r.Host,
		Port:     r.Port,
		User:     r.User,
		Password: r.Password,
		Db:       r.Db,
	})
	if err != nil {
		return errors.New("Redis connect error: " + err.Error())
	}
	return nil
}

// 获取普通业务系统的账号配置，用于正常的业务消息发布、订阅
func getConfig() Config {
	return config
}

// 获取管理员账号配置，用于交换机、队列的处理
func getAdminConfig() Config {
	return adminConfig
}
