/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// organizeDirCmd represents the organizeDir command
var organizeDirCmd = &cobra.Command{
	Use:   "organizeDir",
	Short: "Quickly organize a full directory, non-recursivly, by day",
	Long:  "Util to keep directories clean, by default groups all files in that dir by day,",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("organizeDir called")
		includeDirs, _ := cmd.Flags().GetBool("includeDirs")
		dryRun, _ := cmd.Flags().GetBool("dry")

		_ = cleanup(dryRun, includeDirs)
	},
}

func init() {
	rootCmd.AddCommand(organizeDirCmd)

	organizeDirCmd.Flags().BoolP("dry", "d", false, "Do a dryrun and printout the stats for the operation") // Here you will define your flags and configuration settings.
	organizeDirCmd.Flags().BoolP("wrapper", "w", false, "Add wrapper around clenaup")
	organizeDirCmd.Flags().Bool("includeDirs", false, "Include directories in the cleanup, if true directories will be moved")
}

func cleanup(dryRun bool, moveDirs bool) error {
	// Get the current date
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	// Create the wrapper directory
	dirName := dateStr + " -- cleanup"
	if !dryRun {
		if err := os.Mkdir(dirName, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Map of creation dates to lists of files
	filesByDate := make(map[string][]string)

	// Walk the current directory
	if err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories if moveDirs is false
		if !moveDirs && info.IsDir() {
			return nil
		}

		// Get the file's creation date
		creationTime := info.ModTime()
		creationDate := creationTime.Format("2006-01-02")

		// Add the file to the appropriate list
		filesByDate[creationDate] = append(filesByDate[creationDate], path)

		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	// Move the files into their respective folders
	for date, files := range filesByDate {
		// Skip dates with no files if not a dry run
		if len(files) == 0 && !dryRun {
			continue
		}

		// Create the date directory if not a dry run
		if !dryRun {
			if err := os.Mkdir(filepath.Join(dirName, date), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		}

		// Move each file into the date directory if not a dry run
		for _, file := range files {
			if dryRun {
				fmt.Printf("Would move %s to %s\n", file, filepath.Join(dirName, date))
			} else {
				if err := os.Rename(file, filepath.Join(dirName, date, filepath.Base(file))); err != nil {
					return fmt.Errorf("failed to move file: %v", err)
				}
			}
		}
	}

	return nil
}
