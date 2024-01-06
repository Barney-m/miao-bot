package main

import (
	"flag"
	"miao-bot/bot"
	"miao-bot/config"
	"miao-bot/logger"
)

var (
	debug = flag.Bool("debug", true, "Debug Mode")
)

func main() {
	flag.Parse()

	go logger.InitLogger(*debug)

	// Read Config
	config.ReadConfig()

	bot.Run()
}
