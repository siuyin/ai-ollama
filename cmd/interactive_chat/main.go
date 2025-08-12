package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	client, err := api.ClientFromEnvironment() // eg. OLLAMA_HOST=http://imac2.h:11434
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("enter: {exit} to exit")
	systemPrompt := `You have access to functions. If you decide to invoke any of the function(s),
you MUST put it in the json format of
{"name": function name, "parameters": dictionary of argument name and its value}

You SHOULD NOT include any other text in the response if you call a function.
[
  {
    "name": "time",
    "description": "gets the current time for a given timezone",
    "parameters": {
      "type": "object",
      "properties": {
        "timezone": {
          "type": "string"
	  "default": "UTC"
        }
      }
    }
  },
  {"name":"Weather", "description":"get the weather forecast for a given city",
       "parameters": {"type":"object", "properties": {"city":{"type":"string}},"required": "city" }
  },
  {"name":"LocalEats", "description":"returns a list of recommended eats.", "parameters": null},
  {"name":"LocalAttractions", "description":"returns a list of recommended attractions.","parameters": null },
  {"name":"Parks", "description":"returns a list of nearby parks and gardens.", "parameters": null }
]

Hi
`

	messages := []api.Message{
		{Role: "user", Content: systemPrompt},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    dflt.EnvString("AI_MODEL", "gemma3:1b"),
		Messages: messages,
	}

	respFunc := func(resp api.ChatResponse) error {
		s := strings.Replace(resp.Message.Content, "\n", "", -1)
		fmt.Print(s)
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
