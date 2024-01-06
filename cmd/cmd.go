package cmd

import (
	"miao-bot/constants"
	"miao-bot/services"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	// To store guilds
	guilds      = make(map[string]*services.ActiveGuild)
	guildsMutex = sync.RWMutex{}

	// guildNames is a mapping between guild id and guild name
	guildNames = make(map[string]string)
)

// Get Author Voice Channel
func GetVoiceChannelWhereMessageAuthorIs(bot *discordgo.Session, message *discordgo.MessageCreate) (string, error) {
	guild, err := bot.State.Guild(message.GuildID)
	if err != nil {
		return "", err
	}
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == message.Author.ID {
			return voiceState.ChannelID, nil
		}
	}
	return "", constants.ErrUserNotInVoiceChannel
}

// Retrieve Guild Name By Guild ID
func GetGuildNameByID(bot *discordgo.Session, guildID string) string {
	guildName, ok := guildNames[guildID]
	if !ok {
		guild, err := bot.Guild(guildID)
		if err != nil {
			// Failed to get the guild? Whatever, we'll just use the guild id
			guildNames[guildID] = guildID
			return guildID
		}
		guildNames[guildID] = guild.Name
		return guild.Name
	}
	return guildName
}
