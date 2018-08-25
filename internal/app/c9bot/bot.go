package c9bot

import (
	"github.com/bwmarrin/discordgo"
)

//RunC9Bot sets up and runs the bot
func RunC9Bot() error {

	token := MustEnv("DISCORD_BOT_TOKEN")
	databaseAddress := MustEnv("DATABASE_ADDRESS")
	databaseUsername := MustEnv("DATABASE_USERNAME")
	databasePassword := MustEnv("DATABASE_PASSWORD")

	persistence, err := createmySQLPersistenceLayer(databaseAddress, databaseUsername, databasePassword)
	if err != nil {
		return err
	}
	defer persistence.Close()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	defer dg.Close()

	dg.AddHandler(makeMessageHandler("^!c9", persistence, recordC9))
	dg.AddHandler(makeMessageHandler("^!howlong", persistence, reportLastC9))

	err = dg.Open()
	if err != nil {
		return err
	}

	return nil
}
