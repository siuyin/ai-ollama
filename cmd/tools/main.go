package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

type prop struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}
type param struct {
	Type       string          `json:"type"`
	Required   []string        `json:"required"`
	Properties map[string]prop `json:"properties"`
}

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: "Provide very brief, concise responses",
		},
		{
			Role:    "user",
			Content: "what is the weather in Bukit Batok, Singapore like tomorrow?",
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "gemma3:4b",
		Messages: messages,
		Tools: []api.Tool{{Type: "function",
			Function: api.ToolFunction{
				Name:        "getWeather",
				Description: "Get the weather in a given location",
				Parameters: struct {
					Type       string   `json:"type"`
					Required   []string `json:"required"`
					Properties map[string]struct {
						Type        string   `json:"type"`
						Description string   `json:"description"`
						Enum        []string `json:"enum,omitempty"`
					} `json:"properties"`
				}{Type: "object",
					Properties: map[string]struct {
						Type        string   `json:"type"`
						Description string   `json:"description"`
						Enum        []string `json:"enum,omitempty"`
					}{"location": {"string", "the location to get the weather for", []string{}}},
				}}}},
		Options: map[string]any{
			"temperature": 0,
		},
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}
