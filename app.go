package main

import (
	"Builder/Api/handler"
	"Builder/Api/service"
	"Builder/lib"
	"github.com/spf13/viper"
	"log"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal("error while initialization config")
	}

	dep, err := lib.CheckDependencies()
	if err != nil {
		log.Fatalf("dependency %s not found: %v", dep, err)
	}

	services := service.NewService()
	handlers := handler.NewHandler(services)

	server := new(Server)
	err = server.Run(viper.GetString("port"), handlers.InitRoutes())
	if err != nil {
		log.Fatal("Server are not run")
	}
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
