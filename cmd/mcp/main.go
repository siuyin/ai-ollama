package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
	"github.com/siuyin/mcptry/olamtl"
)

type ToolParams struct {
	Type       string                      `json:"type"`
	Defs       any                         `json:"$defs,omitempty"`
	Items      any                         `json:"items,omitempty"`
	Required   []string                    `json:"required"`
	Properties map[string]api.ToolProperty `json:"properties"`
}

func main() {
	model := dflt.EnvString("MODEL", "qwen3:0.6b")
	host := dflt.EnvString("OLLAMA_HOST", "http://localhost:11434")
	svr := dflt.EnvString("SERVER", "myserver")
	prompt := dflt.EnvString("PROMPT", "My name is Siu Yin.")
	log.Printf("MODEL=%s OLLAMA_HOST=%s SERVER=%s", model, host, svr)

	// Create a new mcpCl, with no features.
	mcpCl := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	ctx := context.Background()
	// Connect to a server over stdin/stdout
	transport := mcp.NewCommandTransport(exec.Command(svr))
	session, err := mcpCl.Connect(ctx, transport)
	if err != nil {
		log.Fatal("connect: ", err)
	}
	defer session.Close()

	// List Tools
	lt, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		log.Fatal("list tools: ", err)
	}

	tools, _ := olamtl.FromMCP(lt.Tools)

	olamCl := getClient()

	messages := []api.Message{
		{
			Role:    "system",
			Content: "You are a receptioninst agent. When meeting someone for the first time, say hi.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("%s", prompt),
		},
	}

	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Tools:    tools,
		Options:  map[string]any{"Temperature": 0.1},
		Think:    &api.ThinkValue{Value: false},
	}

	respFunc := func(resp api.ChatResponse) error {
		if len(resp.Message.ToolCalls) == 0 {
			fmt.Print(resp.Message.Content)
			return nil
		}

		for _, tc := range resp.Message.ToolCalls {
			fn := tc.Function
			log.Printf("Model wants to call tool: %s with args: %v", fn.Name, fn.Arguments)
			toolParam := &mcp.CallToolParams{
				Name:      fn.Name,
				Arguments: fn.Arguments,
			}
			output := mcpCallTool(session, toolParam)
			messages = append(messages, api.Message{
				Role:    "tool",
				Content: output,
			})
		}
		return nil
	}

	err = olamCl.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	req.Messages = messages
	err = olamCl.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func getClient() *api.Client {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal("getClient: ", err)
	}

	return client
}

func mcpCallTool(session *mcp.ClientSession, params *mcp.CallToolParams) string {
	ctx := context.Background()
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}
	s := ""
	for _, c := range res.Content {
		s += c.(*mcp.TextContent).Text
		//log.Print(c.(*mcp.TextContent).Text)
	}
	log.Printf("\tTool: %s called: output: %s", params.Name, s)
	return s
}
