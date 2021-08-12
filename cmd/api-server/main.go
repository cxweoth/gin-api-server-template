/*
 api server main package
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cxweoth/gin-api-server-template/internal/api"
	"github.com/cxweoth/gin-api-server-template/internal/conf"
	"github.com/cxweoth/gin-api-server-template/internal/logger"
)

// This vars help to binding easily.
// [Note: interface and function can binding easily.]
var cfg conf.IConf
var makeLogger = logger.MakeLogger

func init() {
	cfg = &conf.Conf{}
}

func main() {
	os.Exit(RealMain(os.Stdout))
}

func RealMain(out io.Writer) int {

	// Init log output
	log.SetOutput(out)

	// Define cli inputs
	cliConfigPath := flag.String("cfgpath", "", "Input config path")
	flag.Parse()

	// Check input arguments
	if *cliConfigPath == "" {
		log.Printf("argument needs config file path with -cfgpath in cli!")
		return 1
	}
	// Read config file path from arguments
	configFilePath := *cliConfigPath

	// Check config file exists or not
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Printf("config file not exists!")
		return 1
	}

	// Load config file
	if err := cfg.Load(configFilePath); err != nil {
		fmt.Println(err)
		log.Printf("try load config file [%s] error [%s]\n", configFilePath, err.Error())
		return 1
	}

	// Init logger
	logger, err := makeLogger(cfg)
	if err != nil {
		log.Printf("try init logger failed: " + err.Error())
		return 1
	}
	log.Printf("logger already init")

	// Run API server
	logger.Info("Run API Server")
	if err = api.RunServer(cfg, logger); err != nil {
		log.Printf("try run API server failed: " + err.Error())
		return 1
	}

	return 0
}
