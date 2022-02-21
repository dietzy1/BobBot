package function

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

var CheckEloNoWhiteSpace string //Span scraped from op.gg - Contains rank and lp
var Url string                  //op.gg URL
var Person string               //Person assigned to a URL
var Region string               //Region assigned

func SplitString(m *discordgo.MessageCreate) string {
	str := strings.Split(m.Content, " ")[1:]
	if len(str) >= 3 {
		fmt.Println("Error")
		return "Error"
	}
	if len(str) == 0 {
		fmt.Println("error")
		return "Error"
	}
	if len(str) >= 1 {
		temp := ""
		for i := 0; i < len(str); i++ {
			temp = temp + str[i] + "%20"
			Person = Person + str[i] + "%20"
		}
		Person = strings.Trim(temp, "%20") //removes last %20

		return Person
	}
	return Person
}

//Acceps a message from API and removes prefix and indexes into 2 seperate substrings, who are assigned to variables.
//Important that it checks for errors, inputs with higher lenghts or nil length causes runtime errors.

func SplitStringSearch(m *discordgo.MessageCreate) (string, string, error) {
	str := strings.Split(m.Content, " ")[1:]
	if len(str) >= 3 {
		fmt.Println("Error")
		return "Error", "Error", errors.New("Error")
	}
	if len(str) == 0 {
		fmt.Println("Error")
		return "Error", "Error", errors.New("Error")
	}
	if len(str) >= 2 {
		Person = str[0]
		Region = str[1]
		return Person, Region, nil
	}
	if len(str) >= 1 {
		Person = str[0]
		Region = ""
		return Person, Region, nil
	}
	return Person, Region, nil
}

func SplitStringRegion(m *discordgo.MessageCreate) (string, error) {
	str := strings.Split(m.Content, " ")[1:2] //removes index 0 !search and only allows index 1
	if len(str) == 0 {
		fmt.Println("Error")
		return "Input was not readable see !help for further assistance.", errors.New("Error")
	}
	if len(str) >= 1 {
		Region = str[0]
		return Region, nil
	}
	if len(str) >= 2 {
		return "Error", errors.New("Error")
	}
	return Region, nil
}

//Splits from 2nd index to return the Person so it can accept whitespace account names
func SplitStringPerson(m *discordgo.MessageCreate) (string, error) {
	str := strings.Split(m.Content, " ")[2:]
	if len(str) == 0 {
		fmt.Println("Error")
		return "Input was not readable see !help for further assistance.", errors.New("Error")
	}
	if len(str) >= 1 {
		temp := ""
		for i := 0; i < len(str); i++ {
			temp = temp + str[i] + "%20"
			Person = Person + str[i] + "%20"
		}
		Person = strings.Trim(temp, "%20") //removes last %20
		return Person, nil
	}
	return Person, nil
}

//Search function that works together with splitstringperson and splitstringregion
func Search(Region string, Person string) string {
	if Region == "euw" || Region == "euwest" {
		Url := "https://euw.op.gg/summoner/" + "userName=" + Person //Euwest
		return Url
	}
	if Region == "eune" || Region == "northeast" {
		Url := "https://eune.op.gg/summoner/" + "userName=" + Person //EUNE
		return Url
	}
	if Region == "kr" || Region == "korea" {
		Url := "https://www.op.gg/summoner/" + "userName=" + Person //Korea
		return Url
	}
	if Region == "na" || Region == "northamerica" || Region == "murica" {
		Url := "https://na.op.gg/summoner/" + "userName=" + Person //NA
		return Url
	} else {
		return "Error"
	}
}

//Function used together with mongoDB,
func Add(m *discordgo.MessageCreate) (string, string, error) {
	str := strings.Split(m.Content, " ")[1:]
	if len(str) == 0 {
		fmt.Println("Error1")
		return "", "", errors.New("Error")
	}
	if len(str) == 1 {
		fmt.Println("Error2")
		fmt.Println(len(str))
		return "", "", errors.New("Error")

	}
	if len(str) >= 3 {
		fmt.Println("Error2")
		return "", "", errors.New("Error")
	}
	Person = str[0]
	Url = str[1]
	return Person, Url, nil
}

func Delete(m *discordgo.MessageCreate) (string, error) {
	str := strings.Split(m.Content, " ")[1:]
	if len(str) == 0 {
		fmt.Println("Error1")
		return "", errors.New("Error")
	}
	if len(str) == 1 {
		Person = str[0]
	}
	if len(str) >= 2 {
		fmt.Println("Error2")
		return "", errors.New("Error")
	}
	return Person, nil
}

//Function used to webscrabe spans and divs on a website manually.
func spanElo(Url string) string {
	if Url == "Error" {
		return "Error"
	}
	//Open and close the http handle
	response, err := http.Get(Url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	//Statuscode 100-200-300 good // 400-500 bad
	if response.StatusCode > 400 {
		fmt.Println("Status Code:", response.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	checkElo, err := doc.Find("span.LeaguePoints").Html()
	if err != nil {
		fmt.Println(err)
	}
	checkTierRank, err := doc.Find("div.TierRank").Html()
	if err != nil {
		fmt.Println(err)
	}
	//Removesconsequevent whitespace
	checkEloNoWhiteSpace := strings.TrimSpace(checkElo)
	checkTierRankNoWhiteSpace := strings.TrimSpace(checkTierRank)
	combinedTierElo := checkTierRankNoWhiteSpace + " " + checkEloNoWhiteSpace
	fmt.Println(combinedTierElo)
	return combinedTierElo
}
