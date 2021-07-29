package main

import (
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	myapi "github.com/proggcreator/WbWorkDb"
	"github.com/proggcreator/WbWorkDb/config"
	"github.com/proggcreator/WbWorkDb/handler"
	"github.com/proggcreator/WbWorkDb/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	if err := godotenv.Load(); err != nil { //переменные из .evm
		fmt.Println("No .env file found")
	}
}

func initConfig() error { //переменные из config
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) //формат для логгера json
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs")
	}

	conf := config.New()
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: conf.Username,
		Password: conf.Password,
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("faled to initialization %s", err.Error())
	}
	//иерархия зависимостей
	repos := repository.NewRepository(db)
	handlers := handler.NewHandler(repos)
	defer db.Close()

	srv := new(myapi.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}
