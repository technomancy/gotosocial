package main

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/superseriousbusiness/gotosocial/internal/plugin"
)

func pluginCommands() *cobra.Command {
	pluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "run plugins from the command line",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(preRunArgs{cmd: cmd})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.WithValue(cmd.Context(), "args", args)
			return run(ctx, plugin.Run)
		},
	}

	return pluginCmd
}
