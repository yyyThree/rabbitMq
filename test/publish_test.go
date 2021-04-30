package test

import (
	"github.com/yyyThree/rabbitmq"
	"testing"
	"time"
)

func init()  {
	initConfig()
}

// 测试普通发布
func TestPublish(t *testing.T)  {
	err := rabbitmq.Publish(QueueNameDirectKeys[0], "data1")
	if err != nil {
		t.Fatal(err)
	}
	err = rabbitmq.Publish(QueueNameDirectKeys[1], "data2")
	if err != nil {
		t.Fatal(err)
	}
}

// 测试带发布确认模式的发布
func TestPublishWithConfirm(t *testing.T)  {
	err := rabbitmq.PublishWithConfirm(QueueNameDirectKeys[0], "data3")
	if err != nil {
		t.Fatal(err)
	}
	err = rabbitmq.PublishWithConfirm(QueueNameWithDlKeys[0], "data4")
	if err != nil {
		t.Fatal(err)
	}
	// 发布错误的路由
	err = rabbitmq.PublishWithConfirm("errorRouteKey", "data5")
	if err != nil {
		t.Fatal(err)
	}
	// 发布至指定的交换机
	err = rabbitmq.PublishWithConfirm("errorRouteKey2", "data6", "go.topic")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * 1e9)
}