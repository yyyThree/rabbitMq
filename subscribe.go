package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/yyyThree/rabbitmq/helper"
	"github.com/yyyThree/rabbitmq/library/redis"
	"time"
)

const redisPrefix = "Rabbitmq_msg_"

type Cb func(amqp.Delivery)

// 订阅队列
func Subscribe(queueName string, callback Cb) (err error) {
	msgList, err := subscribe(queueName)
	if err != nil {
		return
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgList {
			fmt.Println(111, msg.MessageId)
			callback(msg)
		}
	}()

	<-forever

	return
}

// 订阅队列
// 幂等性消费
func SubscribeIdp(queueName string, callback Cb) (err error) {
	msgList, err := subscribe(queueName)
	if err != nil {
		return
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgList {
			fmt.Println(222, msg.MessageId)
			// 校验消息是否已被消费
			if redis.Client != nil && msg.MessageId != "" {
				key := redisPrefix + msg.MessageId
				if redis.Client.Get(redis.GetCtx(), key).Val() != "" {
					fmt.Println("已被消费", msg, msg.MessageId)
					Reject(msg)
					continue
				}
			}

			callback(msg)
		}
	}()

	<-forever

	return
}

// 订阅队列
func subscribe(queueName string) (msgList <-chan amqp.Delivery, err error) {
	client, err := New(config)
	if err != nil {
		subscribeFailLog(BaseMap{
			"queueName": queueName,
			"err": err,
		})
		return
	}

	msgList, err = client.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	return msgList, err
}

// 确认消费
func Ack(msg amqp.Delivery)  {
	_ = msg.Ack(false)
}

// 幂等确认消费
func AckIdp(msg amqp.Delivery)  {
	_ = msg.Ack(false)
	fmt.Println(22222, redis.Client)
	if redis.Client != nil {
		config = getConfig()
		key := redisPrefix + msg.MessageId
		data := helper.StructToJson(msg)
		res := redis.Client.SetNX(redis.GetCtx(), key, data, time.Duration(config.Ttl.Msg)*time.Microsecond)
		fmt.Println(333, key, res.Val(), res.Err(), string(data))
	}
}

// 否认消费
func Nack(msg amqp.Delivery)  {
	_ = msg.Nack(false, false)
}

// 拒绝消费
func Reject(msg amqp.Delivery)  {
	_ = msg.Reject(false)
}
