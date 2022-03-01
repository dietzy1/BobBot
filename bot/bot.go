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

/* type discordStruct struct {
	s discordgo.Session
	m discordgo.MessageCreate
} */

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
	//FIXED AND FUNCTIONAL
	if strings.HasPrefix(m.Content, "test") {
		C := make(chan string, 2)
		function.TestFunction(m, C)
		message := <-C
		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	}
	//FIXED AND FUNCTIONAL
	if strings.HasPrefix(m.Content, "!elo") {
		C := make(chan string, 2)
		function.SplitString(m)
		db.SearchData(function.Person, C)
		message := <-C
		_, _ = s.ChannelMessageSend(m.ChannelID, message)

	}
	//FIXED AND FUNCTIONAL
	if strings.HasPrefix(m.Content, "!search") {
		C := make(chan string, 2)
		function.SplitStringRegion(m)
		function.SplitStringPerson(m)
		//This function can be added
		function.Search(function.Region, function.Person, C)
		message := <-C
		_, _ = s.ChannelMessageSend(m.ChannelID, message)

	}
	//FIXED AND FUNCTIONAL
	if strings.HasPrefix(m.Content, "!add") {
		C := make(chan string, 2)
		function.Add(m)
		function.ValidateURL(function.Url, C)
		message := <-C
		if message == "Not a valid op.gg you fuckface" {
			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
		if message == "Valid URl" {
			db.StoreData(function.Person, function.Url, C)
			message := <-C
			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}
	// FIXED AND FUNCTIONAL
	if strings.HasPrefix(m.Content, "!delete") {
		C := make(chan string, 2)
		function.Delete(m)
		db.DeleteData(function.Person, C)
		message := <-C
		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	}
	//TODO
	if strings.HasPrefix(m.Content, "!help") {
		fmt.Println("Fcking learn how to do embeds")
	}
	if strings.HasPrefix(m.Content, "!list") {
		fmt.Println("Fcking learn how to use embeds and shit so bot doesn't crash like a retard")
	}
}

/* func rateLimit(u *discordgo.User, m *discordgo.MessageCreate) {

	mes
	discordgo.MessageActivity
	for i := 1; i < m.Author; i++ {

	}
	//m.Author.ID


} */

//Count each time m.Author.ID is seen if set amount of 50 is met before certain timer then their calls should be blocked.

/* func returnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, "Not a valid op.gg you fuckface")

} */

//this shit is absolutely useless
/* _, _ = s.ChannelMessageSend(m.ChannelID, "Here is a list of commands and their shitty syntax: \n")
_, _ = s.ChannelMessageSend(m.ChannelID, "!search, !elo, !delete, !add")
_, _ = s.ChannelMessageSend(m.ChannelID, "!search region username     -Example: !search euw dietzy ")

_, _ = s.ChannelMessageSend(m.ChannelID, "!elo username     -Example: !elo dietzy")
_, _ = s.ChannelMessageSend(m.ChannelID, "!delete name     -Example: !delete dietzy")
_, _ = s.ChannelMessageSend(m.ChannelID, "!add name Op.ggURL     -Example: !add dietzy https://euw.op.gg/summoner/userName=dietzy") */

/* if strings.HasPrefix(m.Content, "!list") {
db.List()
for _, v := range db.ListResult {
	str := fmt.Sprintf("%v", v)
	_, _ = s.ChannelMessageSend(m.ChannelID, str)
} */
