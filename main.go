package main

import (
	"fmt"
	"log"

	"github.com/songchenwen/cloudfront-invalidator/cf"
	"github.com/songchenwen/cloudfront-invalidator/config"
	"github.com/songchenwen/cloudfront-invalidator/server"
	"github.com/spf13/viper"
)

func init() {
	config.Init()
	err := cf.Init()
	if err != nil {
		panic(fmt.Sprintf("Cloudfront init err %v", err))
	}
}

func main() {
	engine := server.New()
	log.Printf("server starting at port %d\n", viper.GetInt("port"))
	engine.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
