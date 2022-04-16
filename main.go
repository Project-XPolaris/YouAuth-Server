package main

import (
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/cli"
	"github.com/sirupsen/logrus"
	"youauth/application/httpapi"
	"youauth/config"
	"youauth/database"
	"youauth/youlog"
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
