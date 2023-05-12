package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	openai "github.com/PullRequestInc/go-gpt3"
	azureopenai "github.com/hubs-ai/terraform-ai/pkg/gpt3"
	"github.com/pkg/errors"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
)

const userRole = "user"

var (
	maxTokensMap = map[string]int{
		"code-davinci-002":   8001,
		"text-davinci-003":   4097,
		"gpt-3.5-turbo-0301": 4096,
		"gpt-3.5-turbo":      4096,
		"gpt-35-turbo-0301":  4096, // for azure
		"gpt-4-0314":         8192,
		"gpt-4-32k-0314":     8192,
	}
	errToken = errors.New("invalid max tokens")
)

type oaiClients struct {
	azureClient  azureopenai.Client
	openAIClient openai.Client
}

func newOAIClients() (oaiClients, error) {
	var (
		oaiClient   openai.Client
		azureClient azureopenai.Client
		err         error
	)

	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		oaiClient = openai.NewClient(*openAIAPIKey)
	} else {
		re := regexp.MustCompile(`^[a-zA-Z0-9]+([_-]?[a-zA-Z0-9]+)*$`)
		if !re.MatchString(*openAIDeploymentName) {
			return oaiClients{}, errors.New("azure openai deployment can only include alphanumeric characters, '_,-', and can't end with '_' or '-'")
		}

		azureClient, err = azureopenai.NewClient(*azureOpenAIEndpoint, *openAIAPIKey, *openAIDeploymentName)
		if err != nil {
			return oaiClients{}, fmt.Errorf("error create Azure client: %w", err)
		}
	}

	clients := oaiClients{
		azureClient:  azureClient,
		openAIClient: oaiClient,
	}

	return clients, nil
}

func completion(ctx context.Context, client oaiClients, prompts []string, deploymentName string, subcommand string) (string, error) {
	temp := float32(*temperature)

	maxTokens, err := calculateMaxTokens(prompts, deploymentName)
	if err != nil {
		return "", fmt.Errorf("error calculate max token: %w", err)
	}

	var prompt strings.Builder
	_, err = fmt.Fprint(&prompt, subcommand)

	if err != nil {
		return "", fmt.Errorf("error prompt string builder: %w", err)
	}

	for _, p := range prompts {
		_, err = fmt.Fprintf(&prompt, "%s\n", p)
		if err != nil {
			return "", fmt.Errorf("error range prompt: %w", err)
		}
	}

	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		if isGptTurbo(deploymentName) || isGpt4(deploymentName) {
			resp, err := client.openaiGptChatCompletion(ctx, prompt, maxTokens, temp)
			if err != nil {
				return "", fmt.Errorf("error openai GptChat completion: %w", err)
			}

			return resp, nil
		}

		resp, err := client.openaiGptCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", fmt.Errorf("error openai Gpt completion: %w", err)
		}

		return resp, nil
	}

	if isGptTurbo35(deploymentName) || isGpt4(deploymentName) {
		resp, err := client.azureGptChatCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", fmt.Errorf("error azure GptChat completion: %w", err)
		}

		return resp, nil
	}

	resp, err := client.azureGptCompletion(ctx, prompt, maxTokens, temp)
	if err != nil {
		return "", fmt.Errorf("error azure Gpt completion: %w", err)
	}

	return resp, nil
}

func isGptTurbo(deploymentName string) bool {
	return deploymentName == "gpt-3.5-turbo-0301" || deploymentName == "gpt-3.5-turbo"
}

func isGptTurbo35(deploymentName string) bool {
	return deploymentName == "gpt-35-turbo-0301" || deploymentName == "gpt-35-turbo"
}

func isGpt4(deploymentName string) bool {
	return deploymentName == "gpt-4-0314" || deploymentName == "gpt-4-32k-0314"
}

func calculateMaxTokens(prompts []string, deploymentName string) (*int, error) {
	maxTokensFinal, ok := maxTokensMap[deploymentName]
	if !ok {
		return nil, errors.Wrapf(errToken, "deploymentName %q not found in max tokens map", deploymentName)
	}

	if *maxTokens > 0 {
		maxTokensFinal = *maxTokens
	}

	encoder, err := gptEncoder.NewEncoder()
	if err != nil {
		return nil, fmt.Errorf("error encode gpt: %w", err)
	}

	// start at 100 since the encoder at times doesn't get it exactly correct
	totalTokens := 100

	for _, prompt := range prompts {
		tokens, err := encoder.Encode(prompt)
		if err != nil {
			return nil, fmt.Errorf("error encode prompt: %w", err)
		}

		totalTokens += len(tokens)
	}

	remainingTokens := maxTokensFinal - totalTokens

	return &remainingTokens, nil
}
