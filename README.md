# go语言rabbitmq客户端
1. 交换机
    - 支持多种交换机：直连交换机/主题交换机/死信交换机
    - 支持交换机创建/删除
    - 支持交换机持久化
2. 队列
    - 支持普通队列/死信队列
    - 支持队列创建/删除/清空/绑定路由与交换机/解绑路由与交换机
    - 支持队列持久化
3. 消息发布
    - 支持普通消息发布（不保证100%发布成功）/带发布确认模式的消息发布（100%发布成功）
    - 支持发布至指定交换机 
    - 支持消息发布失败记录/MQ服务器到达失败记录/MQ服务器入列失败记录
    - 支持消息持久化 
4. 消息消费
    - 支持订阅队列，传入回调函数进行消费
    - 支持多种消息消费确认：ack/nack/reject
    - 支持幂等性消费  
    - 支持订阅队列失败记录/幂等消费失败记录

## 一、安装与引用
   
   ```go
   go get github.com/yyyThree/rabbitmq
   import "github.com/yyyThree/rabbitmq"
   ```

## 二、使用说明
1. 加载配置
   ```go
    rabbitmq.InitConfig(rabbitmq.Config{
        Base: rabbitmq.Base {
            Host: "127.0.0.1",
            Port: 5673,
            User: "go",
            Password: "go",
            Vhost: "go",
        },
        Exchange: rabbitmq.Exchange {
            Direct: "go.direct", // 基础直连交换机，业务系统默认使用
            Topic: "go.topic", // 主题交换机
            DeathLetter: "go.dl", // 死信交换机
        },
        Ttl: rabbitmq.Ttl {
            QueueMsg: 86400 * 1e3, // 队列中消息有效期，毫秒，默认为 86400 * 1e3
            Msg: 86400 * 1e3, // 每条消息的有效期，毫秒，默认为 86400 * 1e3
        },
        // vhost对应的管理员账号，用于交换机、队列的声明
        Admin: rabbitmq.Admin {
            User: "goadmin",
            Password: "goadmin",
        },
        Log: rabbitmq.Log{
            Dir: "rabbitmqLog" // 日志存储地址，默认为rabbitmqLog
        }, 
    })
   ```
2. 交换机
   1. 声明
      ```go
      // 直连交换机
      rabbitmq.ExchangeDeclareDirect()
      
      // 主题交换机
      rabbitmq.ExchangeDeclareTopic()
      
      // 死信交换机
      rabbitmq.ExchangeDeclareDl()
      ```
   2. 删除
      ```go
      rabbitmq.ExchangeDelete(exchangeName)
      ```
3. 队列
   1. 声明
      ```go
      // 直连交换机队列
      rabbitmq.QueueDeclareDirect("queueDirect", []string{"queueDirectKey1", "queueDirectKey2"})
      
      // 主题交换机队列
      rabbitmq.QueueDeclareTopic("queueTopic", []string{"queueTopicKey1", "queueTopicKey2"})
      
      // 带死信参数的直连交换机队列
      rabbitmq.QueueDeclareWithDl("queueWithDl", []string{"queueWithDlKey1", "queueWithDlKey2"}, "queueDlKey")
      
      // 死信交换机队列
      rabbitmq.QueueDeclareDl("queueDl", []string{"queueDlKey"})
      ```
   2. 删除
      ```go
      rabbitmq.QueueDelete(QueueName)
      ```
   3. 清空
      ```go
      rabbitmq.QueuePurge(QueueName)
      ```
   4. 解绑路由
      ```go
      rabbitmq.QueueUnBind(QueueName, []string{"queueDirectKey2"}, exchangeName)
      ```   
4. 消息发布
   1. 普通发布（不保证100%发布成功）
      ```go
      // 默认发布至直连交换机
      rabbitmq.Publish("queueDirectKey1", "data1")
      
      // 指定交换机
      rabbitmq.Publish("queueDirectKey1", "data1", exchangeName)
      ```   
   2. 带发布确认模式的发布（100%发布成功）
      ```go
      // 默认发布至直连交换机
      rabbitmq.PublishWithConfirm("errorRouteKey", "data2")

      // 指定交换机
      rabbitmq.PublishWithConfirm("errorRouteKey", "data2", exchangeName)
      ```   
   3. 错误日志
      - `{Log.Dir}/publishFail/*/*.log`: 消息发布失败记录
      - `{Log.Dir}/publishConfirmFail/*`: MQ服务器到达失败记录
      - `{Log.Dir}/publishReturnFail/*`: MQ服务器入列失败记录
5. 消息订阅
   1. 订阅队列
      ```go
      rabbitmq.Subscribe(QueueName, func(msg amqp.Delivery) {
        doSomething()
        
        // 确认消费
        rabbitmq.Ack(msg)
        
        // 否认消费，自动发布至对应的死信信息队列（如果存在）
        rabbitmq.Ack(msg)      
        
        // 拒绝消费，自动发布至对应的死信信息队列（如果存在）
        rabbitmq.Reject(msg)         
      })      
      ```      
   2. 幂等性消费
      ```go
      rabbitmq.SubscribeIdp(QueueName, func(msg amqp.Delivery) {
        doSomething()

        rabbitmq.Ack(msg)
      })      
      ```
      自动校验消息幂等性，保证每条消息仅被消费一次。如果幂等校验不通过，消息将自动被`reject`并记录日志。         
   3. 错误日志
      - `{Log.Dir}/subscribeFail/*/*.log`: 订阅队列失败记录/幂等消费失败记录

## 三、下一步开发计划
1. rabbitmq连接、channel异常的监听和后续处理