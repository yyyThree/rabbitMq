package rabbitmq

import "errors"

// 声明三种交换机
// 直连交换机 - 业务中主要使用的交换机
// 主题交换机 - 特殊业务可用
// 死信交换机 - 本质也是直连交换机，用于消息未处理完成后的回调
func ExchangeDeclare() error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()

	// 声明直连交换机
	err = rabbitmq.channel.ExchangeDeclare(config.exchange.direct, "direct", true, false, false, false, nil)
	if err != nil {
		return errors.New("ExchangeDeclare direct err:" + err.Error())
	}
	// 声明主题交换机
	err = rabbitmq.channel.ExchangeDeclare(config.exchange.topic, "topic", true, false, false, false, nil)
	if err != nil {
		return errors.New("ExchangeDeclare topic err:" + err.Error())
	}
	// 声明死信交换机
	err = rabbitmq.channel.ExchangeDeclare(config.exchange.deathLetter, "direct", true, false, false, false, nil)
	if err != nil {
		return errors.New("ExchangeDeclare deathLetter err:" + err.Error())
	}

	_ = rabbitmq.conn.Close()
	return nil
}

// 删除某个交换机，如果交换机已绑定队列则不允许删除
func ExchangeDelete(exchange string) error {
	config := getAdminConfig()
	rabbitmq, err := New(config)
	if err != nil {
		return err
	}
	defer func() {
		rabbitmq.Close()
	}()

	err = rabbitmq.channel.ExchangeDelete(exchange, true, false)
	if err != nil {
		return errors.New("ExchangeDelete fail" + err.Error())
	}
	return nil
}
