package main

import (
	"fmt"

	"github.com/dietzy1/discord/bot"
	"github.com/dietzy1/discord/config"
)

func main() {

	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bot.Start()

	<-make(chan struct{})
	return
}
