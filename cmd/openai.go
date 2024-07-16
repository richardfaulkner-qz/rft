package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	utils "github.com/richardfaulkner-qz/rft/internal"
	openai "github.com/sashabaranov/go-openai"
)

type UserInteraction struct {
	InteractionType   LLMInteractionInterface
	ShouldIterate     bool
	EvalString        string
	AdditionalContext string
}

// Read in the openai token from an environment variable
var openAIToken string = os.Getenv("OPENAI_API_KEY")

func haveCoversation(interaction *UserInteraction) {
	promptString, err := interaction.InteractionType.GeneratePrompt(&interaction.EvalString, &interaction.AdditionalContext)
	if err != nil {
		fmt.Println(err)
		return
	}
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: promptString},
	}

	for {
		fmt.Println("\nAssistant: ")
		response := makeOpenAICall(&messages, interaction.InteractionType.GetMaxTokens())
		if !interaction.ShouldIterate || len(response) == 0 {
			break
		}
		messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: response})

		if usermessage, err := utils.GetUserChatInput(); err != nil || len(usermessage) == 0 {
			break
		} else {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: usermessage,
			})
		}
	}

}

func makeOpenAICall(messages *[]openai.ChatCompletionMessage, maxTokens int) string {
	c := openai.NewClient(openAIToken)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    *messages,
		MaxTokens:   maxTokens,
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
