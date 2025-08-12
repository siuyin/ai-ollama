package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	model := dflt.EnvString("MODEL", "qwen3:0.6b")
	log.Printf("MODEL=%s", model)

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
			Content: "what is the weather in Bukit Batok, Singapore?",
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Tools: []api.Tool{
			{Type: "function",
				Function: api.ToolFunction{
					Name:        "getWeather",
					Description: "Get the weather in a given location",
					Parameters: struct {
						Type       string                      `json:"type"`
						Defs       any                         `json:"$defs,omitempty"`
						Items      any                         `json:"items,omitempty"`
						Required   []string                    `json:"required"`
						Properties map[string]api.ToolProperty `json:"properties"`
					}{
						Type:     "object",
						Required: []string{"location"},
						Properties: map[string]api.ToolProperty{
							"location": api.ToolProperty{Type: []string{"string"}},
						},
					},
				},
			},
		},
		Options: map[string]any{
			"temperature": 0,
		},
		Think: &api.ThinkValue{Value: false},
	}

	respFunc := func(resp api.ChatResponse) error {
		if len(resp.Message.ToolCalls) == 0 {
			fmt.Print(resp.Message.Content)
			return nil
		}
		tc := resp.Message.ToolCalls[0].Function
		log.Printf("Model wants to call tool: %s with args: %v", tc.Name, tc.Arguments)

		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}
