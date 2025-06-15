package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
	"github.com/yuin/goldmark"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model: "gemma3:1b",
		Prompt: `how many dwarf planets larger than Ceres are there?
		First gather your thoughts.
		Then get the diameter of Ceres.
		When cheking diameters ignore the "," which is used as a thousands separator.
		For example: 1,234 km should be read as 1234 km.
		Second example: 22,123 km should be read as 22123 km.
		Then check each answer against your knowlege by getting the diameter of each dwarf planet.
		Exclude any answers that you deem false.
		Include those answer that you have double checked to be true.
		Each candidate dwarf planet must be larger than Ceres in diameter.
		Count the number of answers that are true. Then include it in your response.
		Limit your response to a number and your reasoning.
		Output in json format.
		Example: {"Number":23, "Reason":"I've only considered the IAU definition of dwarf planet" }`,
		Options: map[string]any{"temperature": 0, "top_k": 1, "top_p": 0.2},
	}

	ctx := context.Background()
	output := ""
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.

		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Print(resp.Response)
		output += resp.Response
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	if err := goldmark.Convert([]byte(output), &b); err != nil {
		log.Fatal(err)
	}
	os.WriteFile("/tmp/junk.html", b.Bytes(), 0644)
}
