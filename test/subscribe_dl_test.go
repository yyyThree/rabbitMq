package test

import (
	"fmt"
	"github.com/streadway/amqp"
	rabbitmq "github.com/yyyThree/rabbitMq"
	"testing"
)

func init()  {
	initConfig()
}

// 订阅死信交换机
// 消费成功
func TestSubscribeDl(t *testing.T)  {
	err := rabbitmq.Subscribe(QueueNameDl, func(msg amqp.Delivery) {
		fmt.Println("TestSubscribeDl", msg.RoutingKey, string(msg.Body))
		rabbitmq.Ack(msg)
	})
	if err != nil {
		t.Fatal(err)
	}
}