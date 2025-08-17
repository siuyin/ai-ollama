package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	model := dflt.EnvString("MODEL", "qwen3:1.7b")
	host := dflt.EnvString("OLLAMA_HOST", "http://localhost:8080")
	prompt := dflt.EnvString("PROMPT", "Galpathi Golan (Ms) was born on Feb 28, 1962. She worked as a teacher at RGS for 18 years and then an education consultant for a further 21 years.")
	log.Printf("OLLAMA_HOST=%s MODEL=%s PROMPT=%s", model, host, prompt)

	client, err := api.ClientFromEnvironment() // eg. OLLAMA_HOST=http://imac2.h:11434
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		{
			Role: "system",
			Content: `Analyze the input below and use the following json schema for output:
			{name: string, name_without_salutation: string, sex: string (M|F), date_of_birth: string (dd/mm/yyyy), years_of_working_experience: int, industry: string, full_text: string}`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Input:
			%s`, prompt),
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
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
