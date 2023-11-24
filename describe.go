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

type AI struct {
	llm *openai.LLM
}

func NewAI() (*AI, error) {
	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
	)
	if err != nil {
		return nil, err
	}

	return &AI{
		llm: llm,
	}, nil
}

func (ai *AI) Describes(object string) (string, error) {
	answer, err := ai.generateAnswer(object)
	if err != nil {
		return "", err
	}

	jAnswer, err := ai.translateIntoJapanese(answer)
	if err != nil {
		return "", err
	}

	return jAnswer, nil
}

func (ai *AI) generateAnswer(object string) (string, error) {
	duckduckgoTool, err := duckduckgo.New(3, duckduckgo.DefaultUserAgent)
	if err != nil {
		return "", err
	}
	agentTools := []tools.Tool{
		duckduckgoTool,
	}

	executor, err := agents.Initialize(
		ai.llm,
		agentTools,
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(3),
	)
	if err != nil {
		return "", err
	}

	input := "Describe " + object + "."
	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func (ai *AI) translateIntoJapanese(text string) (string, error) {
	llmChain := chains.NewLLMChain(ai.llm, prompts.NewPromptTemplate(
		"Translate `{{.text}}` into Japanese.",
		[]string{"text"},
	))

	outputValues, err := chains.Call(context.Background(), llmChain, map[string]any{
		"text": text,
	})
	if err != nil {
		return "", err
	}

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		return "", fmt.Errorf("invalid chain return")
	}

	return out, nil
}
