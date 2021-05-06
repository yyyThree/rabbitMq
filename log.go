package rabbitmq

import (
	"github.com/yyyThree/zap"
	"strings"
)

var (
	publishFailLogger *zap.Logger // 发布消息失败日志处理器
	publishConfirmFailLogger *zap.Logger // 确认发布消息失败日志处理器
	publishReturnFailLogger *zap.Logger // 发布消息入列失败日志处理器
	subscribeFailLogger *zap.Logger // 订阅失败日志处理器
)

// 初始化日志记录
func initLog(dir string)  {
	if dir == "" {
		dir = defaultLogDir
	}

	publishFailLogger = zap.New(zap.Config{
		LogDir: strings.TrimRight(dir, "/") + "/" + "publishFail",
	})

	publishConfirmFailLogger = zap.New(zap.Config{
		LogDir: strings.TrimRight(dir, "/") + "/" + "publishConfirmFail",
	})

	publishReturnFailLogger = zap.New(zap.Config{
		LogDir: strings.TrimRight(dir, "/") + "/" + "publishReturnFail",
	})

	subscribeFailLogger = zap.New(zap.Config{
		LogDir: strings.TrimRight(dir, "/") + "/" + "subscribeFail",
	})
}

func publishFailLog(data map[string]interface{})  {
	publishFailLogger.Error("publishFail", data)
}

func publishConfirmFailLog(data map[string]interface{})  {
	publishConfirmFailLogger.Error("publishConfirmFail", data)
}

func publishReturnFailLog(data map[string]interface{})  {
	publishReturnFailLogger.Error("publishReturnFail", data)
}

func subscribeFailLog(data map[string]interface{})  {
	subscribeFailLogger.Error("subscribeFail", data)
}

func subscribeIdpFailLog(data map[string]interface{})  {
	subscribeFailLogger.Error("subscribeIdp", data)
}

func subscribeConnErrorLog(data map[string]interface{})  {
	subscribeFailLogger.Error("subscribeConnErrorLog", data)
}