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

func main() {
	model := dflt.EnvString("MODEL", "qwen3.5:2b")
	host := dflt.EnvString("OLLAMA_HOST", "http://localhost:11434")
	svr := dflt.EnvString("SERVER", "myserver")
	prompt := dflt.EnvString("PROMPT", "What is the UTC time.")
	thinkStr := dflt.EnvString("THINK", "false")
	log.Printf("MODEL=%s OLLAMA_HOST=%s THINK=%s SERVER=%s PROMPT=%q", model, host, thinkStr, svr, prompt)

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

	tools, err := olamtl.FromMCP(lt.Tools)
	if err != nil {
		log.Fatal(err)
	}

	listOllam(tools)

	olamCl := getClient()

	messages := []api.Message{
		{
			Role: "system",
			Content: `You are a professional assistant with access to tools.
			If you do not know, say so.`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("%s", prompt),
		},
	}

	think := false
	if thinkStr != "false" {
		think = true
	}
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Tools:    tools,
		Options:  map[string]any{"Temperature": 0.1},
		Think:    &api.ThinkValue{Value: think},
	}

	respFunc := func(resp api.ChatResponse) error {
		for _, tc := range resp.Message.ToolCalls {
			fn := tc.Function
			log.Printf("Model wants to call tool: %s with args: %v", fn.Name, dumpFuncArgs(fn.Arguments))
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
		fmt.Print(resp.Message.Content)
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

func listOllam(tools []api.Tool) {
	for _, t := range tools {
		f := t.Function

		fmt.Printf("Name: %s, Description: %s\n\tParamsType: %v Properties: %v\n",
			f.Name, f.Description, f.Parameters.Type, dump(f.Parameters.Properties))
	}
	fmt.Println()
}

func dump(tpm *api.ToolPropertiesMap) []string {
	ret := []string{}
	for str, prop := range tpm.All() {
		ret = append(ret, fmt.Sprintf("%s: %s", str, prop.Type))
	}

	return ret
}

func dumpFuncArgs(args api.ToolCallFunctionArguments) []string {
	ret := []string{}
	for _, arg := range args.All() {
		ret = append(ret, arg.(string))
	}
	return ret
}
