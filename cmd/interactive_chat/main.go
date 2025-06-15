package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	client, err := api.ClientFromEnvironment() // eg. OLLAMA_HOST=http://imac2.h:11434
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("enter: {exit} to exit")

	messages := []api.Message{
		{
			Role: "system",
			Content: `You are a helpful assistant  with access to the tools listed below.
			Determine the user intent and decide if you can respond directly or
			if you need to call a tool.

			If you decide to call a tool(s), first list the tools to be called then
			compose a response with only a json array of tool calls.
			Eg. [getWeather("Singapore")]

			When responding directly you may use free form text.
			 
			AVAILABLE TOOLS:
			// getWeather gets the current weather
			func getWeather(location string)  string 

			// getTime returns the current time in UTC
			func getTime() time

			// getLocalEats returns a list of recommended eats.
			func getLocalEats() string

			// getLocalAttractions() returns a list of recommended local attractions.
			func getLocalAttractions() string

			// getLocalParks() return a list of recommended local parks.
			func getLocalParks() string

			`,
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    dflt.EnvString("AI_MODEL", "gemma3:1b"),
		Messages: messages,
	}

	respFunc := func(resp api.ChatResponse) error {
		s := strings.Replace(resp.Message.Content, "\n", "", -1)
		fmt.Print(s)
		messages = append(messages, resp.Message)
		return nil
	}

	for err := client.Chat(ctx, req, respFunc); err == nil; err = client.Chat(ctx, req, respFunc) {
		fmt.Printf("\n> ")
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		txt := sc.Text()
		if strings.Contains(txt, "{exit}") {
			break
		}
		msg := api.Message{Role: "user", Content: sc.Text()}
		messages = append(messages, msg)
		req.Messages = messages
	}

}
