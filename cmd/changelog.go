/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type ParsedCommit struct {
	Hash, ShortHash, Body string
	Date                  time.Time
	ParsedBody            conventionalcommits.Message
}

func (p ParsedCommit) String() string {
	return fmt.Sprintf("Date: %s\tHash: %s\tShortHash: %s\tBody: %s\t Parsing OK : %t", p.Date.Format("2006-01-02 15:04:05"), p.Hash, p.ShortHash, p.Body, p.ParsedBody.Ok())
}

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
		ref := plumbing.NewHashReference(plumbing.ReferenceName(currentBranch), headRef.Hash())

		err = repo.Storer.SetReference(ref)
		if err != nil {
			log.Fatal(err)
		}

		format := "%H\xbd%h\xbd%ct\xbd%B"

		gitLog := exec.Command("git", "log", "--no-merges", "--format="+format, currentBranch, "^main")
		output, err := gitLog.Output()
		if err != nil {
			log.Fatal(err)
		}

		p := parser.NewMachine(parser.WithBestEffort(), conventionalcommits.WithTypes(conventionalcommits.TypesFalco))

		var parsedCommits []ParsedCommit
		for _, v := range strings.Split(string(output), "\n") {
			if v == "" {
				continue
			}
			splitted := strings.Split(v, "\xbd")
			fmt.Println(splitted, len(splitted))
			timestamp, err := strconv.ParseInt(splitted[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			parsedBody, err := p.Parse([]byte(splitted[3]))
			if err != nil {
				log.Println(err)
			}
			commit := &ParsedCommit{
				Hash:       splitted[0],
				ShortHash:  splitted[1],
				Date:       time.Unix(timestamp, 0),
				Body:       splitted[3],
				ParsedBody: parsedBody,
			}
			parsedCommits = append(parsedCommits, *commit)
		}

		lo.ForEach(parsedCommits, func(pc ParsedCommit, _ int) {
			fmt.Println(pc)
		})

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
