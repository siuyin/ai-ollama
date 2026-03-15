package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

type ToolParams struct {
	Type       string                      `json:"type"`
	Defs       any                         `json:"$defs,omitempty"`
	Items      any                         `json:"items,omitempty"`
	Required   []string                    `json:"required"`
	Properties map[string]api.ToolProperty `json:"properties"`
}

func main() {
	model := dflt.EnvString("MODEL", "qwen3.5:0.8b")
	host := dflt.EnvString("OLLAMA_HOST", "http://localhost:11434")
	thinkStr := dflt.EnvString("THINK", "false")
	log.Printf("MODEL=%s OLLAMA_HOST=%s THINK=%s", model, host, thinkStr)

	prompt := dflt.EnvString("PROMPT", "1. what is the weather in Bukit Batok, Singapore?\n 2. What is the UTC time?")
	log.Printf("PROMPT=%q", prompt)

	client := getClient()

	messages := []api.Message{
		{
			Role:    "system",
			Content: "Provide very brief, concise responses. Check all necessary tool calls have been made.",
		},
		{
			Role: "user",
			//Content: fmt.Sprintf("what is the UTC time?"),
			//Content: fmt.Sprintf("what is the weather in %s?", loc),
			Content: prompt,
		},
	}
	gwtProps := api.NewToolPropertiesMap()
	gwtProps.Set("location", api.ToolProperty{
		Type: []string{"string"},
	})
	getWeatherTool := api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "getWeather",
			Description: "Get the weather for a given location",
			Parameters: api.ToolFunctionParameters{
				Type:       "object",
				Required:   []string{"location"},
				Properties: gwtProps,
			},
		},
	}

	getTimeTool := api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        "getTime",
			Description: "Get the current time in UTC.",
			Parameters: api.ToolFunctionParameters{
				Type:     "object",
				Required: []string{},
			},
		},
	}

	think := false
	if thinkStr != "false" {
		think = true
	}

	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Tools:    []api.Tool{getWeatherTool, getTimeTool},
		Options:  map[string]any{"Temperature": 0.1},
		Think:    &api.ThinkValue{Value: think},
	}

	respFunc := func(resp api.ChatResponse) error {
		for _, tc := range resp.Message.ToolCalls {
			fn := tc.Function
			log.Printf("Model wants to call tool: %s with args %v", fn.Name, dump(fn.Arguments))
			switch fn.Name {
			case "getWeather":
				loc, ok := fn.Arguments.Get("location")
				if !ok {
					log.Fatal("error geting location")
				}
				output, err := getWeather(loc.(string))
				if err != nil {
					log.Fatalf("error executing tool: %v", err)
				}

				// Add the tool's output to the messages list as a new "tool" role message.
				messages = append(messages, api.Message{
					Role:    "tool",
					Content: output,
				})
			case "getTime":
				output := getTime()
				messages = append(messages, api.Message{
					Role:    "tool",
					Content: output,
				})
			default:
				log.Fatalf("invalid function: %q", fn.Name)
			}
		}

		fmt.Print(resp.Message.Thinking)
		fmt.Print(resp.Message.Content)

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
	res := fmt.Sprintf("It is currently 30°C in %s. Humidity is 80%%. Rain is expected later.\n", loc)
	log.Printf("\tTool: getWeather called with arg: %s. resp: %s", loc, res)
	return res, nil
}

func getClient() *api.Client {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal("getClient: ", err)
	}

	return client
}

func getTime() string {
	res := time.Now().UTC().Format("15:04:05 UTC")
	log.Printf("\tTool: getTime called. resp: %s", res)
	return res
}

func dump(args api.ToolCallFunctionArguments) []string {
	ret := []string{}
	for _, arg := range args.All() {
		ret = append(ret, arg.(string))
	}
	return ret
}
