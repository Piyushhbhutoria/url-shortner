package main

import (
	"github.com/Piyushhbhutoria/url-shortner/logger"
	"github.com/Piyushhbhutoria/url-shortner/server"
)

func main() {
	logger.Init()
	server.Init()
}
