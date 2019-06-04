package main

import (
	"github.com/spf13/cobra"
	"github.com/terryding77/gocz/commitizen/commit"
	"github.com/terryding77/gocz/commitizen/config"
)

func setupGoczCommit() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cz(commitizen)",
	}
	commitArgs := commit.NewDefaultArgs()
	cmd.Flags().BoolVarP(&commitArgs.Quiet, "quiet", "q", commitArgs.Quiet, "suppress summary after successful commit")
	cmd.Flags().BoolVarP(&commitArgs.Verbose, "verbose", "v", commitArgs.Verbose, "show diff in commit message template")
	cmd.Flags().StringVar(&commitArgs.Author, "author", commitArgs.Author, "override author for commit")
	cmd.Flags().StringVar(&commitArgs.Date, "date", commitArgs.Date, "override date for commit")
	cmd.Flags().BoolVarP(&commitArgs.All, "all", "a", commitArgs.All, "commit all changed files")
	cmd.Flags().BoolVar(&commitArgs.Amend, "amend", commitArgs.Amend, "amend previous commit")

	cmd.RunE = func(cobraCmd *cobra.Command, args []string) (err error) {
		// setup config
		cfg := config.New()
		if err := cfg.Setup(""); err != nil {
			return err
		}

		// setup commit
		cmt := commit.New()
		if err := cmt.Setup(cfg); err != nil {
			return err
		}
		if _, err = cmt.Read(); err != nil {
			return err
		}
		if err = cmt.Execute(commitArgs); err != nil {
			return err
		}

		return nil
	}
	return cmd
}

func main() {
	setupGoczCommit().Execute()
}
