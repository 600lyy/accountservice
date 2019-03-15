package main

import (
	"fmt"

	"github.com/600lyy/go_study/accountservice/dbclient"
	"github.com/600lyy/go_study/accountservice/service"
)

var _appName = "accountservice"

func main() {
	fmt.Printf("Starting service %v\n", _appName)
	initializeBoltClient()
	service.StartWebServer("7676")

}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
