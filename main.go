package main

import (
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/cli"
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
