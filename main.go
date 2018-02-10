package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func MustEnv(env string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		log.Panicf("Environment variable $%v not set", env)
	}
	return v
}

func main() {

	Token := MustEnv("DISCORD_BOT_TOKEN")
	err := createDBFromEnv()
	if err != nil {
		log.Fatal("error creating database connection", err)
	}
	defer db.Close()

	err = initDB()
	if err != nil {
		log.Fatal("error initializing database connection", err)
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session", err)
	}
	defer dg.Close()

	dg.AddHandler(makeMessageHandler("!c9", recordC9))
	dg.AddHandler(makeMessageHandler("!howlong", reportLastC9))

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection", err)
	}
	log.Println("Bot started. Press CTRL-C to exit")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit

}
