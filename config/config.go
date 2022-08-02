package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	aws_access_key_id     = "aws_access_key_id"
	aws_secret_access_key = "aws_secret_access_key"
	debug                 = "debug"
	aws_region            = "aws_region"
)

func Init() {
	viper.SetDefault("port", 8080)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/cloudfront-validator")
	viper.AddConfigPath("$HOME/.config/cloudfront-validator")
	viper.SetDefault("debug", false)
	viper.SetDefault(aws_region, "us-west-2")

	viper.AutomaticEnv()

	pflag.IntP("port", "P", 8080, "port to listen to")
	pflag.Bool("debug", false, "debug mode")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
	return
}

func IsDebug() bool {
	return viper.GetBool(debug)
}

func AwsKeyId() string {
	return viper.GetString(aws_access_key_id)
}

func AwsKeySec() string {
	return viper.GetString(aws_secret_access_key)
}

func AwsRegion() string {
	return viper.GetString(aws_region)
}
