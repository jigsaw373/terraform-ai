package cli

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	openai "github.com/PullRequestInc/go-gpt3"
	azureopenai "github.com/ia-ops/terraform-ai/pkg/gpt3"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
)

const userRole = "user"

var maxTokensMap = map[string]int{
	"code-davinci-002":   8001,
	"text-davinci-003":   4097,
	"gpt-3.5-turbo-0301": 4096,
	"gpt-3.5-turbo":      4096,
	"gpt-35-turbo-0301":  4096, // for azure
}

type oaiClients struct {
	azureClient  azureopenai.Client
	openAIClient openai.Client
}

func newOAIClients() (oaiClients, error) {
	var oaiClient openai.Client
	var azureClient azureopenai.Client
	var err error

	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		oaiClient = openai.NewClient(*openAIAPIKey)
	} else {
		re := regexp.MustCompile(`^[a-zA-Z0-9]+([_-]?[a-zA-Z0-9]+)*$`)
		if !re.MatchString(*openAIDeploymentName) {
			return oaiClients{}, errors.New("azure openai deployment can only include alphanumeric characters, '_,-', and can't end with '_' or '-'")
		}

		azureClient, err = azureopenai.NewClient(*azureOpenAIEndpoint, *openAIAPIKey, *openAIDeploymentName)
		if err != nil {
			return oaiClients{}, err
		}
	}

	clients := oaiClients{
		azureClient:  azureClient,
		openAIClient: oaiClient,
	}
	return clients, nil
}

func gptCompletion(ctx context.Context, client oaiClients, prompts []string, deploymentName string) (string, error) {
	temp := float32(*sensitivity)
	maxTokens, err := calculateMaxTokens(prompts, deploymentName)
	if err != nil {
		return "", err
	}

	var prompt strings.Builder
	fmt.Fprintf(&prompt, "You are a Terraform HCL generator, only generate valid Terraform HCL templates.")
	for _, p := range prompts {
		fmt.Fprintf(&prompt, "%s\n", p)
	}

	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		if isGptTurbo(deploymentName) {
			resp, err := client.openaiGptChatCompletion(ctx, prompt, maxTokens, temp)
			if err != nil {
				return "", err
			}
			return resp, nil
		}

		resp, err := client.openaiGptCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", err
		}
		return resp, nil
	}

	if isGptTurbo35(deploymentName) {
		resp, err := client.azureGptChatCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", err
		}
		return resp, nil
	}

	resp, err := client.azureGptCompletion(ctx, prompt, maxTokens, temp)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func isGptTurbo(deploymentName string) bool {
	return deploymentName == "gpt-3.5-turbo-0301" || deploymentName == "gpt-3.5-turbo"
}

func isGptTurbo35(deploymentName string) bool {
	return deploymentName == "gpt-35-turbo-0301" || deploymentName == "gpt-35-turbo"
}

func calculateMaxTokens(prompts []string, deploymentName string) (*int, error) {
	maxTokensFinal, ok := maxTokensMap[deploymentName]
	if !ok {
		return nil, fmt.Errorf("deploymentName %q not found in max tokens map", deploymentName)
	}

	if *maxTokens > 0 {
		maxTokensFinal = *maxTokens
	}

	encoder, err := gptEncoder.NewEncoder()
	if err != nil {
		return nil, err
	}

	// start at 100 since the encoder at times doesn't get it exactly correct
	totalTokens := 100
	for _, prompt := range prompts {
		tokens, err := encoder.Encode(prompt)
		if err != nil {
			return nil, err
		}
		totalTokens += len(tokens)
	}

	remainingTokens := maxTokensFinal - totalTokens
	return &remainingTokens, nil
}
