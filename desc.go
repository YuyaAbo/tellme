package tellme

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/duckduckgo"
)

func Run(query string) error {
	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
	)
	if err != nil {
		return err
	}

	duckduckgoTool, err := duckduckgo.New(3, duckduckgo.DefaultUserAgent)
	if err != nil {
		return err
	}
	agentTools := []tools.Tool{
		duckduckgoTool,
	}

	executor, err := agents.Initialize(
		llm,
		agentTools,
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(3),
	)
	if err != nil {
		return err
	}

	input := "Describe " + query + "."
	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		return err
	}

	llmChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
		"Translate `{{.answer}}` into Japanese.",
		[]string{"answer"},
	))

	outputValues, err := chains.Call(context.Background(), llmChain, map[string]any{
		"answer": answer,
	})
	if err != nil {
		return err
	}

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		return fmt.Errorf("invalid chain return")
	}
	fmt.Println(out)

	return nil
}
