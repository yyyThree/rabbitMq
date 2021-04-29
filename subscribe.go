package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

var publishClient *rabbitmq

// 发布消息
// 不保证100%发布成功
func Publish(routeKey, data string) (err error) {
	config := getConfig()
	if publishClient == nil || publishClient.conn.IsClosed() {
		publishClient, err = New(config)
		if err != nil {
			return err
		}
	}
	exchange := config.Exchange.Direct
	err = publishClient.channel.Publish(exchange, routeKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent, // 持久化
		Expiration: strconv.Itoa(config.Ttl.Msg), // 每条消息的有效期
		Body: []byte(data),
	})
	if err != nil {
		err = errors.New("Publish err:" + err.Error())
		return
	}
	return
}


// 发布消息
// 开启发布确认模式
func PublishWithConfirm(routeKey, data string) (err error) {
	config := getConfig()
	// 创建MQ连接
	if publishClient == nil || publishClient.conn.IsClosed() {
		fmt.Println("publishClient")
		publishClient, err = New(config)
		if err != nil {
			return
		}
	}

	// 发布确认模式下需要使用独立的channel
	channelWithConfirm, err := publishClient.conn.Channel()
	if err != nil {
		err = errors.New("PublishWithConfirm open channel error: " + err.Error())
		return
	}
	// 开启发布确认模式
	err = channelWithConfirm.Confirm(false)
	if err != nil {
		return errors.New("PublishWithConfirm open confirm mode err:" + err.Error())
	}
	// MQ服务器到达回调
	notifyConfirm := channelWithConfirm.NotifyPublish(make(chan amqp.Confirmation))
	// MQ服务器入列失败回调
	notifyReturn := channelWithConfirm.NotifyReturn(make(chan amqp.Return))
	// 协程处理入列失败回调通知
	go NotifyReturn(notifyReturn, channelWithConfirm)

	// mandatory：是否对无法路由的消息进行返回处理
	// 设置为true，用于监听入列失败回调
	// false，无法入列时直接丢弃消息
	exchange := config.Exchange.Direct
	fmt.Println("exchange", exchange)
	err = channelWithConfirm.Publish(exchange, routeKey, true, false, amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent, // 持久化
		Expiration: strconv.Itoa(config.Ttl.Msg), // 每条消息的有效期
		Body: []byte(data),
	})
	if err != nil {
		return errors.New("PublishWithConfirm err:" + err.Error())
	}

	// 阻塞获取到达结果
	select {
	case confirm := <-notifyConfirm:
		if confirm.Ack {
			log.Println("NotifyConfirm suc", routeKey, data, confirm)
			return
		}else {
			log.Println("NotifyConfirm fail", routeKey, data, confirm)
			return
		}
	}
}

// 发布确认 - 入列失败回调
func NotifyReturn(notifyReturn chan amqp.Return, channelWithConfirm *amqp.Channel) {
	// 7. 监听抵达Broker无误后的确认信息,设置5秒超时
	ticker := time.NewTicker(5 * time.Second)
	select {
	case returnChan := <-notifyReturn:
		fmt.Println("NotifyReturn", string(returnChan.Body), returnChan.RoutingKey)
	case <-ticker.C:
		fmt.Println("out of limit time!")
		_ = channelWithConfirm.Close()
	}
}