package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/unsuman/discord-bot/config"
	"github.com/unsuman/discord-bot/ollamaorcalite"
)

var (
	BotID string
	goBot *discordgo.Session
)

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	goBot.AddHandler(messageCreate)

	err = goBot.Open()
	defer goBot.Close()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.Content[0] != ';' {
		return
	}

	if m.Content == ";ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!w")
	}

	if m.Content == ";pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if strings.HasPrefix(m.Content, ";orca") {
		prompt := strings.SplitN(m.Content, ";orca", 2)
		response := ollamaorcalite.ChatAI(prompt[1])
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
