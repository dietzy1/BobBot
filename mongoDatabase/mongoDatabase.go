package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dietzy1/discord/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStruct struct {
	Person string `bson:"Name"`
	Url    string `bson:"Url"`
}
type configStruct struct {
	Token string `json:"Token"`
}

var User UserStruct
var config *configStruct
var UrlString string
var Person string

var NameString string
var ListResult []bson.M
var ListEmoteResult []bson.M

//Used together with !add - Also checks prior if name already exists in database if it does its discarded.
func StoreData(Person, Url, GuildID string, C chan string) {
	if Url == "" {
		message := "That op.gg is not valid you fucktard"
		C <- message
	}
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userDatabase := client.Database(GuildID + "OpggDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	User := UserStruct{
		Person: Person,
		Url:    Url,
	}
	filterCursor, err := userCollection.Find(ctx, bson.M{"Name": Person})
	if err != nil {
		log.Fatal(err)
	}
	var result []bson.M
	if err = filterCursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}
	//Means no result was found in the database
	if len(result) == 0 {
		insertResult, err := userCollection.InsertOne(ctx, User)
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult.InsertedID, "Cba dealing with this shit for now TODO")
		client.Disconnect(ctx)
		message := "Yo I added someone to the database"
		C <- message

	}
	//means a result was found in the database
	if len(result) >= 1 {
		fmt.Println(User)
		//var interfaceToString interface{}
		interfaceToString := result[0]["Name"]
		//interfaceToString = result[0]["Name"]
		UrlString := fmt.Sprintf("%v", interfaceToString)
		if err != nil {
			panic(err)
		}
		//Checks
		if UrlString != Person {
			insertResult, err := userCollection.InsertOne(ctx, User)
			if err != nil {
				panic(err)
			}
			fmt.Println(insertResult.InsertedID, "Cba dealing with this shit for now TODO 🚀")
			client.Disconnect(ctx)
			message := "Succesfull! but honestly no fcking clue why this part of the function is here"
			C <- message
		}
		//Username alreadt exists in data base - no action is taken
		if UrlString == Person {
			message := "Not succesfull! - Username already exists in database for fcks sake"
			C <- message
			client.Disconnect(ctx)
		}
	}
}

//Used together with !delete
func DeleteData(Person, GuildID string, C chan string) {
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	//TODO NEEDS TO ADD USERDATABASED BASED ON GUILD.ID
	userDatabase := client.Database(GuildID + "OpggDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	result, err := userCollection.DeleteOne(ctx, bson.M{"Name": Person})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
	if result.DeletedCount > 1 {
		message := "Yo I deleted" + function.Person
		C <- message
	}
	if result.DeletedCount < 1 {
		message := "Yo I didn't delete shit"
		C <- message
	}
}

//Used together with !search
func SearchData(Person, GuildID string, C chan string) {
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
	userDatabase := client.Database(GuildID + "OpggDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	filterCursor, err := userCollection.Find(ctx, bson.M{"Name": Person})
	if err != nil {
		log.Fatal(err)
	}
	var result []bson.M
	if err = filterCursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(result))
	if len(result) == 0 {
		message := "For fcks sake can you pls input a prober name, or atleast add them to database"
		C <- message
	}
	if len(result) == 1 {
		var interfaceToString interface{}
		interfaceToString = result[0]["Url"]
		UrlString := fmt.Sprintf("%v", interfaceToString)
		message := function.Person + " is absolutely pisslow " + UrlString
		C <- message
		if err != nil {
			panic(err)
		}
	}
}

func List(GuildID string) {
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
	userDatabase := client.Database(GuildID + "OpggDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	filterCursor, err := userCollection.Find(ctx, bson.M{})
	if err = filterCursor.All(ctx, &ListResult); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		panic(err)
	}
}
func AddFromSearch() {

}

//Used to hide the mongoDB APPLYURI
func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./mongoConfig.json")

	if err != nil {
		fmt.Println((err.Error()))
		return err
	}
	//fmt.Println(string(file))
	json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
	Token := config.Token
	_ = Token
	//Token = config.Token
	return nil
}
