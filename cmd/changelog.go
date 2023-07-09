/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.PlainOpen(".")
		if err != nil {
			log.Fatal(err)
		}

		headRef, err := repo.Head()
		if err != nil {
			log.Fatal(err)
		}

		currentBranch := headRef.Name().Short()
		fmt.Println(currentBranch)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(currentBranch), headRef.Hash())

		err = repo.Storer.SetReference(ref)
		if err != nil {
			log.Fatal(err)
		}

		gitLog := exec.Command("git", "log", "--no-merges", "^main")
		output, err := gitLog.Output()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(output))

		// repo.
		// commits, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// err = commits.ForEach(func(c *object.Commit) error {
		// 	fmt.Printf("Commit: %s\n", c.Hash)
		// 	fmt.Printf("Author: %s <%s>\n", c.Author.Name, c.Author.Email)
		// 	fmt.Printf("Date: %s\n", c.Author.When)
		// 	fmt.Printf("Message: %s\n\n", c.Message)
		// 	return nil
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changelogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changelogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
