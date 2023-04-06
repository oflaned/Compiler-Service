package main

import (
	"Builder/Api/handler"
	"Builder/Api/service"
	"Builder/lib"
	"github.com/spf13/viper"
	"log"
)

const code = "#include <iostream>\n#include <cstdlib>\n#include <cmath>\nint main()\n{\n  using namespace std;\n  double a = 0, b = 0;\n  cin >> a;\n cin >> b;\n  cout.precision(16); \n  cout << \"a to b power  = \" << pow(a, b) << endl;\n  return 0;\n}"
const lang = "c++"
const stdin = "2\n3\n"

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
