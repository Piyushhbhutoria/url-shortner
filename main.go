package main

import (
	"github.com/Piyushhbhutoria/go-gin-boilerplate/logger"
	"github.com/Piyushhbhutoria/go-gin-boilerplate/server"
	"github.com/Piyushhbhutoria/go-gin-boilerplate/store"
)

func main() {
	logger.Init()
	store.Init()
	server.Init()
}
