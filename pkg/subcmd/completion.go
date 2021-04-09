package subcmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "shell completion",
	}

	cmd.AddCommand(
		newBashCompletionCmd(),
		newZshCompletionCmd(),
	)

	return cmd
}

func newBashCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenBashCompletion(os.Stdout)
		},
	}

	return cmd
}

func newZshCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenZshCompletion(os.Stdout)
		},
	}

	return cmd
}
