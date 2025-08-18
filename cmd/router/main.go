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
	prompt := dflt.EnvString("PROMPT", "Write me a story about a magic backpack.")
	log.Printf("OLLAMA_HOST=%s MODEL=%s PROMPT=%s", model, host, prompt)

	client, err := api.ClientFromEnvironment() // eg. OLLAMA_HOST=http://imac2.h:11434
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: "You are an routing agent. Respond with ONLY one word: 'creative', 'weather', 'time', 'flight' or 'none' if no agent is applicable.",
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Based on the following user request, determine which agent or team or agents should handle it:
User request: %s
Available agents:
- 'weather' : Weather forecast agent that gets the current weather i.e. rain, wind, fog, temperature for various cities or locations.
- 'time': Select this agent if the user request includes the word 'time'. This agent reports the current time for various cities or locations, including UTC.
- 'flight': Flight booking agent, able to answer questions on available flights, fares, airport enquiries and book flights.
- 'creative': Creative writing agent that writes stories, poems or other creative text. Only select this agent if the user actually asks to write something.
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
