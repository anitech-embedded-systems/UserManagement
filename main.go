package main

import (
	config "main/Config"
	server "main/Server"
	mysqlclient "main/client/mysql"
	dataimpl "main/impl"
	"main/service"
)

func main() {

	cfg := config.Get()
	sqlclient := mysqlclient.New(cfg)

	defer sqlclient.Close()
	userdata1, _ := dataimpl.New(sqlclient)

	userservice1, _ := service.New(userdata1)

	serv1, _ := server.New(8080, *userservice1)

	serv1.Start()
	defer serv1.Stop()

	// wait for sigkill
	// close all resources
	// shutdown
}
