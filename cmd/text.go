/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var additionalContext *string
var shouldIterate *bool

// textCmds represents the textAnalysis command
var textCmds = &cobra.Command{
	Use:   "text",
	Short: "Tools for text analysis and other operations",
	Long: `
	Tools for text analysis and other operations`,
}

func init() {
	textCmds.AddCommand(textToneCmd)
	textCmds.AddCommand(emojiCmd)
	textCmds.AddCommand(brainstormCmd)

	additionalContext = textCmds.PersistentFlags().StringP("context", "c", "", "Additional context to be used in analysis")
	shouldIterate = textCmds.PersistentFlags().BoolP("iterable", "i", true, "Whether or not you want to iterate after the first response, true by default")
}

var textToneCmd = &cobra.Command{
	Use:   "tone",
	Short: "Sentiment/Tone analysis of the text provided",
	Run: func(cmd *cobra.Command, args []string) {
		if evalString, err := getUserText(args); err == nil {
			haveCoversation(&UserInteraction{
				InteractionType:   GetInteractionSetup(Tone),
				ShouldIterate:     *shouldIterate,
				EvalString:        evalString,
				AdditionalContext: *additionalContext})
		} else {
			fmt.Printf("\n Error getting user input: %v", err)
		}
	},
}

var emojiCmd = &cobra.Command{
	Use:   "emoji",
	Short: "Emoji analysis of the text provided",
	Run: func(cmd *cobra.Command, args []string) {
		if evalString, err := getPromptUserForInput(); err == nil {
			haveCoversation(&UserInteraction{
				InteractionType:   GetInteractionSetup(Emoji),
				ShouldIterate:     *shouldIterate,
				EvalString:        evalString,
				AdditionalContext: *additionalContext})
		} else {
			fmt.Printf("\n Error getting user input: %v", err)
		}
	},
}

var brainstormCmd = &cobra.Command{
	Use:   "brainstorm",
	Short: "Brainstorming tool back by an LLM -- always iterative",
	Long:  "A brainstorming tool primed to help ask clarifying questions and generate ideas. The primary purpose of this tool is to help you 'get the creative juices flowing'.",
	Run: func(cmd *cobra.Command, args []string) {
		if evalString, err := getPromptUserForInput(); err == nil {
			haveCoversation(&UserInteraction{
				InteractionType:   GetInteractionSetup(Brainstorm),
				ShouldIterate:     true,
				EvalString:        evalString,
				AdditionalContext: *additionalContext})
		} else {
			fmt.Printf("\n Error getting user input: %v", err)
		}
	},
}

func getUserText(args []string) (string, error) {
	if len(args) < 1 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter Text: ")
		scanner.Scan()

		fmt.Println(scanner.Text())
		return scanner.Text(), scanner.Err()
	} else {
		return args[0], nil
	}
}

func getPromptUserForInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Text: ")
	scanner.Scan()

	fmt.Println(scanner.Text())
	return scanner.Text(), scanner.Err()
}
