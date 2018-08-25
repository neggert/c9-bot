package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/neggert/c9-bot/internal/app/c9bot"
)

func main() {

	bot, err := c9bot.RunC9Bot()
	defer bot.Close()
	if err != nil {
		log.Fatal("error creating bot", err)
	}
	log.Print("Successfully started C9 Bot")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit

}
