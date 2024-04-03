/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var utilsKillPortCmd = &cobra.Command{
	Use:   "kill-port",
	Short: "Kill a pesky task running on a port",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(1, 10)(cmd, args); err != nil {
			return err
		}

		for i := range len(args) {
			if _, err := strconv.Atoi(args[i]); err != nil {
				return fmt.Errorf("Invalid port number supplied '%v'", args[i])
			}
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			fmt.Printf("killing port: %v \n", v)

			cmd := exec.Command("kill", fmt.Sprintf("$(lsof -ti:%v)", v))

			if err := cmd.Run(); err != nil {
				fmt.Println("error killing port. all excecution halted", err)
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(utilsKillPortCmd)
}
