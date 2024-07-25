/*
Copyright © 2023 Rick
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	utils "github.com/richardfaulkner-qz/rft/internal"
	"github.com/spf13/cobra"
)

type config struct {
	dryRun      bool
	force       bool
	includeDirs bool
	silence     bool
}

var dateFormatStr string = "2006-01-02"
var c config

// organizeDirCmd represents the organizeDir command
var organizeDirCmd = &cobra.Command{
	Use:   "organizeDir",
	Short: "Quickly organize a full directory, non-recursivly, by day",
	Long:  "This util is design to quickly organize a directory by day, non-recursivly. ",
	Run: func(cmd *cobra.Command, args []string) {
		if !c.silence {
			cmd.Println("Organizing directory...")
		}

		err := cleanup(c)

		if !c.silence {
			if err != nil {
				cmd.Printf("Error: %v\n", err)
			} else {
				cmd.Println("Done")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(organizeDirCmd)

	c.dryRun = *organizeDirCmd.Flags().BoolP("dry", "d", false, "Do a dry run and printout the stats for the operation")
	c.force = *organizeDirCmd.Flags().BoolP("force", "f", false, "Force the cleanup, if true the cleanup will push through errors to continue")
	c.includeDirs = *organizeDirCmd.Flags().BoolP("includeDirs", "r", false, "Include directories in the cleanup, if true directories will be moved")
	c.silence = *organizeDirCmd.Flags().BoolP("silent", "s", false, "Silence text outputs")
}

func cleanup(c config) error {
	dateStr := time.Now().Format(dateFormatStr)

	// Create the wrapper directory
	dirName := dateStr + " -- cleanup"

	// Map of creation dates to lists of files
	filesByDate := make(map[string][]string)
	fileCount := 0

	// Walk the current directory
	if err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, "/") { // skip subdirectories
			return nil
		}
		fileCount++
		if err != nil {
			return err
		}

		// Skip directories if moveDirs is false
		if !c.includeDirs && info.IsDir() {
			return nil
		}

		// Get the file's creation date
		creationDate := info.ModTime().Format(dateFormatStr)

		// Add the file to the appropriate list
		filesByDate[creationDate] = append(filesByDate[creationDate], path)

		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	if !c.silence {
		fmt.Printf("Will move %+v files to %+v directories \n", fileCount, len(filesByDate))
	}

	if !c.force {
		if resp, err := utils.GetYNUserInput(); err != nil || !resp {
			if resp {
				if !c.silence {
					fmt.Println("Cleanup canceled")
				}
				return nil
			}
			return err
		}

	}

	if !c.dryRun {
		if err := os.Mkdir(dirName, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Move the files into their respective folders
	for date, files := range filesByDate {
		if len(files) == 0 {
			continue
		}

		// Create the date directory if not a dry run
		if !c.dryRun {
			if err := os.Mkdir(filepath.Join(dirName, date), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		}

		// Move each file into the date directory if not a dry run
		for _, file := range files {
			if !c.dryRun {
				if err := os.Rename(file, filepath.Join(dirName, date, filepath.Base(file))); err != nil {
					return fmt.Errorf("failed to move file: %v", err)
				}
			}
		}
	}

	return nil
}
