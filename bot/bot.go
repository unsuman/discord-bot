package bot

import (
	"encoding/json"
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

type Res struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

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

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.Content[0] != ';' {
		return
	}

	if m.Content == ";ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == ";pong" {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Ping!")
	}

	if strings.HasPrefix(m.Content, ";orca") {
		prompt := strings.SplitN(m.Content, ";orca", 2)
		jsonResp := ollamaorcalite.ChatAI(prompt[1])
		generateMessage(s, m, jsonResp)
	}
}

func generateMessage(s *discordgo.Session, m *discordgo.MessageCreate, jsonResp []string) {
	var res []Res

	for _, item := range jsonResp {
		var r Res
		err := json.Unmarshal([]byte(item), &r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		res = append(res, r)
	}

	r := m.Author.Mention() + res[0].Response
	sum := 0
	msg, _ := s.ChannelMessageSend(m.ChannelID, r)

	for index, item := range res {
		if item.Done {
			s.ChannelMessageEdit(m.ChannelID, msg.ID, r)
			break
		}

		if index == 0 {
			continue
		}
		r += item.Response
		sum += len(item.Response)

		if sum >= 50 {
			s.ChannelMessageEdit(m.ChannelID, msg.ID, r)
			sum = 0
		}
	}

}
