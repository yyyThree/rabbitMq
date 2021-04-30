package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Cb func(amqp.Delivery)

// 订阅队列
func Subscribe(queueName string, callback Cb) (err error) {
	client, err := New(config)
	if err != nil {
		subscribeFailLog(BaseMap{
			"queueName": queueName,
			"err": err,
		})
		return
	}

	msgList, err := client.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		subscribeFailLog(BaseMap{
			"queueName": queueName,
			"err": err,
		})
		return
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgList {
			callback(msg)
		}
	}()

	<-forever

	return
}

// 确认消费
func Ack(msg amqp.Delivery)  {
	_ = msg.Ack(false)
}

// 否认消费
func Nack(msg amqp.Delivery)  {
	_ = msg.Nack(false, false)
}

// 拒绝消费
func Reject(msg amqp.Delivery)  {
	_ = msg.Reject(false)
}
