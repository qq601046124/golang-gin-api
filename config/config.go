package config

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/orandin/sentrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

// 读取配置
func (c *Config) InitConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	// 从环境变量中读取
	viper.AutomaticEnv()
	viper.SetEnvPrefix("web")
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	return viper.ReadInConfig()
}

// 监控配置改动
func (c *Config) WatchConfig(change chan int) {
	viper.WatchConfig()
	// TODO: 这个会触发两次, 考虑使用限流模式, 第一次是无效的
	// https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("配置已经被改变: %s", e.Name)

		// 非常有可能读到空的
		if err := viper.ReadInConfig(); err != nil || viper.GetString("db.addr") == "" {
			if err == nil {
				logrus.Warnf("配置更新后读取失败: 未读到数据")
			} else {
				logrus.Warnf("配置更新后读取失败: %s", err)
			}

			return
		}
		change <- 1
	})
}

// 初始化日志
func (c *Config) InitLog() {
	// log.logrus_json
	if viper.GetBool("log.logrus_json") {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// log.logrus_level
	switch viper.GetString("log.logrus_level") {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	// 可以设置日志单独打印到elesticsearch
	//esurl := viper.GetString("eslog.esurl")
	//esusername := viper.GetString("eslog.username")
	//espassword := viper.GetString("eslog.password")
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//httpsclient := &http.Client{Transport: tr} // 自定义transport
	//// 创建elasticsearch客户端
	//client, err := elastic.NewClient(elastic.SetHttpClient(httpsclient), elastic.SetSniff(false), elastic.SetURL(esurl), elastic.SetBasicAuth(esusername, espassword))
	//if err != nil {
	//	//logrus.Panic(err)
	//	logrus.Info(err)
	//}
	//hostname, _ := os.Hostname()
	//
	//// 将logrus和elastic绑定，localhost 是指定该程序执行时的ip
	//hook, err := elogrus.NewAsyncElasticHookWithFunc(client, hostname, logrus.DebugLevel, IndexName)
	//if err != nil {
	//	logrus.Info(err)
	//}
	//logrus.AddHook(hook)
	logrus.AddHook(sentrus.NewHook(
		[]logrus.Level{logrus.ErrorLevel},

		// Optional: add tags to add in each Sentry event
		sentrus.WithTags(map[string]string{"foo": "bar"}),

		// Optional: set custom CaptureLog function
		sentrus.WithCustomCaptureLog(sentrus.DefaultCaptureLog),
	))

	// log.logrus_file
	logrusFile := viper.GetString("log.logrus_file")
	os.MkdirAll(filepath.Dir(logrusFile), os.ModePerm)

	file, err := os.OpenFile(logrusFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		if viper.GetBool("log.logrus_console") {
			logrus.SetOutput(io.MultiWriter(file, os.Stdout))
		} else {
			logrus.SetOutput(file)
		}
	}

	// log.gin_file & log.gin_console
	ginFile := viper.GetString("log.gin_file")
	os.MkdirAll(filepath.Dir(ginFile), os.ModePerm)

	file, err = os.OpenFile(ginFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		if viper.GetBool("log.gin_console") {
			gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
		} else {
			gin.DefaultWriter = io.MultiWriter(file)
		}
	}

	// default
	logrus.SetReportCaller(true)
}
