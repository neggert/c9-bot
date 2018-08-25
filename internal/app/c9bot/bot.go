package c9bot

import (
	"github.com/bwmarrin/discordgo"
)

// Bot contains all the bits of the bot that need to be cleand up on exit
type Bot struct {
	persistence persistenceLayer
	discord     *discordgo.Session
}

// Close cleans up the bot
func (b Bot) Close() {
	if b.persistence != nil {
		b.persistence.Close()
	}
	if b.discord != nil {
		b.discord.Close()
	}
}

//RunC9Bot sets up and runs the bot
func RunC9Bot() (Bot, error) {

	b := Bot{}

	token := MustEnv("DISCORD_BOT_TOKEN")
	databaseAddress := MustEnv("DATABASE_ADDRESS")
	databaseUsername := MustEnv("DATABASE_USERNAME")
	databasePassword := MustEnv("DATABASE_PASSWORD")

	var err error
	b.persistence, err = createmySQLPersistenceLayer(databaseAddress, databaseUsername, databasePassword)
	if err != nil {
		return b, err
	}

	b.discord, err = discordgo.New("Bot " + token)
	if err != nil {
		return b, err
	}

	b.discord.AddHandler(makeMessageHandler("^!c9", b.persistence, recordC9))
	b.discord.AddHandler(makeMessageHandler("^!howlong", b.persistence, reportLastC9))

	err = b.discord.Open()
	if err != nil {
		return b, err
	}

	return b, nil
}
