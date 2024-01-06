package bot

import (
	"fmt"
	"miao-bot/cmd"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Run() {
	discord, err := discordgo.New("Bot " + viper.GetString("discord.bot.token"))

	if err != nil {
		zlog.Error().Msg(err.Error())
	}

	discord.AddHandler(newMessage)
	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Escape once no prefix found
	command := strings.Split(message.Content, " ")[0]
	if command[:1] != viper.GetString("discord.bot.prefix") {
		return
	}

	switch command[1:] {
	case "help":
		cmd.ExecuteHelp(discord, message)
	case "chat":
		cmd.ExecuteChat(discord, message)
	case "music", "youtube", "yt":
		cmd.ExecuteMusic(discord, message)
	case "bye":
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
		// add more cases if required
	}
}
