package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/yyyThree/rabbitmq/helper"
	"github.com/yyyThree/rabbitmq/library/redis"
	"time"
)

const redisPrefix = "Rabbitmq_msg_"

type Cb func(amqp.Delivery)

// 订阅队列
func Subscribe(queueName string, callback Cb) (err error) {
	msgList, closeChan, err := subscribe(queueName)
	if err != nil {
		return
	}

	forever := make(chan bool)

	go func() {
		for {
			isClosed := false
			select {
			case msg := <-msgList:
				callback(msg)
			case e := <-closeChan:
				subscribeConnErrorLog(BaseMap{
					"queueName": queueName,
					"err": e,
				})
				_ = Subscribe(queueName, callback)
				isClosed = true
			}
			if isClosed {
				break
			}
		}
	}()

	<-forever

	return
}

// 订阅队列
// 幂等性消费
func SubscribeIdp(queueName string, callback Cb) (err error) {
	msgList, closeChan, err := subscribe(queueName)
	if err != nil {
		return
	}

	forever := make(chan bool)

	go func() {
		for {
			isClosed := false
			select {
			case msg := <-msgList:
				// 校验消息是否已被消费
				if !checkMsgIdp(msg) {
					subscribeIdpFailLog(BaseMap{
						"queueName": queueName,
						"msg": BaseMap{
							"MessageId": msg.MessageId,
							"Exchange": msg.Exchange,
							"RoutingKey": msg.RoutingKey,
							"Body": string(msg.Body),
						},
					})
					Reject(msg)
					continue
				}
				callback(msg)
			case e := <-closeChan:
				subscribeConnErrorLog(BaseMap{
					"queueName": queueName,
					"err": e,
				})
				_ = SubscribeIdp(queueName, callback)
				isClosed = true
			}
			if isClosed {
				break
			}
		}
	}()

	<-forever

	return
}

// 校验消息幂等性
func checkMsgIdp(msg amqp.Delivery) bool {
	if redis.Client == nil || msg.MessageId == "" {
		return true
	}
	key := redisPrefix + msg.MessageId
	if redis.Client.Get(redis.GetCtx(), key).Val() == "" {
		return true
	}

	return false
}

// 订阅队列
func subscribe(queueName string) (msgList <-chan amqp.Delivery, notifyClose chan *amqp.Error, err error) {
	client, err := New(config)
	if err != nil {
		subscribeFailLog(BaseMap{
			"queueName": queueName,
			"err": err,
		})
		return
	}
	// 通道关闭通知
	notifyClose = client.channel.NotifyClose(make(chan *amqp.Error))

	msgList, err = client.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	return msgList, notifyClose, err
}

// 确认消费
func Ack(msg amqp.Delivery)  {
	_ = msg.Ack(false)
}

// 幂等确认消费
func AckIdp(msg amqp.Delivery)  {
	_ = msg.Ack(false)
	if redis.Client != nil {
		config = getConfig()
		key := redisPrefix + msg.MessageId
		data := helper.MapToJson(BaseMap{
			"MessageId": msg.MessageId,
			"Exchange": msg.Exchange,
			"RoutingKey": msg.RoutingKey,
			"Body": string(msg.Body),
		})
		redis.Client.SetNX(redis.GetCtx(), key, data, time.Duration(config.Ttl.Msg)*time.Millisecond)
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
