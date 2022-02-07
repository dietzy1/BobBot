package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dietzy1/discord/config"
	"github.com/dietzy1/discord/function"
	db "github.com/dietzy1/discord/mongoDatabase"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}
	BotID = u.ID
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	if strings.HasPrefix(m.Content, "!elo") {
		function.SplitString(m)
		_, _ = s.ChannelMessageSend(m.ChannelID, function.Person+" is "+function.CheckElo(function.Url))
	}
	if strings.HasPrefix(m.Content, "!search") {
		//	function.SplitStringSearch(m) //real function just testing the other
		function.SplitStringRegion(m)
		function.SplitStringPerson(m)
		_, _ = s.ChannelMessageSend(m.ChannelID, function.Search(function.Region, function.Person))
	}
	if strings.HasPrefix(m.Content, "!add") {
		function.Add(m)
		db.StoreData(function.Person, function.Url)
		_, _ = s.ChannelMessageSend(m.ChannelID, "User has succesfully been stored in the database, use '!elo user' to check")
	}
	if strings.HasPrefix(m.Content, "!delete") {
		function.Delete(m)
		db.DeleteData(function.Person)
		_, _ = s.ChannelMessageSend(m.ChannelID, "User has been succesfully deleted from the database")
	}

}
