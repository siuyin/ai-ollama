package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	model := dflt.EnvString("MODEL", "qwen3.5:0.8b")
	log.Printf("MODEL=%s", model)

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		api.Message{
			Role:    "system",
			Content: "Do not over-think, provide very brief, concise responses",
		},
		api.Message{
			Role:    "user",
			Content: "Name some unusual animals",
		},
		api.Message{
			Role:    "assistant",
			Content: "Monotreme, platypus, echidna",
		},
		api.Message{
			Role:    "user",
			Content: "which of these is the most dangerous?",
		},
	}

	think := false
	thinkEnv := dflt.EnvString("THINK", "false")
	if thinkEnv != "false" {
		think = true
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Think:    &api.ThinkValue{Value: think},
		Options:  map[string]any{"Temperature": 0.0, "Seed": 123},
	}

	respFunc := func(resp api.ChatResponse) error {
		if resp.Message.Thinking != "" {
			fmt.Print(resp.Message.Thinking)
		}

		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
