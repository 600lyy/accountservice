package main

import (
	"fmt"

	"github.com/600lyy/accountservice/dbclient"
	"github.com/600lyy/accountservice/service"
)

var _appName = "accountservice"

func init() {
	initializeBoltClient()
	initializeRedisClient()
}

func main() {
	fmt.Printf("Starting service %v\n", _appName)
	service.StartWebServer("7676")
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}

func initializeRedisClient() {
	service.RedisClient = &dbclient.RedisClient{}
	service.RedisClient.OpenRedisDB()
}
