package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Execute() {
	var command = &cobra.Command{
		Use:   "post",
		Short: "Post API Application",
		Long:  `A modular blog post API built with Go, Gin, and GORM.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(apiCmd())

	if err := command.Execute(); err != nil {
		log.Fatal().Msgf("failed run app: %s", err.Error())
	}
}
