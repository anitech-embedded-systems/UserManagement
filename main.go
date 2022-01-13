package main

import (
	dbm "main/Data"
	server "main/Server"
)

func main() {
	db := dbm.ConnectDB()
	if db == nil {
		return
	}
	server.HandleRequest()
	defer db.Close()
}
