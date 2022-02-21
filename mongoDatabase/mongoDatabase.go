package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
var Boolio bool

//Used together with !add - Also checks prior if name already exists in database if it does its discarded.
func StoreData(Person, Url string) {
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	userDatabase := client.Database("userDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	User := UserStruct{
		Person: Person,
		Url:    Url,
	}
	filterCursor, err := userCollection.Find(ctx, bson.M{"Name": Person})
	var result []bson.M
	if err = filterCursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}
	if len(result) == 0 {
		insertResult, err := userCollection.InsertOne(ctx, User)
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult.InsertedID, "Cba dealing with this shit for now TODO")
		client.Disconnect(ctx)
		Boolio = true
	}
	if len(result) >= 1 {
		fmt.Println(User)
		var interfaceToString interface{}
		interfaceToString = result[0]["Name"]
		UrlString := fmt.Sprintf("%v", interfaceToString)
		if err != nil {
			panic(err)
		}
		if UrlString != Person {
			Boolio = true
			insertResult, err := userCollection.InsertOne(ctx, User)
			if err != nil {
				panic(err)
			}
			fmt.Println(insertResult.InsertedID, "Cba dealing with this shit for now TODO 🚀")
			client.Disconnect(ctx)
		}
		if UrlString == Person {
			Boolio = false
			client.Disconnect(ctx)
		}
	}
}

//Used together with !delete
func DeleteData(Person string) {
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	userDatabase := client.Database("userDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	result, err := userCollection.DeleteOne(ctx, bson.M{"Name": Person})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
	newBool := result.DeletedCount != 0
	Boolio = newBool
	fmt.Println(Boolio)
}

//Used together with !search
func SearchData(Person string) string {
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	userDatabase := client.Database("userDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	filterCursor, err := userCollection.Find(ctx, bson.M{"Name": Person})

	var result []bson.M
	if err = filterCursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}
	//DO NOTHING

	fmt.Println(len(result))
	if len(result) == 0 {
		Boolio = false
	}

	if len(result) == 1 {
		var interfaceToString interface{}
		interfaceToString = result[0]["Url"]
		UrlString := fmt.Sprintf("%v", interfaceToString)
		Boolio = true
		if err != nil {
			panic(err)
		}
		return UrlString
	}
	return UrlString

}

func List() {
	ReadConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Token))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	userDatabase := client.Database("userDatabase")
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
