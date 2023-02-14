package main

import (
	"context"
	"github.com/amperov/basic-auth-service/app/internal"
	"github.com/amperov/basic-auth-service/app/internal/redis"
	"github.com/amperov/basic-auth-service/app/internal/service"
	"github.com/amperov/basic-auth-service/app/internal/storage"
	"github.com/amperov/basic-auth-service/app/internal/transport"
	"github.com/amperov/basic-auth-service/app/internal/transport/grpc"
	"github.com/amperov/basic-auth-service/app/pkg/db"
	"github.com/amperov/basic-auth-service/app/pkg/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	grpc2 "google.golang.org/grpc"
	"net"
)

func main() {

	err := InitConfig()
	if err != nil {
		logrus.Fatalln("Reading config due error: ", err.Error())
		return
	}

	DBConfig, err := db.InitPGConfig()
	if err != nil {
		logrus.Fatalln("Reading DBConfig due error: ", err.Error())
		return
	}

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}

	PGClient, err := db.GetPGClient(context.Background(), DBConfig)
	if err != nil {
		return
	}

	TokenManager := tools.NewTokenManager()
	hasher := internal.NewHasher()

	AuthStorage := storage.NewAuthStorage(PGClient)

	RedisClient := redis.GetRedisClient()

	AuthService := service.NewAuthService(AuthStorage, RedisClient, TokenManager, hasher)

	AuthServer := transport.NewGRPCServer(AuthService)
	GRPCServer := grpc2.NewServer()
	grpc.RegisterAuthorizationServer(GRPCServer, AuthServer)

	err = GRPCServer.Serve(listen)
	if err != nil {
		logrus.Fatalln(err.Error())
		return
	}
}
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
