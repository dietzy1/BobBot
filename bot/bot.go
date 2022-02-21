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
		db.SearchData(function.Person)
		if db.Boolio == true {
			_, _ = s.ChannelMessageSend(m.ChannelID, function.Person+" is absolutely fcking pisslow elo")
			_, _ = s.ChannelMessageSend(m.ChannelID, db.SearchData(function.Person))
			//The problem is right here
		}
		if db.Boolio == false {
			_, _ = s.ChannelMessageSend(m.ChannelID, "For fcks sake can you pls input a prober name, or atleast add them to database")
		}
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
		if db.Boolio == false {
			fmt.Println(db.UrlString, db.Person)
			_, _ = s.ChannelMessageSend(m.ChannelID, "Not succesful! - Username already exists in database")
		}
		if db.Boolio == true {
			fmt.Println(db.UrlString, db.Person)
			_, _ = s.ChannelMessageSend(m.ChannelID, "User has succesfully been stored in the database, use '!elo user' to check")
		}
	}
	if strings.HasPrefix(m.Content, "!delete") {
		function.Delete(m)
		db.DeleteData(function.Person)
		if db.Boolio == true {
			_, _ = s.ChannelMessageSend(m.ChannelID, "User has been succesfully deleted from the database")
		}
		if db.Boolio == false {
			_, _ = s.ChannelMessageSend(m.ChannelID, "idk that retard aint in the database bro")
		}
	}
	if strings.HasPrefix(m.Content, "!help") {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Here is a list of commands and their shitty syntax: \n")
		_, _ = s.ChannelMessageSend(m.ChannelID, "!search, !elo, !delete, !add")
		_, _ = s.ChannelMessageSend(m.ChannelID, "!search region username     -Example: !search euw dietzy ")

		_, _ = s.ChannelMessageSend(m.ChannelID, "!elo username     -Example: !elo dietzy")
		_, _ = s.ChannelMessageSend(m.ChannelID, "!delete name     -Example: !delete dietzy")
		_, _ = s.ChannelMessageSend(m.ChannelID, "!add name Op.ggURL     -Example: !add dietzy https://euw.op.gg/summoner/userName=dietzy")
	}
	if strings.HasPrefix(m.Content, "!list") {
		db.List()
		for _, v := range db.ListResult {
			str := fmt.Sprintf("%v", v)
			_, _ = s.ChannelMessageSend(m.ChannelID, str)
		}
	}
}
