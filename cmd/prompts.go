package cmd

import (
	"errors"
	"fmt"
)

const (
	Tone       = "tone"
	Emoji      = "emoji"
	Brainstorm = "brainstorm"
)

type LLMInteractionInterface interface {
	GeneratePrompt(*string, *string) (string, error)
	GetMaxTokens() int
}

type DefaultInteractionSetup struct {
	promptTemplate    string
	MaxTokens         int  `default:"200"`
	requiresUserInput bool `default:"true"`
}

var prompts = map[string]LLMInteractionInterface{
	Tone: DefaultInteractionSetup{
		promptTemplate: "You are an expert in sentiment and text analysis. First what is the general sentiment and tone of this string: ```%s``` \n %s",
	},
	Emoji: DefaultInteractionSetup{
		promptTemplate: "You are an expert in text analysis and adding relivant emojies to drive home the sentiment and or meaning of the text. Analyze this string and add emojis to it. %s string: ```%s``` \n",
	},
	Brainstorm: DefaultInteractionSetup{
		promptTemplate: "You are an expert in brainstorming. Brainstorm ideas for this string: ```%s``` prioritize asking follow up questions that help guide the user to deeper thinking around the subjsect. For any counter arguments or objections you give present reasons why those arguments are applicible and should be thought of. Your goal is not to shut down or discurage the user, its to brainstorm a larger volume of ideas so the user can get their ideas out and start to recognize the different paths available. \n",
		MaxTokens:      500,
	},
}

func GetInteractionSetup(key string) LLMInteractionInterface {
	return prompts[key]
}

func (p DefaultInteractionSetup) GeneratePrompt(userInput *string, additionalContext *string) (string, error) {
	if userInput == nil && p.requiresUserInput {
		return "", errors.New("No user input provided")
	}
	additionalContextString := generateAdditionalContext(additionalContext)

	prompt := fmt.Sprintf(p.promptTemplate, *userInput, additionalContextString)
	return prompt, nil
}

func (p DefaultInteractionSetup) GetMaxTokens() int {
	return p.MaxTokens
}

func generateAdditionalContext(additionalContext *string) string {
	if additionalContext != nil && len(*additionalContext) > 0 {
		return fmt.Sprintf("Some additional context for this analysis: ```%s```", *additionalContext)
	}
	return ""
}
