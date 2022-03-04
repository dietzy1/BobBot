package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dietzy1/discord/config"
	"github.com/dietzy1/discord/embedHelp"
	"github.com/dietzy1/discord/function"
	db "github.com/dietzy1/discord/mongoDatabase"
)

var BotID string
var goBot *discordgo.Session

/* type GobotStruct struct {
	s *discordgo.Session
} */

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

//This is experimental shit
/* func (bob *GobotStruct) HandleMessageCreate() interface{} {

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
	}
} */

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
		embed := embedHelp.NewEmbed()
		embed.SetTitle("BobBot guide 2k22 visualized🚀🚀🚀🚀")
		embed.AddField("!search region username", "Example: !search euw twtvkibbylol")
		embed.AddField("!add name op.ggURL", "Example: !add kibby https://euw.op.gg/summoner/userName=twtvkibbylol ")
		embed.AddField("!elo assignedName", "Example: !elo kibby")
		embed.AddField("!delete assignedName", "!delete kibby")
		embed.AddField("🦖", "🦖")

		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
	}
	if strings.HasPrefix(m.Content, "!list") {
		embed := embedHelp.NewEmbed()
		embed.SetTitle("Pls dont fcking crash🚀🚀🚀🚀")
		db.List()
		for _, v := range db.ListResult {
			str := fmt.Sprintf("%v", v)
			strsplit := strings.Split(str, "[")
			embed.AddFieldnoValue(strsplit[1])
		}
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
	}
	if strings.HasSuffix(m.Content, "bonk") {
		embed := embedHelp.NewEmbed()
		embed.SetTitle("Fuck off weeb")
		embed.SetImage("https://c.tenor.com/yHX61qy92nkAAAAC/yoshi-mario.gif")
		embed.SetColor(0x00ff00)
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)

	}

	/* if strings.HasPrefix(m.Content, "bonk") {
		embed := embedHelp.NewEmbed()
		embed.SetTitle("Fuck off weeb")
		embed.SetImage("https://c.tenor.com/yHX61qy92nkAAAAC/yoshi-mario.gif")
		embed.SetColor(0x00ff00)
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
	} */
}

/* func (&embed.Embed) embedShit() {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: "This is a discordgo embed",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "I am a field",
				Value:  "I am a value",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "I am a second field",
				Value:  "I am a value",
				Inline: true,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     "I am an Embed",
	}
	fmt.Println(embed)
} */

/* session.ChannelMessageSendEmbed(channelid, embed)
}
*/

//I dont think so
//Its just some data structure homding user id, command and some sort of cooldown logic
//Could be a date or a remaining numbers in ms

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

/* if strings.HasPrefix(m.Content, "!list") {
db.List()
for _, v := range db.ListResult {
	str := fmt.Sprintf("%v", v)
	_, _ = s.ChannelMessageSend(m.ChannelID, str)
} */
