package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
)

// 通用队列声明
func queueDeclare(name string) error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()


	_, err = client.channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-message-ttl": config.Ttl.QueueMsg,
	})
	if err != nil {
		return errors.New("queueDeclare err:" + err.Error())
	}

	return nil
}

// 直连交换机队列声明
// 绑定交换机和路由key
func QueueDeclareDirect(name string, routingKeys []string) (err error) {
	// 声明队列
	if err = queueDeclare(name); err != nil {
		return
	}

	// 绑定队列路由键值
	if err = QueueBind(name, routingKeys, config.Exchange.Direct); err != nil {
		return
	}

	return
}

// 主题交换机队列声明
// 绑定交换机和路由key
func QueueDeclareTopic(name string, routingKeys []string) (err error) {
	// 声明队列
	if err = queueDeclare(name); err != nil {
		return
	}

	// 绑定队列路由键值
	if err = QueueBind(name, routingKeys, config.Exchange.Topic); err != nil {
		return
	}

	return
}

// 带死信参数的直连交换机队列声明
// 绑定交换机和路由key
// dlRouteKey：死信消息发布至死信交换机使用的路由key
func QueueDeclareWithDl(name string, routingKeys []string, dlRouteKey string) (err error) {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	// 队列声明
	_, err = client.channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-message-ttl": config.Ttl.QueueMsg,
		"x-dead-letter-exchange": config.Exchange.DeathLetter,
		"x-dead-letter-routing-key": dlRouteKey,
	})
	if err != nil {
		return errors.New("QueueDeclareWithDl err:" + err.Error())
	}

	// 绑定队列路由键值
	if err = QueueBind(name, routingKeys, config.Exchange.Direct); err != nil {
		return
	}

	return
}

// 死信交换机队列声明
// 绑定交换机和路由key
func QueueDeclareDl(name string, routingKeys []string) (err error) {
	// 声明队列
	if err = queueDeclare(name); err != nil {
		return
	}

	// 绑定队列路由键值
	if err = QueueBind(name, routingKeys, config.Exchange.DeathLetter); err != nil {
		return
	}

	return
}

// 队列删除
// 存在消费者、剩余未消费消息的不允许删除
func QueueDelete(name string) error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	_, err = client.channel.QueueDelete(name, true, true, false)
	if err != nil {
		return errors.New("QueueDelete err:" + err.Error())
	}

	return nil
}

// 队列清空
// 清空未投递至消费者的消息
func QueuePurge(name string) error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	_, err = client.channel.QueuePurge(name, false)
	if err != nil {
		return errors.New("QueuePurge err:" + err.Error())
	}

	return nil
}

// 队列绑定路由
func QueueBind(name string, routingKeys []string, exchange string) error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	for _, routeKey := range routingKeys {
		err = client.channel.QueueBind(name, routeKey, exchange, false, nil)
		if err != nil {
			return errors.New("QueueBind err:" + err.Error())
		}
	}

	return nil
}

// 队列解绑路由
func QueueUnBind(name string, routingKeys []string, exchange string) error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	for _, routeKey := range routingKeys {
		err = client.channel.QueueUnbind(name, routeKey, exchange, nil)
		if err != nil {
			return errors.New("QueueUnBind err:" + err.Error())
		}
	}

	return nil
}