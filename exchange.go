package rabbitmq

import (
	"errors"
)

// 声明直连交换机
// 业务中主要使用的交换机
func ExchangeDeclareDirect() error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	// 声明直连交换机
	err = client.channel.ExchangeDeclare(config.Exchange.Direct, "direct", true, false, false, false, nil)
	if err != nil {
		return errors.New("ExchangeDeclare direct err:" + err.Error())
	}

	return nil
}

// 声明主题交换机
// 特殊业务可用
func ExchangeDeclareTopic() error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	// 声明主题交换机
	err = client.channel.ExchangeDeclare(config.Exchange.Topic, "topic", true, false, false, false, nil)
	if err != nil {
		return errors.New("ExchangeDeclare topic err:" + err.Error())
	}

	return nil
}

// 声明死信交换机
// 本质也是直连交换机，用于消息未处理完成后的回调
func ExchangeDeclareDl() error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	// 声明死信交换机
	err = client.channel.ExchangeDeclare(config.Exchange.DeathLetter, "direct", true, false, false, false, nil)
	if err != nil {
		return errors.New("ExchangeDeclare deathLetter err:" + err.Error())
	}

	return nil
}

// 删除某个交换机
// 如果交换机已绑定队列则不允许删除
func ExchangeDelete(exchange string) error {
	config := getAdminConfig()
	client, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	err = client.channel.ExchangeDelete(exchange, true, false)
	if err != nil {
		return errors.New("ExchangeDelete fail" + err.Error())
	}
	return nil
}
