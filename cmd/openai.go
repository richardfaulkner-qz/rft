package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Read in the openai token from an environment variable
var openAIToken string = os.Getenv("OPENAI_API_KEY")

func ToneEval(evalString *string, additionalContext *string, shouldIterate *bool) {
	prompt := fmt.Sprintf("You are an expert in sentiment and text analysis. First what is the general sentiment and tone of this string: ```%s``` \n %s", *evalString, generateAdditionalContext(additionalContext))

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: prompt},
	}

	for {
		fmt.Println("\nAssistant: ")
		response := makeOpenAICall(&messages)
		if !*shouldIterate || len(response) == 0 {
			break
		}
		messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: response})

		// Get user input`
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("\n\nUser: ")
		scanner.Scan()

		if len(scanner.Text()) == 0 {
			break
		}
		if scanner.Err() != nil {
			fmt.Println("Error: ", scanner.Err())
			break
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: scanner.Text(),
		})

	}
}

func generateAdditionalContext(additionalContext *string) string {
	if additionalContext != nil && len(*additionalContext) > 0 {
		return fmt.Sprintf("Some additional context for this analysis: ```%s```", *additionalContext)
	}
	return ""
}

func makeOpenAICall(messages *[]openai.ChatCompletionMessage) string {
	c := openai.NewClient(openAIToken)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    *messages,
		MaxTokens:   200,
		Temperature: 0.7,
		Stream:      false,
	}

	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("CompletionStream error: %v\n", err)
		return ""
	}
	defer stream.Close()

	responseMessage := ""
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return responseMessage
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return ""
		}

		responseMessage += response.Choices[0].Delta.Content
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
