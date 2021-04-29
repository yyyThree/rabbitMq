package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
)

// 队列声明，绑定交换机和路由key
func QueueDeclare(name string, routingKeys []string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()


	_, err = rabbitmq.channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-message-ttl": config.ttl.queueMsg,
	})
	if err != nil {
		return errors.New("QueueDeclare err:" + err.Error())
	}
	for _, routeKey := range routingKeys {
		err = rabbitmq.channel.QueueBind(name, routeKey, config.exchange.direct, false, nil)
		if err != nil {
			return errors.New("QueueBind err:" + err.Error())
		}
	}
	return nil
}

// 带死信参数的队列声明，绑定交换机和路由key
// dlRouteKey 死信消息发布至死信交换机使用的路由key
func QueueDeclareWithDl(name string, routingKeys []string, dlRouteKey string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()


	_, err = rabbitmq.channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-message-ttl": config.ttl.queueMsg,
		"x-dead-letter-exchange": config.exchange.deathLetter,
		"x-dead-letter-routing-key": dlRouteKey,
	})
	if err != nil {
		return errors.New("QueueDeclare err:" + err.Error())
	}
	for _, routeKey := range routingKeys {
		err = rabbitmq.channel.QueueBind(name, routeKey, config.exchange.direct, false, nil)
		if err != nil {
			return errors.New("QueueBind err:" + err.Error())
		}
	}
	return nil
}

// 死信队列声明，绑定交换机和路由key
func DlQueueDeclare(name string, dlRouteKeys []string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()


	_, err = rabbitmq.channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-message-ttl": config.ttl.queueMsg,
	})
	if err != nil {
		return errors.New("DlQueueDeclare err:" + err.Error())
	}
	for _, dlRouteKey := range dlRouteKeys {
		err = rabbitmq.channel.QueueBind(name, dlRouteKey, config.exchange.deathLetter, false, nil)
		if err != nil {
			return errors.New("DlQueueBind err:" + err.Error())
		}
	}
	return nil
}

// 队列删除
// 存在消费者、剩余未消费消息的不允许删除
func QueueDelete(name string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()

	_, err = rabbitmq.channel.QueueDelete(name, true, true, false)
	if err != nil {
		return errors.New("QueueDelete err:" + err.Error())
	}

	return nil
}

// 队列清空
// 清空未投递至消费者的消息
func QueuePurge(name string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()

	_, err = rabbitmq.channel.QueuePurge(name, false)
	if err != nil {
		return errors.New("QueuePurge err:" + err.Error())
	}

	return nil
}

// 队列绑定路由
func QueueBind(name string, routingKeys []string, exchange string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()

	for _, routeKey := range routingKeys {
		err = rabbitmq.channel.QueueBind(name, routeKey, exchange, false, nil)
		if err != nil {
			return errors.New("QueueBind err:" + err.Error())
		}
	}

	return nil
}

// 队列解绑路由
func QueueUnBind(name string, routingKeys []string, exchange string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()

	for _, routeKey := range routingKeys {
		err = rabbitmq.channel.QueueUnbind(name, routeKey, exchange, nil)
		if err != nil {
			return errors.New("QueueUnBind err:" + err.Error())
		}
	}

	return nil
}