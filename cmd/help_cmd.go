package cmd

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	zLog "github.com/rs/zerolog/log"
)

// Get help from us
// Command: help
func ExecuteHelp(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Retrieving help message from markdown
	content, err := getHelpMd()
	if err != nil {
		log.Fatal(err)
		zLog.Error().Msg(err.Error())
		discord.ChannelMessageSend(message.ChannelID, "Encountered Error!")
		return
	}

	discord.ChannelMessageSend(message.ChannelID, string(content))
}

// Get help message from markdown file
func getHelpMd() (content []byte, err error) {
	// Retrieve current project root directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		zLog.Error().Msg(err.Error())
	}

	// Filepath: {PROJECT_ROOT_DIR}/help.md
	content, err = os.ReadFile(dir + "/help.md")
	if err != nil {
		return nil, err
	}
	return content, nil
}
