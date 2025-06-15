package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

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
			Content: "Name some unusual animals",
		},
		{
			Role:    "assistant",
			Content: "Monotreme, platypus, echidna",
		},
		{
			Role:    "user",
			Content: "which of these is the most dangerous? Briefly explain in one paragraph.",
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    dflt.EnvString("AI_MODEL", "gemma3:1b"),
		Messages: messages,
		// Options: map[string]any{
		// 	"temperature": 0, "top_k": 5, "top_p": 1,
		// 	"stop": []string{"<end_of_turn>"},
		// },
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
