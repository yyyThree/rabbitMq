package test

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/yyyThree/rabbitmq"
	"testing"
)

func init()  {
	initConfig()
}

// 订阅直连交换机
// 消费成功
func TestSubscribeDirect(t *testing.T)  {
	err := rabbitmq.Subscribe(QueueNameDirect, func(msg amqp.Delivery) {
		fmt.Println("TestSubscribeDirect", msg.RoutingKey, string(msg.Body))
		rabbitmq.Ack(msg)
	})
	if err != nil {
		t.Fatal(err)
	}
}