package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	common "go-ecommerce/common"

	"github.com/julienschmidt/httprouter"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	env := flag.String("env", "dev", "env=[dev, test, prod]")
	configFilePath := flag.String("app_config", "app.yaml", "the absolute path of APP's config file")
	flag.Parse()

	configBytes, err := os.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("can't read config file: %s", *configFilePath)
	}

	var configMap map[string]common.AppConfig
	err = yaml.Unmarshal(configBytes, &configMap)
	if err != nil {
		log.Fatalf("can't unmarshal config file: %s, msg=%s", *configFilePath, err.Error())
	}
	config, ok := configMap[*env]
	if !ok {
		log.Fatalf("config of env[%s] is missing in config file: %s", *env, *configFilePath)
	}
	for k := range configMap {
		delete(configMap, k)
	}

	config.Env = *env
	err = config.Init()
	if err != nil {
		log.Panic(err.Error())
	}
	common.Logger.WithField("config", config).Info()

	var port = 8080
	if config.ServerPort > 0 {
		port = config.ServerPort
	}
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome!\n")
	})
	common.Logger.Info("listening to " + strconv.Itoa(port) + ".......")
	common.Logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), router))

	// val, err := common.Redis.Get(context.Background(), "aid:1054034957").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// common.Logger.WithField("acccount", val).Info()

	// // var count int
	// // row := common.DB.
	// // row.Scan(&count)
	// var uid uid.Uid
	// common.DB.First(&uid)
	// common.Logger.WithField("uid", uid.String()).Info()
}
