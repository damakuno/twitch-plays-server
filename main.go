package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	twitchchat "github.com/dimorinny/twitch-chat-api"
)

type Creds struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	OauthToken   string `json:"oauth_token"`
}

func main() {
	var creds Creds
	// Open our credsJSON
	credsJSON, err1 := os.Open("./secrets/creds.json")
	// if we os.Open returns an error then handle it
	if err1 != nil {
		fmt.Println(err1)
	}
	// defer the closing of our credsJSON so that we can parse it later on
	defer credsJSON.Close()

	byteValue, _ := ioutil.ReadAll(credsJSON)
	json.Unmarshal(byteValue, &creds)

	fmt.Println(creds.ClientId)
	var config = twitchchat.NewConfiguration(
		"damakuno",
		creds.OauthToken,
		"damakuno",
	)
	twitch := twitchchat.NewChat(config)

	stop := make(chan struct{})
	defer close(stop)

	disconnected := make(chan struct{})
	connected := make(chan struct{})
	message := make(chan string)

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
				stop <- struct{}{}
			case <-connected:
				fmt.Println("Connected")
			case newMessage := <-message:
				fmt.Println(newMessage)
			}
		}
	}()

	if err := twitch.ConnectWithChannels(connected, disconnected, message); err != nil {
		return
	}

	<-stop
}
