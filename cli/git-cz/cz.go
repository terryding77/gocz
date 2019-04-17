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
		commitArgs := commit.NewDefaultArgs()
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
