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
// 幂等消费
// 消费成功
func TestSubscribeDirectIdp(t *testing.T)  {
	err := rabbitmq.SubscribeIdp(QueueNameDirect, func(msg amqp.Delivery) {
		fmt.Println("TestSubscribeDirectIdp", msg.RoutingKey, string(msg.Body))
		rabbitmq.AckIdp(msg)
	})
	if err != nil {
		t.Fatal(err)
	}
}