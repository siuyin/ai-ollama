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
	host := dflt.EnvString("OLLAMA_HOST", "http://localhost:8080")
	prompt := dflt.EnvString("PROMPT", "I am planning a 3-day holiday to Kuala Lumpur, Malaysia with a budget of $1000 per person. Research places to visit. Food to try. Places to stay and walking and/or bicycling tours of the city.")
	log.Printf("OLLAMA_HOST=%s MODEL=%s PROMPT=%s", model, host, prompt)

	client, err := api.ClientFromEnvironment() // eg. OLLAMA_HOST=http://imac2.h:11434
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: "You are a planning agent. Carefully read the user request below. Break it down into action steps. Check again if all the action steps are covered. Also check if they are in the right order. When you have the plan in the correct sequence, present it in concise point form.",
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Based on the following user request, write out an action plan. 
User request: %s
`, prompt),
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Options:  map[string]any{"temperature": 0},
		Think:    &api.ThinkValue{Value: false},
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
