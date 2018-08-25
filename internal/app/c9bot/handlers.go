package c9bot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/when"
)

type c9botSession interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
}

func makeMessageHandler(messageRegexp string, p persistenceLayer, f func(c9botSession, *discordgo.MessageCreate, persistenceLayer)) func(*discordgo.Session, *discordgo.MessageCreate) {
	re := regexp.MustCompile(messageRegexp)
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if re.MatchString(m.Content) {
			f(s, m, p)
		}
	}
}

func recordC9(s c9botSession, m *discordgo.MessageCreate, p persistenceLayer) {
	msg := strings.Trim(m.Content[3:], " ")

	if len(msg) > 0 {
		r, err := when.EN.Parse(msg, time.Now())
		if err != nil {
			log.Println("Error parsing C9 time")
			r = nil
		}
		if r != nil {
			err := p.insertOccurence(m.ChannelID, r.Time)
			if err != nil {
				log.Printf("Could not insert C9 at %s on channel %s. Error: %s\n", r.Time, m.ChannelID, err)
				sendErrorMessage(s, m.ChannelID)
			}

			returnMsg := fmt.Sprintf("Logged a C9 %s ago", durationString(time.Since(r.Time)))
			s.ChannelMessageSend(m.ChannelID, returnMsg)
			return
		}
	}

	currentTime := time.Now()
	err := p.insertOccurence(m.ChannelID, currentTime)
	if err != nil {
		log.Printf("Could not insert C9 at %s on channel %s. Error: %s\n", currentTime, m.ChannelID, err)
		sendErrorMessage(s, m.ChannelID)
	}
	s.ChannelMessageSend(m.ChannelID, "Reset the C9 counter :(")
}

func reportLastC9(s c9botSession, m *discordgo.MessageCreate, p persistenceLayer) {
	mostRecent, err := p.getMostRecentOccurrence(m.ChannelID)
	switch {
	case err == errNoOccurrence:
		s.ChannelMessageSend(m.ChannelID, "No C9s logged.")
		return
	case err != nil:
		log.Printf("Error getting most recent C9 on channel %s. Error: %s", m.ChannelID, err)
		sendErrorMessage(s, m.ChannelID)
	default:
		d := time.Since(mostRecent)
		msg := fmt.Sprintf("It has been %s without a C9.", durationString(d))

		longestGap, err := p.getLongestGap(m.ChannelID)
		switch {
		case err == errNoOccurrence || longestGap < int(d.Hours())/24:
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

func durationString(d time.Duration) string {
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

func sendErrorMessage(s c9botSession, c string) {
	s.ChannelMessageSend(c, "An error occurred")
}

func parseChannelID(channelID string) (uint64, error) {
	return strconv.ParseUint(channelID, 10, 64)
}
