package main

import (
	"encoding/json"

	"log"
	"os"

	ginGonic "github.com/gin-gonic/gin"
	cors "github.com/krakendio/krakend-cors/v2/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/router/gin"
	clientPlugin "github.com/luraproject/lura/v2/transport/http/client/plugin"
	luraHTTPServer "github.com/luraproject/lura/v2/transport/http/server"
	serverPlugin "github.com/luraproject/lura/v2/transport/http/server/plugin"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/common/telemetry"
	_ "go.uber.org/automaxprocs"
)

var (
	traceFilename = "service_area_traces.txt"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile("./config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("err", err.Error(), "message", "Can't find config file")
	}
}

func main() { //nolint:funlen

	//	Open Telemetry
	l := log.New(os.Stdout, "", 0)
	f, err := os.Create(traceFilename)
	if err != nil {
		helper.CheckError(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	tel := telemetry.NewTelemetry(f, l)
	tel.InitProvider()
	tel.StartTracerProvider()
	defer tel.Shutdown()

	serviceConfig := loadConfigurations(GeneratedLuraConf)
	if err := serviceConfig.Init(); err != nil {
		panic(err)
	}
	serviceConfig.Normalize()
	logger, _ := logging.NewLogger("INFO", os.Stdout, "asd")
	loadPlugins(serviceConfig.Plugin.Folder, serviceConfig.Plugin.Pattern, logger)

	routerFactory := gin.NewFactory(gin.Config{
		Engine:      gin.NewEngine(serviceConfig, gin.EngineOptions{}),
		Middlewares: []ginGonic.HandlerFunc{},
		RunServer:   gin.RunServerFunc(cors.NewRunServer(luraHTTPServer.RunServer)),
	})

	x := routerFactory.New()
	x.Run(serviceConfig)
}

func loadConfigurations(raw []byte) config.ServiceConfig {
	var serviceConfig config.ServiceConfig
	var b map[string]interface{}
	err := json.Unmarshal(raw, &b)
	if err != nil {
		panic(err)
	}
	err = mapstructure.Decode(b, &serviceConfig)
	if err != nil {
		panic(err)
	}
	return serviceConfig
}

// loadPlugins loads and registers the plugins so they can be used if enabled at the configuration
func loadPlugins(directoryName, filePattern string, logger logging.Logger) {
	directoryName = directoryName
	n, err := clientPlugin.LoadWithLogger(
		directoryName,
		filePattern,
		clientPlugin.RegisterClient,
		logger,
	)
	if err != nil {
		logger.Warning(
			"plugin_http-client", "loading",
			"err", err.Error(),
			"message", err.Error(),
		)
	}

	logger.Info(
		"plugin_http-client", "done",
		"loaded_plugins", n,
		"message", "plugin/http-client are loaded",
	)

	nServer, errServer := serverPlugin.LoadWithLogger(directoryName, filePattern, serverPlugin.RegisterHandler, logger)
	if errServer != nil {
		logger.Warning(
			"plugin_http-server", "loading",
			"err", errServer.Error(),
			"message", errServer.Error(),
		)
	}

	logger.Info(
		"plugin_http-server", "done",
		"loaded_plugins", nServer,
		"message", "plugin/http-server are loaded",
	)
}
