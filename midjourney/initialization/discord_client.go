package initialization

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
)

var discordClient *discord.Session

func LoadDiscordClient(create func(s *discord.Session, m *discord.MessageCreate), update func(s *discord.Session, m *discord.MessageUpdate)) {
	var err error
	discordClient, err = discord.New("Bot " + config.DISCORD_BOT_TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = discordClient.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// defer dg.Close()
	discordClient.AddHandler(create)
	discordClient.AddHandler(update)
}

func GetDiscordClient() *discord.Session {
	return discordClient
}
