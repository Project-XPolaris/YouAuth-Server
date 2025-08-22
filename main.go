package main

import (
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/cli"
	"github.com/allentom/harukap/plugins/nacos"
	httpapi "github.com/projectxpolaris/youauth/application/httpapi"
	"github.com/projectxpolaris/youauth/config"
	"github.com/projectxpolaris/youauth/database"
	"github.com/projectxpolaris/youauth/plugins/youlog"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.InitConfigProvider()
	if err != nil {
		logrus.Fatal(err)
	}
	err = youlog.DefaultYouLogPlugin.OnInit(config.DefaultConfigProvider)
	if err != nil {
		logrus.Fatal(err)
	}
	appEngine := harukap.NewHarukaAppEngine()
	appEngine.ConfigProvider = config.DefaultConfigProvider
	appEngine.LoggerPlugin = youlog.DefaultYouLogPlugin

	// 初始化并挂载 Nacos 插件（可选）
	defaultServiceName := config.DefaultConfigProvider.Manager.GetString("service.name")
	if defaultServiceName == "" {
		defaultServiceName = "youauth"
	}
	if nacosPlugin, err := nacos.NewNacosPluginFromYAML(config.DefaultConfigProvider, defaultServiceName, 8602); err != nil {
		logrus.Warnf("init nacos plugin failed: %v", err)
	} else {
		appEngine.UsePlugin(nacosPlugin)
	}
	appEngine.UsePlugin(database.DefaultPlugin)
	appEngine.HttpService = httpapi.GetEngine()
	if err != nil {
		logrus.Fatal(err)
	}
	appWrap, err := cli.NewWrapper(appEngine)
	if err != nil {
		logrus.Fatal(err)
	}
	appWrap.RunApp()
}
