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

// 订阅带死信参数的直连交换机队列
// 否认消费
func TestSubscribeDirectWithDl(t *testing.T)  {
	err := rabbitmq.Subscribe(QueueNameWithDl, func(msg amqp.Delivery) {
		fmt.Println("TestSubscribeDirectWithDl", msg.RoutingKey, string(msg.Body))
		rabbitmq.Nack(msg)
	})
	if err != nil {
		t.Fatal(err)
	}
}