/*
 * @Notice: edit notice here
 * @Author: zhulei
 * @Date: 2022-09-22 11:00:10
 * @LastEditors: zhulei
 * @LastEditTime: 2022-09-22 15:42:38
 */
package model

import (
	// "fmt"
	// "net/http"

	"fmt"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// sentry 初始化
func InitSentry(app *gin.Engine) {

	err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("sentry.url"),
		// 发送之前的回调
		// BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
		// 	if hint.Context != nil {
		// 		if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
		// 			// You have access to the original Request here
		// 			// logrus.Infof(req)
		// 		}
		// 	}
		// 	return event
		// },
		AttachStacktrace: true,                                          // 详细的跟踪信息
		ServerName:       viper.GetString("sentry.server_name"),         // 要上报的服务名称
		SampleRate:       viper.GetFloat64("sentry.sample_rate"),        // 事件提交的采样率， (0.0 - 1.0, defaults to 1.0)
		TracesSampleRate: viper.GetFloat64("sentry.traces_sample_rate"), // 性能监控事件采样率 1.0 --> 100%， 生产根据性能调整， (defaults to 0.0 (meaning performance monitoring disabled))
		Environment:      viper.GetString("sentry.environment"),         // 设置环境，测试test、开发dev、生产prod
		Debug:            viper.GetBool("sentry.debug"),                 // 启用打印sdk debug消息
	})

	if err != nil {
		logrus.Errorf("sentry init err:%v", err)
		return
	}

	app.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	//
	// app.Use(func(ctx *gin.Context) {
	// 	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
	// 		fmt.Println("222222222222222")
	// 		hub.Scope().SetTag("someRandomTag", "这是windows")
	// 		hub.CaptureMessage("函数异常退出")
	// 	}
	// 	ctx.Next()
	// })
}

// 使用 model.SendSentryMsg(ctx, "这里报错了") 调用
// SendSentryMsg 发送sentry消息
func SendSentryMsg(ctx *gin.Context, e interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("cover LogSentry error %v", err)
		}
	}()
	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			// scope.SetExtra("unwantedQuery", "someQUeryDataMaybe") //添加额外业务自定义信息， scope.SetExtra("code", statusCode) -- scope.SetExtra("resp", blw.body.String())
			// hub.CaptureMessage(fmt.Sprintf("%s", e))
			// hub.CaptureMessage(e.(string))
			hub.CaptureMessage(fmt.Sprintf("%s", e)) // 发送的报错信息
		})
	}
}

// SendSentryMsgExtra 发送sentry消息带额外消息
func SendSentryMsgExtra(ctx *gin.Context, e string, extKey string, extValue string) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("cover LogSentry error %v", err)
		}
	}()
	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra(extKey, extValue) //添加额外业务自定义信息， scope.SetExtra("code", statusCode) -- scope.SetExtra("resp", blw.body.String())
			// hub.CaptureMessage(fmt.Sprintf("%s", e))
			// hub.CaptureMessage(e.(string))
			hub.CaptureMessage(e) // 发送的报错信息
		})
	}
}
