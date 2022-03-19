package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO figure out if its possible to call the discord search function through the discord API
//Not fcking possible shit
//Takes all emote input and adds to the database with ++
func TrackEmote(Emote string, s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return
	}
	author := m.Author
	person := strings.ToLower(author.Username)
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//Needs to implement a search feature to check if emote already exists in the database
	userDatabase := client.Database(guild.ID + "EmoteDatabase")
	userCollection := userDatabase.Collection(person)

	//var emoteCount int
	filterCursor, err := userCollection.Find(ctx, bson.M{"Emote": Emote})
	if err != nil {
		log.Fatal(err)
	}
	var result []bson.M
	if err = filterCursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	if len(result) == 0 {
		emoteResult, err := userCollection.InsertOne(ctx, bson.D{
			{Key: "Emote", Value: Emote},
			{Key: "EmoteCount", Value: 1},
			{Key: "AuthorUserName", Value: person}})
		fmt.Println(emoteResult)
		if err != nil {
			log.Fatal(err)
		}
		client.Disconnect(ctx)
	}

	if len(result) >= 1 {
		fmt.Println("this emote already exists")

		emoteCount := result[0]["EmoteCount"].(int32)
		emoteCount = emoteCount + 1
		objectID := result[0]["_id"]
		result, err := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			bson.D{
				{"$set", bson.D{{"EmoteCount", emoteCount}}}})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
		client.Disconnect(ctx)
	}
}

//needs to accept a name and kind of emote
func ListEmote(Person string, s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return
	}
	author := m.Author
	author.Username = strings.ToLower(author.Username)

	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	userDatabase := client.Database(guild.ID + "EmoteDatabase")
	userCollection := userDatabase.Collection(Person)
	opts := options.Find()
	opts.SetSort(bson.D{{"EmoteCount", -1}})
	sortCursor, err := userCollection.Find(ctx, bson.D{{"EmoteCount", bson.D{{"$gt", 0}}}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = sortCursor.All(ctx, &ListEmoteResult); err != nil {
		log.Fatal(err)
	}
}

//Needs to accept emote
func LeaderBoard(Emote string, s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Leaderboard function called")
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return
	}
	author := m.Author
	author.Username = strings.ToLower(author.Username)

	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	userDatabase := client.Database(guild.ID + "EmoteDatabase")
	result, err := userDatabase.ListCollectionNames(ctx,
		bson.D{{ /* "options.capped", true */ }})
	if err != nil {
		log.Fatal(err)
	}
	for _, coll := range result {
		userCollection := userDatabase.Collection(coll)
		//	fmt.Println(userCollection)
		//	opts := options.Find()
		filterCursor, err := userCollection.Find(ctx, bson.M{"Emote": Emote})
		if err != nil {
			log.Fatal(err)
		}
		var result []bson.M
		if err = filterCursor.All(ctx, &result); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Result of filtercurser: ", result)

		for _, v := range result {
			m, m1, m2 := v["Emote"], v["EmoteCount"], v["AuthorUserName"]
			//result2 := result["EmoteCount"]
			//result3 := result["AuthorUserName"]
			userTempCollection := userDatabase.Collection("temp")
			userTempCollection.InsertOne(ctx, bson.D{
				{Key: "Emote", Value: m},
				{Key: "EmoteCount", Value: m1},
				{Key: "AuthorUserName", Value: m2}})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	userTempCollection := userDatabase.Collection("temp")
	opts := options.Find()
	opts.SetSort(bson.D{{"EmoteCount", -1}})
	sortCursor, err := userTempCollection.Find(ctx, bson.D{{"EmoteCount", bson.D{{"$gt", 0}}}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = sortCursor.All(ctx, &ListEmoteResult); err != nil {
		log.Fatal(err)

	}
	userTempCollection.DeleteMany(ctx, bson.M{"Emote": result})
	userDatabase.RunCommand(ctx, bson.M{"drop": "temp"})
	//userDatabase := client.Database(GuildID + "EmoteDatabase")
	//userTempCollection := userDatabase.Collection("temp")

}

//Might be able to just write all the data to a temporary file which can then be read and deleted afterworths

//Todo -- everything from under here is copy pasta from https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.8.2/mongo#Collection.Aggregate

/* var coll *mongo.Collection

	// Specify a pipeline that will return the number of times each name appears
	// in the collection.
	// Specify the MaxTime option to limit the amount of time the operation can
	// run on the server.
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$name"},
			{"numTimes", bson.D{
				{"$sum", 1},
			}},
		}},
	}
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := coll.Aggregate(
		context.TODO(),
		mongo.Pipeline{groupStage},
		opts)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Printf(
			"name %v appears %v times\n",
			result["_id"],
			result["numTimes"])
	}
} */

/*
userDatabase.aggregate([{
	$lookup: {
	 from: guild.ID,
	}}]);


	fmt.Println(userCollection)s
	/* opts := options.Find()
	opts.SetSort(bson.D{{"EmoteCount", -1}})
/* 	sortCursor, err := userCollection.Find(ctx, bson.D{{"EmoteCount", bson.D{{"$gt", 0}}}}, opts) */
/* if err != nil {
	log.Fatal(err)
} */
/*  if err = sortCursor.All(ctx, &ListEmoteResult); err != nil {
		log.Fatal(err)

}
*/
