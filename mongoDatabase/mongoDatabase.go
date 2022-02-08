package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Person string `bson:"Name"`
	Url    string `bson:"Url"`
}
type configStruct struct {
	Token string `json:"Token"`
}

var user User
var config *configStruct

// TODO needs to check if if name is already existing in the database
func StoreData(Person string, Url string) {
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
	user := User{
		Person: Person,
		Url:    Url,
	}
	insertResult, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)
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

/* func SearchData() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://MartinVad:xykjo6-fizbuh-Deqmim@cluster0.pdrpb.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	userDatabase := client.Database("userDatabase")
	userCollection := userDatabase.Collection("UserStructs")
	cursor, err := userCollection.Find(ctx, bson.M{"Name":}
	if err != nil {
		panic(err)
	}
	if err = cursor.All(ctx, &userCollection); err != nil {
		panic(err)
	}
/* 	fmt.Println(userCollection) */
/*
func DeleteData(Person string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://MartinVad:xykjo6-fizbuh-Deqmim@cluster0.pdrpb.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

}
*/
