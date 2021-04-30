package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
	"github.com/yyyThree/rabbitmq/helper"
	"log"
	"strconv"
	"time"
)

var publishClient *rabbitmq

// 发布消息
// 不保证100%发布成功
// 复用conn和channel
// 默认发布至直连交换机
func Publish(routeKey, data string, exchanges ...string) (err error) {
	exchange := config.Exchange.Direct
	if len(exchanges) > 0 {
		exchange = exchanges[0]
	}

	config := getConfig()
	// TODO amqp包暂未提供校验channel状态的方法
	if publishClient == nil || publishClient.conn.IsClosed() {
		publishClient, err = New(config)
		if err != nil {
			publishFailLog(BaseMap{
				"routeKey": routeKey,
				"data": data,
				"err": err,
			})
			return err
		}
	}
	err = publishClient.channel.Publish(exchange, routeKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent, // 持久化
		Expiration: strconv.Itoa(config.Ttl.Msg), // 每条消息的有效期
		Body: []byte(data),
	})
	if err != nil {
		err = errors.New("Publish err:" + err.Error())
		publishFailLog(BaseMap{
			"routeKey": routeKey,
			"data": data,
			"err": err,
		})
		return
	}
	return
}


// 发布消息
// 开启发布确认模式
// 使用独立的conn和channel
// 默认发布至直连交换机
func PublishWithConfirm(routeKey, data string, exchanges ...string) (err error) {
	exchange := config.Exchange.Direct
	if len(exchanges) > 0 {
		exchange = exchanges[0]
	}

	config := getConfig()
	// 创建MQ连接
	client, err := New(config)
	if err != nil {
		publishFailLog(BaseMap{
			"routeKey": routeKey,
			"data": data,
			"err": err,
		})
		return
	}

	// 开启发布确认模式
	err = client.channel.Confirm(false)
	if err != nil {
		err = errors.New("PublishWithConfirm open confirm mode err:" + err.Error())
		publishFailLog(BaseMap{
			"routeKey": routeKey,
			"data": data,
			"err": err,
		})
		return
	}

	// 设置MQ服务器到达回调
	notifyConfirm := client.channel.NotifyPublish(make(chan amqp.Confirmation))
	// 设置MQ服务器入列失败回调
	notifyReturn := client.channel.NotifyReturn(make(chan amqp.Return))
	// 协程处理入列失败回调通知
	go NotifyReturn(notifyReturn, client)

	// mandatory：是否对无法路由的消息进行返回处理
	// 设置为true，用于监听入列失败回调
	// false，无法入列时直接丢弃消息
	err = client.channel.Publish(exchange, routeKey, true, false, amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent, // 持久化
		Expiration: strconv.Itoa(config.Ttl.Msg), // 每条消息的有效期
		Body: []byte(data),
	})
	if err != nil {
		err = errors.New("PublishWithConfirm err:" + err.Error())
		publishFailLog(BaseMap{
			"routeKey": routeKey,
			"data": data,
			"err": err,
		})
		return
	}

	// 阻塞获取到达结果
	select {
	case confirm := <-notifyConfirm:
		if confirm.Ack {
			log.Println("NotifyConfirm suc", routeKey, data, confirm)
			return
		}else {
			publishConfirmFailLog(BaseMap{
				"routeKey": routeKey,
				"data": data,
			})
			return
		}
	}
}

// 发布确认 - 入列失败回调
func NotifyReturn(notifyReturn chan amqp.Return, client *rabbitmq) {
	// 设置5s超时
	// 在超时时间内读取到入列失败通知，记录处理并关闭连接
	// 触发超时时间后仍未读取到通知，关闭连接
	ticker := time.NewTicker(5 * time.Second)
	select {
	case returnChan := <-notifyReturn:
		if !helper.IsEmpty(returnChan) {
			publishReturnFailLog(BaseMap{
				"exchange": returnChan.Exchange,
				"routeKey": returnChan.RoutingKey,
				"data": string(returnChan.Body),
			})
		}
	case <-ticker.C:
	}
	client.Close()
}