package test

import (
	rabbitmq "github.com/yyyThree/rabbitMq"
	"testing"
)

func TestInit(t *testing.T) {
	initConfig()
}

// 直连交换机队列声明
func TestQueueDeclareDirect(t *testing.T) {
	err := rabbitmq.QueueDeclareDirect(QueueNameDirect, QueueNameDirectKeys)
	if err != nil {
		t.Fatal(err)
	}
}

// 主题交换机队列声明
func TestQueueDeclareTopic(t *testing.T) {
	err := rabbitmq.QueueDeclareTopic(QueueNameTopic, QueueNameTopicKeys)
	if err != nil {
		t.Fatal(err)
	}
}

// 带死信参数的直连交换机队列声明
func TestQueueDeclareWithDl(t *testing.T) {
	err := rabbitmq.QueueDeclareWithDl(QueueNameWithDl, QueueNameWithDlKeys, QueueNameDlKey)
	if err != nil {
		t.Fatal(err)
	}
}

// 死信交换机队列声明
func TestQueueDeclareDl(t *testing.T) {
	err := rabbitmq.QueueDeclareDl(QueueNameDl, []string{QueueNameDlKey})
	if err != nil {
		t.Fatal(err)
	}
}

// 队列删除
func TestQueueDelete(t *testing.T) {
	err := rabbitmq.QueueDelete(QueueNameTopic)
	if err != nil {
		t.Fatal(err)
	}
}

// 队列清空
func TestQueuePurge(t *testing.T) {
	err := rabbitmq.QueuePurge(QueueNameDirect)
	if err != nil {
		t.Fatal(err)
	}
}

// 队列解绑路由键值
func TestQueueUnBind(t *testing.T) {
	err := rabbitmq.QueueUnBind(QueueNameDirect, QueueNameDirectKeys[1:], "go.direct")
	if err != nil {
		t.Fatal(err)
	}
}