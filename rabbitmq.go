package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type BaseMap map[string]interface{}

type rabbitmq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(configs ...Config) (*rabbitmq, error) {
	config := getConfig()
	if len(configs) > 0 {
		config = configs[0]
	}
	conn, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Base.User, config.Base.Password, config.Base.Host, config.Base.Port), amqp.Config{
		Vhost:     config.Base.Vhost,
		Heartbeat: 0 * time.Second, // 如果小于1会使用服务端默认的600s，amqp客户端已实现心跳机制
		Locale:    "en_US",
	})
	if err != nil {
		return nil, errors.New("rabbitmq connect err:" + err.Error())
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.New("rabbitmq get channel err:" + err.Error())
	}

	return &rabbitmq{
		conn:    conn,
		channel: channel,
	}, nil
}

func (rabbitmq *rabbitmq) Close() {
	_ = rabbitmq.channel.Close()
	_ = rabbitmq.conn.Close()
}
