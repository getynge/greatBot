package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func RateUser(good string, bad string) func(session *discordgo.Session, message *discordgo.MessageCreate) {
	return func(session *discordgo.Session, message *discordgo.MessageCreate) {
		target := message.Author

		fmt.Printf("Command from %s\n", message.Author)
		fmt.Printf("Command body %s\n", message.Content)
		if len(message.Mentions) > 0 && len(strings.Split(message.Content, " ")) > 1 {
			target = message.Mentions[0]
		}

		nickname := target.Username
		channel, err := session.State.Channel(message.ChannelID)
		if err != nil {
			fmt.Printf("Failed to get channel from channel id, using author username as nickname")
		} else {
			guild, err := session.State.Guild(channel.GuildID)
			if err != nil {
				fmt.Printf("Failed to get guild from guild id, using author username as nickname")
			} else {
				for _, m := range guild.Members {
					if m.User.ID == target.ID && m.Nick != "" {
						nickname = m.Nick
					}
				}
			}
		}
		if strings.Compare(target.Username, "getynge") != 0 {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s %s", nickname, bad))
		} else {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s %s", nickname, good))
		}
	}
}
