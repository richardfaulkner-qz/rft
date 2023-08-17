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
	textCmds.AddCommand(textSentimentCmd)

	additionalContext = textCmds.PersistentFlags().StringP("context", "c", "", "Additional context to be used in analysis")
	shouldIterate = textCmds.PersistentFlags().BoolP("iterable", "i", false, "Whether or not you want to iterate after the first response")
}

var textSentimentCmd = &cobra.Command{
	Use:   "tone",
	Short: "Sentiment/Tone analysis of the text provided",
	Run: func(cmd *cobra.Command, args []string) {
		textToAnalyze := ""

		if len(args) < 1 {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter Text: ")
			scanner.Scan()

			fmt.Println(scanner.Text())
			textToAnalyze = scanner.Text()

			if scanner.Err() != nil {
				fmt.Println("Error: ", scanner.Err())
				return
			}
		} else {
			textToAnalyze = args[0]
		}

		// Make a call to the sentiment analysis API
		ToneEval(&textToAnalyze, additionalContext, shouldIterate)
	},
}
