package test

import (
	"github.com/yyyThree/rabbitMq"
	"testing"
)

func init()  {
	initConfig()
}

// 测试声明直连交换机
func TestExchangeDeclareDirect(t *testing.T)  {
	err := rabbitmq.ExchangeDeclareDirect()
	if err != nil {
		t.Fatal(err)
	}
}

// 测试声明主题交换机
func TestExchangeDeclareTopic(t *testing.T)  {
	err := rabbitmq.ExchangeDeclareTopic()
	if err != nil {
		t.Fatal(err)
	}
}

// 测试声明死信交换机
func TestExchangeDeclareDl(t *testing.T)  {
	err := rabbitmq.ExchangeDeclareDl()
	if err != nil {
		t.Fatal(err)
	}
}

// 测试删除交换机
func TestExchangeDelete(t *testing.T)  {
	err := rabbitmq.ExchangeDelete("go.topic")
	if err != nil {
		t.Fatal(err)
	}
	// 重新声明
	_ = rabbitmq.ExchangeDeclareTopic()
}