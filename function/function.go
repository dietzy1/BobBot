package function

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

var Url string    //op.gg URL
var Person string //Person assigned to a URL
var Region string //Region assigned

/* type discordStruct struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
}

var MessageHandler *discordStruct */

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
		lowercaseTemp := strings.ToLower(temp)
		Person = strings.Trim(lowercaseTemp, "%20") //removes last %20
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
		lowercaseTemp := strings.ToLower(temp)
		Person = strings.Trim(lowercaseTemp, "%20") //removes last %20
		return Person, nil
	}
	return Person, nil
}

//Search function that works together with splitstringperson and splitstringregion
func Search(Region string, Person string, C chan string) {
	if Region == "euw" || Region == "euwest" {
		Url := "https://euw.op.gg/summoner/" + "userName=" + Person //Euwest
		C <- Url
	}
	if Region == "eune" || Region == "northeast" {
		Url := "https://eune.op.gg/summoner/" + "userName=" + Person //EUNE
		C <- Url
	}
	if Region == "kr" || Region == "korea" {
		Url := "https://www.op.gg/summoner/" + "userName=" + Person //Korea
		C <- Url
	}
	if Region == "na" || Region == "northamerica" || Region == "murica" {
		Url := "https://na.op.gg/summoner/" + "userName=" + Person //NA
		C <- Url
	} else {
		message := "Ran into some error while searching shit"
		C <- message
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
	temp := str[0]
	Person = strings.ToLower(temp)
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

//https?:\/\/(euw)?(na)?(kr)?(eune)?(tr)?(las)?(lan)?\.op\.gg\/.+
//https?://(euw)?(na)?(kr)?(eune)?(tr)?(las)?(lan)?.op.gg/.+

func ValidateURL(Url string, C chan string) {
	validate, err := regexp.MatchString(`https?://(euw)?(na)?(kr)?(eune)?(tr)?(las)?(lan)?.op.gg/.+`, Url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(validate)
	if validate {
		response, err := http.Get(Url)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()
		if response.StatusCode > 400 {
			fmt.Println("Status Code:", response.StatusCode)
		}
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		validatedURL, err := doc.Find("h2.header__title").Html()
		if err != nil {
			fmt.Println(err)
		}
		validatedURLNoWhiteSpace := strings.TrimSpace(validatedURL)
		if validatedURLNoWhiteSpace == "This summoner is not registered at OP.GG. Please check spelling." {
			message := "Not a valid op.gg you fuckface"
			C <- message
		} else {
			message := "Valid URl"
			C <- message
		}
	}
	if !validate {
		message := "Not a valid op.gg you fuckface"
		C <- message
	}
}

func TestFunction(m *discordgo.MessageCreate, C chan string) {
	// TODO I need to format all the messages into the correct format pre channel liftoff
	message := "shut the fuck up kibby"
	C <- message

}
