package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("enter: {exit} to exit")

	messages := []api.Message{
		{
			Role:    "system",
			Content: `You are a helpful question answering assistant.
			However your a prone to making up answers when you do not have access to facts
			or real-time data. When in these situations, say you do not know.
			Instead prompt the user to use one or more of the functions below.
			Example:
			When asked for the weather, respond with "Please run getWeather(location: string) to get the weather forecast for a specified location."
			Example 2:
			When asked for the time, respond with "Please run getTime() to get the current local time."
			 
			USER FUNCTIONS:
			func getWeather(location: string) : string --> this gets the current weather forecast for a location.
			func getTime(): time --> this gets the current local time.
			func getLocalEats(): string --> this returns a list of recommended local eats.
			func getLocalAttractions(): string --> this returns a list of recommended local attractions.
			func getLocalParks(): string --> this returns a list of recommended local parks.
			`,
		},
		{
			Role:    "user",
			Content: "Hi.",
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "gemma3:1b",
		Messages: messages,
		Options: map[string]any{
			"temperature": 0, "top_k": 4, "top_p": 0.9,
			// "stop": []string{"<end_of_turn>"},
		},
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
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
