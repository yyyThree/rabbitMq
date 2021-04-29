package test

import (
	"fmt"
	rabbitmq "github.com/yyyThree/rabbitMq"
)


const QueueNameDirect = "queueDirect"

var QueueNameDirectKeys = []string{"queueDirectKey1", "queueDirectKey2"}

const QueueNameTopic = "queueTopic"

var QueueNameTopicKeys = []string{"queueTopicKey1", "queueTopicKey2"}

const QueueNameWithDl = "queueWithDl"

var QueueNameWithDlKeys = []string{"queueWithDlKey1", "queueWithDlKey2"}

const QueueNameDl = "queueDl"

var QueueNameDlKey = "queueDlKey"

// 初始化MQ配置
func initConfig() {
	fmt.Println("initConfig")
	err := rabbitmq.InitConfig(rabbitmq.Config{
		Base: rabbitmq.Base {
			Host: "192.168.3.53",
			Port: 5673,
			User: "go",
			Password: "go",
			Vhost: "go",
		},
		Exchange: rabbitmq.Exchange {
			Direct: "go.direct",
			Topic: "go.topic",
			DeathLetter: "go.dl",
		},
		Ttl: rabbitmq.Ttl {
			QueueMsg: 86400 * 1e3,
			Msg: 86400 * 1e3,
		},
		Admin: rabbitmq.Admin {
			User: "goadmin",
			Password: "goadmin",
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}