package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

var publisherClient *rabbitmq
var notifyConfirm chan amqp.Confirmation // 到达MQ服务器通知
var notifyReturn chan amqp.Return // 入列MQ服务器通知

// 发布消息
// 开启发布确认模式
func Publish(routeKey, data string) error {
	config := getConfig()
	if publisherClient == nil || publisherClient.conn.IsClosed() {
		fmt.Println(1111)
		client, err := New(config)
		if err != nil {
			return err
		}
		publisherClient = client
		// 开启发布确认模式
		err = publisherClient.channel.Confirm(false)
		if err != nil {
			return errors.New("Publish open confirm mode err:" + err.Error())
		}
		notifyConfirm = publisherClient.channel.NotifyPublish(make(chan amqp.Confirmation))
		notifyReturn = publisherClient.channel.NotifyReturn(make(chan amqp.Return))
	}

	// 到达成功回调
	go NotifyConfirm(notifyConfirm, routeKey, data)
	// 入列失败回调
	go NotifyReturn(notifyReturn, routeKey, data)

	// mandatory 设置为true，用于监听入列失败回调
	err := publisherClient.channel.Publish(config.Exchange.Direct, routeKey, true, false, amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent, // 持久化
		Expiration: strconv.Itoa(config.Ttl.Msg), // 每条消息的有效期
		Body: []byte(data),
	})
	if err != nil {
		return errors.New("Publish err:" + err.Error())
	}
	return nil
}

// 发布确认 - 到达服务器 回调
func NotifyConfirm(notifyConfirm chan amqp.Confirmation, routeKey, data string) {
	ret := <-notifyConfirm
	if ret.Ack {
		log.Println("NotifyConfirm suc")
	} else {
		log.Println("NotifyConfirm fail")
		fmt.Println(ret, routeKey, data)
	}
}

// 发布确认 - 入列 回调
func NotifyReturn(notifyReturn chan amqp.Return, routeKey, data string) {
	ret := <-notifyReturn
	fmt.Println("NotifyReturn", string(ret.Body), routeKey, data)
}
