package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
	"time"
)

func makeMessageHandler(messageRegexp string, f func(*discordgo.Session, *discordgo.MessageCreate)) func(*discordgo.Session, *discordgo.MessageCreate) {
	re := regexp.MustCompile(messageRegexp)
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if re.MatchString(m.Content) {
			f(s, m)
		}
	}
}

func recordC9(s *discordgo.Session, m *discordgo.MessageCreate) {
	currentTime := time.Now()
	err := insertOccurence(m.ChannelID, currentTime)
	if err != nil {
		log.Printf("Could not insert C9 at %s on channel %s. Error: %s\n", currentTime, m.ChannelID, err)
		sendErrorMessage(s, m.ChannelID)
	}
	s.ChannelMessageSend(m.ChannelID, "Reset the C9 counter :(")
}

func reportLastC9(s *discordgo.Session, m *discordgo.MessageCreate) {
	mostRecent, err := getMostRecentOccurrence(m.ChannelID)
	switch {
	case err == ErrNoOccurrence:
		s.ChannelMessageSend(m.ChannelID, "No C9s logged.")
		return
	case err != nil:
		log.Printf("Error getting most recent C9 on channel %s. Error: %s", m.ChannelID, err)
		sendErrorMessage(s, m.ChannelID)
	default:
		d := time.Since(mostRecent)
		msg := fmt.Sprintf("It has been %s without a C9.", DurationString(d))

		longestGap, err := getLongestGap(m.ChannelID)
		switch {
		case err == ErrNoOccurrence || longestGap < int(d.Hours())/24:
			msg += " A new record!"
		case err != nil:
			log.Printf("Error getting longest gap on channel %s. Error: %s", m.ChannelID, err)
		default:
			msg += fmt.Sprintf(" Current record is %d day", longestGap)
			if longestGap != 1 {
				msg += "s"
			}
			msg += "."
		}

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func DurationString(d time.Duration) string {
	var howMany int
	var units string
	switch {
	case d >= 24*time.Hour:
		hours := int(d.Hours())
		howMany = hours / 24
		units = "day"
	default:
		return "less than a day"
	}
	if howMany > 1 {
		units += "s"
	}
	return fmt.Sprintf("%d %s", howMany, units)
}

func sendErrorMessage(s *discordgo.Session, c string) {
	s.ChannelMessageSend(c, "An error occurred")
}

func parseChannelId(channelId string) (uint64, error) {
	return strconv.ParseUint(channelId, 10, 64)
}
