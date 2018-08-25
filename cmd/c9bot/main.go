package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/neggert/c9-bot/internal/app/c9bot"
)

func main() {

	err := c9bot.RunC9Bot()
	if err != nil {
		log.Fatal("error creating bot", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit

}
