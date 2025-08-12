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

	client := getClient()

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
		Options: map[string]any{"Temperature": 0.1},
		Think:   &api.ThinkValue{Value: false},
	}

	respFunc := func(resp api.ChatResponse) error {
		if len(resp.Message.ToolCalls) == 0 {
			fmt.Print(resp.Message.Content)
			return nil
		}

		tc := resp.Message.ToolCalls[0].Function
		log.Printf("Model wants to call tool: %s with args: %v", tc.Name, tc.Arguments)
		switch tc.Name {
		case "getWeather":
			loc := tc.Arguments["location"].(string)
			output, err := getWeather(loc)
			if err != nil {
				log.Fatalf("error executing tool: %v", err)
			}

			// Add the tool's output to the messages list as a new "tool" role message.
			messages = append(messages, api.Message{
				Role:    "tool",
				Content: output,
			})
		default:
			log.Fatalf("invalid function: %q", tc.Name)
		}

		return nil
	}

	ctx := context.Background()
	err := client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	req.Messages = messages
	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func getWeather(loc string) (string, error) {
	log.Printf("Tool: getWeather called with arg: %s", loc)
	return fmt.Sprintf("It is currently 30°C in %s. Humidity is 80%. Rain is expected later.\n", loc), nil
}

func getClient() *api.Client {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal("getClient: ", err)
	}

	return client
}
