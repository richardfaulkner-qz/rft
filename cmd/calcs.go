/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// add two floating point var
	num_t           = 0
	val     float32 = 0
	verbose         = false
)

// calcsCmd represents the calcs command
var calcsCmd = &cobra.Command{
	Use:   "calcs",
	Short: "Calculates the total dollar value needed based on the starting value and the number of territories remaining.",
	Run: func(cmd *cobra.Command, args []string) {
		var sum float32 = 0.0
		for i := 0; i < num_t; i++ {
			sum += val
			val = val * 1.05
			if verbose {
				fmt.Printf("Territory %v has a value of %v\n", i+1, val)
			}
		}

		fmt.Printf("Total value of all territories is %v\n", sum)
	},
}

func init() {
	rootCmd.AddCommand(calcsCmd)

	calcsCmd.Flags().BoolP("verbose", "v", false, "print in verbose mode")
	calcsCmd.Flags().IntVarP(&num_t, "num_t", "n", 0, "Number of territories use")
	calcsCmd.Flags().Float32VarP(&val, "start_val", "s", 800.0, "Initial value of territory")
}
